package separatorBasedProtocol

import (
	"bufio"
	"bytes"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	tcpServer "github.com/thewizardplusplus/go-tcp-server"
	defaultProtocol "github.com/thewizardplusplus/go-tcp-server/protocols/default"
	defaultProtocolModels "github.com/thewizardplusplus/go-tcp-server/protocols/default/models"
)

func TestProtocol_interface(test *testing.T) {
	assert.Implements(
		test,
		(*tcpServer.ServerProtocol[
			defaultProtocolModels.Request,
			defaultProtocolModels.Response,
		])(nil),
		Protocol{},
	)
	assert.Implements(
		test,
		(*tcpServer.ClientProtocol[
			defaultProtocolModels.Request,
			defaultProtocolModels.Response,
		])(nil),
		Protocol{},
	)
}

func TestNewProtocol(test *testing.T) {
	type args struct {
		options SeparationParams
	}

	for _, data := range []struct {
		name string
		args args
		want Protocol
	}{
		{
			name: "success",
			args: args{
				options: SeparationParams{
					MessageSeparator:        []byte("\n"),
					MessagePartSeparator:    []byte("|"),
					HeaderSeparator:         []byte("&"),
					HeaderKeyValueSeparator: []byte("="),
				},
			},
			want: Protocol{
				BaseProtocol: defaultProtocol.NewBaseProtocol(
					defaultProtocol.BaseProtocolOptions{
						MessageFormat: NewMessageFormat(SeparationParams{
							MessageSeparator:        []byte("\n"),
							MessagePartSeparator:    []byte("|"),
							HeaderSeparator:         []byte("&"),
							HeaderKeyValueSeparator: []byte("="),
						}),
					},
				),
				options: SeparationParams{
					MessageSeparator:        []byte("\n"),
					MessagePartSeparator:    []byte("|"),
					HeaderSeparator:         []byte("&"),
					HeaderKeyValueSeparator: []byte("="),
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

func TestProtocol_ExtractToken(test *testing.T) {
	type fields struct {
		options SeparationParams
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
				options: SeparationParams{
					MessageSeparator:        []byte("\n"),
					MessagePartSeparator:    []byte("|"),
					HeaderSeparator:         []byte("&"),
					HeaderKeyValueSeparator: []byte("="),
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
				options: SeparationParams{
					MessageSeparator:        []byte("\n"),
					MessagePartSeparator:    []byte("|"),
					HeaderSeparator:         []byte("&"),
					HeaderKeyValueSeparator: []byte("="),
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
				options: SeparationParams{
					MessageSeparator:        []byte("\n"),
					MessagePartSeparator:    []byte("|"),
					HeaderSeparator:         []byte("&"),
					HeaderKeyValueSeparator: []byte("="),
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
