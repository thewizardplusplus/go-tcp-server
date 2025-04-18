package separatorBasedTCPServerProtocol

import (
	"bufio"
	"bytes"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	tcpServer "github.com/thewizardplusplus/go-tcp-server"
)

func TestProtocol_interface(test *testing.T) {
	assert.Implements(
		test,
		(*tcpServer.ServerProtocol[Request, Response])(nil),
		Protocol{},
	)
}

func TestNewProtocol(test *testing.T) {
	type args struct {
		options ProtocolOptions
	}

	for _, data := range []struct {
		name string
		args args
		want Protocol
	}{
		{
			name: "success",
			args: args{
				options: ProtocolOptions{
					SeparationParams: SeparationParams{
						MessageSeparator:        []byte("\n"),
						MessagePartSeparator:    []byte("|"),
						HeaderSeparator:         []byte("&"),
						HeaderKeyValueSeparator: []byte("="),
					},
				},
			},
			want: Protocol{
				options: ProtocolOptions{
					SeparationParams: SeparationParams{
						MessageSeparator:        []byte("\n"),
						MessagePartSeparator:    []byte("|"),
						HeaderSeparator:         []byte("&"),
						HeaderKeyValueSeparator: []byte("="),
					},
				},
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got := NewProtocol(data.args.options)

			assert.Equal(test, data.want, got)
		})
	}
}

func TestProtocol_InitialScannerBufferSize(test *testing.T) {
	type fields struct {
		options ProtocolOptions
	}

	for _, data := range []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "success",
			fields: fields{
				options: ProtocolOptions{
					SeparationParams: SeparationParams{
						MessageSeparator:        []byte("\n"),
						MessagePartSeparator:    []byte("|"),
						HeaderSeparator:         []byte("&"),
						HeaderKeyValueSeparator: []byte("="),
					},
				},
			},
			want: 4 * 1024,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			protocol := Protocol{
				options: data.fields.options,
			}
			got := protocol.InitialScannerBufferSize()

			assert.Equal(test, data.want, got)
		})
	}
}

func TestProtocol_MaxTokenSize(test *testing.T) {
	type fields struct {
		options ProtocolOptions
	}

	for _, data := range []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "success",
			fields: fields{
				options: ProtocolOptions{
					SeparationParams: SeparationParams{
						MessageSeparator:        []byte("\n"),
						MessagePartSeparator:    []byte("|"),
						HeaderSeparator:         []byte("&"),
						HeaderKeyValueSeparator: []byte("="),
					},
				},
			},
			want: 64 * 1024,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			protocol := Protocol{
				options: data.fields.options,
			}
			got := protocol.MaxTokenSize()

			assert.Equal(test, data.want, got)
		})
	}
}

func TestProtocol_ExtractToken(test *testing.T) {
	type fields struct {
		options ProtocolOptions
	}
	type scannerParams struct {
		initialScannerBufferSize int
		maxTokenSize             int
	}

	for _, data := range []struct {
		name           string
		fields         fields
		scannerParams  scannerParams
		scannerInput   []byte
		wantTokens     [][]byte
		wantScannerErr assert.ErrorAssertionFunc
	}{
		{
			name: "success/with a trailing separator",
			fields: fields{
				options: ProtocolOptions{
					SeparationParams: SeparationParams{
						MessageSeparator:        []byte("\n"),
						MessagePartSeparator:    []byte("|"),
						HeaderSeparator:         []byte("&"),
						HeaderKeyValueSeparator: []byte("="),
					},
				},
			},
			scannerParams: scannerParams{
				initialScannerBufferSize: 4 * 1024,
				maxTokenSize:             64 * 1024,
			},
			scannerInput: []byte("dummy #0\ndummy #1\ndummy #2\n"),
			wantTokens: [][]byte{
				[]byte("dummy #0"),
				[]byte("dummy #1"),
				[]byte("dummy #2"),
			},
			wantScannerErr: assert.NoError,
		},
		{
			name: "success/without a trailing separator",
			fields: fields{
				options: ProtocolOptions{
					SeparationParams: SeparationParams{
						MessageSeparator:        []byte("\n"),
						MessagePartSeparator:    []byte("|"),
						HeaderSeparator:         []byte("&"),
						HeaderKeyValueSeparator: []byte("="),
					},
				},
			},
			scannerParams: scannerParams{
				initialScannerBufferSize: 4 * 1024,
				maxTokenSize:             64 * 1024,
			},
			scannerInput: []byte("dummy #0\ndummy #1\ndummy #2"),
			wantTokens: [][]byte{
				[]byte("dummy #0"),
				[]byte("dummy #1"),
				[]byte("dummy #2"),
			},
			wantScannerErr: assert.NoError,
		},
		{
			name: "success/with allocations",
			fields: fields{
				options: ProtocolOptions{
					SeparationParams: SeparationParams{
						MessageSeparator:        []byte("\n"),
						MessagePartSeparator:    []byte("|"),
						HeaderSeparator:         []byte("&"),
						HeaderKeyValueSeparator: []byte("="),
					},
				},
			},
			scannerParams: scannerParams{
				initialScannerBufferSize: 2,
				maxTokenSize:             64 * 1024,
			},
			scannerInput: []byte("dummy #0\ndummy #1\ndummy #2"),
			wantTokens: [][]byte{
				[]byte("dummy #0"),
				[]byte("dummy #1"),
				[]byte("dummy #2"),
			},
			wantScannerErr: assert.NoError,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			protocol := Protocol{
				options: data.fields.options,
			}

			scanner := bufio.NewScanner(bytes.NewReader(data.scannerInput))
			scanner.Buffer(
				make([]byte, data.scannerParams.initialScannerBufferSize),
				data.scannerParams.maxTokenSize,
			)
			scanner.Split(protocol.ExtractToken)

			var gotTokens [][]byte
			for scanner.Scan() {
				gotTokens = append(gotTokens, slices.Clone(scanner.Bytes()))
			}

			scannerErr := scanner.Err()

			assert.Equal(test, data.wantTokens, gotTokens)
			data.wantScannerErr(test, scannerErr)
		})
	}
}

func TestProtocol_ParseRequest(test *testing.T) {
	type fields struct {
		options ProtocolOptions
	}
	type args struct {
		data []byte
	}

	for _, data := range []struct {
		name    string
		fields  fields
		args    args
		want    Request
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success/regular request",
			fields: fields{
				options: ProtocolOptions{
					SeparationParams: SeparationParams{
						MessageSeparator:        []byte("\n"),
						MessagePartSeparator:    []byte("|"),
						HeaderSeparator:         []byte("&"),
						HeaderKeyValueSeparator: []byte("="),
					},
				},
			},
			args: args{
				data: []byte("action|one=two&three=four|body"),
			},
			want: Request{
				Action: []byte("action"),
				Headers: map[string][]byte{
					"6f6e65":     []byte("two"),
					"7468726565": []byte("four"),
				},
				Body: []byte("body"),
			},
			wantErr: assert.NoError,
		},
		{
			name: "success/with escaped separators",
			fields: fields{
				options: ProtocolOptions{
					SeparationParams: SeparationParams{
						MessageSeparator:        []byte("\n"),
						MessagePartSeparator:    []byte("|"),
						HeaderSeparator:         []byte("&"),
						HeaderKeyValueSeparator: []byte("="),
					},
				},
			},
			args: args{
				data: []byte(
					"dummy%7caction%0a|" +
						"dummy%3done%0a=dummy%26two%0a&dummy%3dthree%0a=dummy%26four%0a|" +
						"dummy%7cbody%0a",
				),
			},
			want: Request{
				Action: []byte("dummy|action\n"),
				Headers: map[string][]byte{
					"64756d6d793d6f6e650a":     []byte("dummy&two\n"),
					"64756d6d793d74687265650a": []byte("dummy&four\n"),
				},
				Body: []byte("dummy|body\n"),
			},
			wantErr: assert.NoError,
		},
		{
			name: "error",
			fields: fields{
				options: ProtocolOptions{
					SeparationParams: SeparationParams{
						MessageSeparator:        []byte("\n"),
						MessagePartSeparator:    []byte("|"),
						HeaderSeparator:         []byte("&"),
						HeaderKeyValueSeparator: []byte("="),
					},
				},
			},
			args: args{
				data: []byte("dummy"),
			},
			want:    Request{},
			wantErr: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			protocol := Protocol{
				options: data.fields.options,
			}
			got, err := protocol.ParseRequest(data.args.data)

			assert.Equal(test, data.want, got)
			data.wantErr(test, err)
		})
	}
}

func TestProtocol_ParseResponse(test *testing.T) {
	type fields struct {
		options ProtocolOptions
	}
	type args struct {
		data []byte
	}

	for _, data := range []struct {
		name    string
		fields  fields
		args    args
		want    Response
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success/regular response",
			fields: fields{
				options: ProtocolOptions{
					SeparationParams: SeparationParams{
						MessageSeparator:        []byte("\n"),
						MessagePartSeparator:    []byte("|"),
						HeaderSeparator:         []byte("&"),
						HeaderKeyValueSeparator: []byte("="),
					},
				},
			},
			args: args{
				data: []byte("status|one=two&three=four|body"),
			},
			want: Response{
				Status: []byte("status"),
				Headers: map[string][]byte{
					"6f6e65":     []byte("two"),
					"7468726565": []byte("four"),
				},
				Body: []byte("body"),
			},
			wantErr: assert.NoError,
		},
		{
			name: "success/with escaped separators",
			fields: fields{
				options: ProtocolOptions{
					SeparationParams: SeparationParams{
						MessageSeparator:        []byte("\n"),
						MessagePartSeparator:    []byte("|"),
						HeaderSeparator:         []byte("&"),
						HeaderKeyValueSeparator: []byte("="),
					},
				},
			},
			args: args{
				data: []byte(
					"dummy%7cstatus%0a|" +
						"dummy%3done%0a=dummy%26two%0a&dummy%3dthree%0a=dummy%26four%0a|" +
						"dummy%7cbody%0a",
				),
			},
			want: Response{
				Status: []byte("dummy|status\n"),
				Headers: map[string][]byte{
					"64756d6d793d6f6e650a":     []byte("dummy&two\n"),
					"64756d6d793d74687265650a": []byte("dummy&four\n"),
				},
				Body: []byte("dummy|body\n"),
			},
			wantErr: assert.NoError,
		},
		{
			name: "error",
			fields: fields{
				options: ProtocolOptions{
					SeparationParams: SeparationParams{
						MessageSeparator:        []byte("\n"),
						MessagePartSeparator:    []byte("|"),
						HeaderSeparator:         []byte("&"),
						HeaderKeyValueSeparator: []byte("="),
					},
				},
			},
			args: args{
				data: []byte("dummy"),
			},
			want:    Response{},
			wantErr: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			protocol := Protocol{
				options: data.fields.options,
			}
			got, err := protocol.ParseResponse(data.args.data)

			assert.Equal(test, data.want, got)
			data.wantErr(test, err)
		})
	}
}

func TestProtocol_MarshalRequest(test *testing.T) {
	type fields struct {
		options ProtocolOptions
	}
	type args struct {
		request Request
	}

	for _, data := range []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success/regular request",
			fields: fields{
				options: ProtocolOptions{
					SeparationParams: SeparationParams{
						MessageSeparator:        []byte("\n"),
						MessagePartSeparator:    []byte("|"),
						HeaderSeparator:         []byte("&"),
						HeaderKeyValueSeparator: []byte("="),
					},
				},
			},
			args: args{
				request: Request{
					Action: []byte("action"),
					Headers: map[string][]byte{
						"6f6e65":     []byte("two"),
						"7468726565": []byte("four"),
					},
					Body: []byte("body"),
				},
			},
			want:    []byte("action|one=two&three=four|body"),
			wantErr: assert.NoError,
		},
		{
			name: "success/with escaped separators",
			fields: fields{
				options: ProtocolOptions{
					SeparationParams: SeparationParams{
						MessageSeparator:        []byte("\n"),
						MessagePartSeparator:    []byte("|"),
						HeaderSeparator:         []byte("&"),
						HeaderKeyValueSeparator: []byte("="),
					},
				},
			},
			args: args{
				request: Request{
					Action: []byte("dummy|action\n"),
					Headers: map[string][]byte{
						"64756d6d793d6f6e650a":     []byte("dummy&two\n"),
						"64756d6d793d74687265650a": []byte("dummy&four\n"),
					},
					Body: []byte("dummy|body\n"),
				},
			},
			want: []byte(
				"dummy%7caction%0a|" +
					"dummy%3done%0a=dummy%26two%0a&dummy%3dthree%0a=dummy%26four%0a|" +
					"dummy%7cbody%0a",
			),
			wantErr: assert.NoError,
		},
		{
			name: "error",
			fields: fields{
				options: ProtocolOptions{
					SeparationParams: SeparationParams{
						MessageSeparator:        []byte("\n"),
						MessagePartSeparator:    []byte("|"),
						HeaderSeparator:         []byte("&"),
						HeaderKeyValueSeparator: []byte("="),
					},
				},
			},
			args: args{
				request: Request{
					Action: []byte("action"),
					Headers: map[string][]byte{
						"invalid": []byte("dummy"),
					},
					Body: []byte("body"),
				},
			},
			want:    nil,
			wantErr: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			protocol := Protocol{
				options: data.fields.options,
			}
			got, err := protocol.MarshalRequest(data.args.request)

			assert.Equal(test, data.want, got)
			data.wantErr(test, err)
		})
	}
}

func TestProtocol_MarshalResponse(test *testing.T) {
	type fields struct {
		options ProtocolOptions
	}
	type args struct {
		response Response
	}

	for _, data := range []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success/regular response",
			fields: fields{
				options: ProtocolOptions{
					SeparationParams: SeparationParams{
						MessageSeparator:        []byte("\n"),
						MessagePartSeparator:    []byte("|"),
						HeaderSeparator:         []byte("&"),
						HeaderKeyValueSeparator: []byte("="),
					},
				},
			},
			args: args{
				response: Response{
					Status: []byte("status"),
					Headers: map[string][]byte{
						"6f6e65":     []byte("two"),
						"7468726565": []byte("four"),
					},
					Body: []byte("body"),
				},
			},
			want:    []byte("status|one=two&three=four|body"),
			wantErr: assert.NoError,
		},
		{
			name: "success/with escaped separators",
			fields: fields{
				options: ProtocolOptions{
					SeparationParams: SeparationParams{
						MessageSeparator:        []byte("\n"),
						MessagePartSeparator:    []byte("|"),
						HeaderSeparator:         []byte("&"),
						HeaderKeyValueSeparator: []byte("="),
					},
				},
			},
			args: args{
				response: Response{
					Status: []byte("dummy|status\n"),
					Headers: map[string][]byte{
						"64756d6d793d6f6e650a":     []byte("dummy&two\n"),
						"64756d6d793d74687265650a": []byte("dummy&four\n"),
					},
					Body: []byte("dummy|body\n"),
				},
			},
			want: []byte(
				"dummy%7cstatus%0a|" +
					"dummy%3done%0a=dummy%26two%0a&dummy%3dthree%0a=dummy%26four%0a|" +
					"dummy%7cbody%0a",
			),
			wantErr: assert.NoError,
		},
		{
			name: "error",
			fields: fields{
				options: ProtocolOptions{
					SeparationParams: SeparationParams{
						MessageSeparator:        []byte("\n"),
						MessagePartSeparator:    []byte("|"),
						HeaderSeparator:         []byte("&"),
						HeaderKeyValueSeparator: []byte("="),
					},
				},
			},
			args: args{
				response: Response{
					Status: []byte("status"),
					Headers: map[string][]byte{
						"invalid": []byte("dummy"),
					},
					Body: []byte("body"),
				},
			},
			want:    nil,
			wantErr: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			protocol := Protocol{
				options: data.fields.options,
			}
			got, err := protocol.MarshalResponse(data.args.response)

			assert.Equal(test, data.want, got)
			data.wantErr(test, err)
		})
	}
}
