package defaultProtocol

import (
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	defaultBaseProtocolMocks "github.com/thewizardplusplus/go-tcp-server/mocks/github.com/thewizardplusplus/go-tcp-server/protocols/default"
	defaultProtocolModels "github.com/thewizardplusplus/go-tcp-server/protocols/default/models"
	defaultProtocolModelValueTypes "github.com/thewizardplusplus/go-tcp-server/protocols/default/models/value-types"
)

func TestNewBaseProtocol(test *testing.T) {
	type args struct {
		options func(test *testing.T) BaseProtocolOptions
	}

	for _, data := range []struct {
		name string
		args args
		want func(test *testing.T) BaseProtocol
	}{
		{
			name: "success",
			args: args{
				options: func(test *testing.T) BaseProtocolOptions {
					return BaseProtocolOptions{
						MessageFormat: defaultBaseProtocolMocks.NewMockMessageFormat(test),
					}
				},
			},
			want: func(test *testing.T) BaseProtocol {
				return BaseProtocol{
					options: BaseProtocolOptions{
						MessageFormat: defaultBaseProtocolMocks.NewMockMessageFormat(test),
					},
				}
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got := NewBaseProtocol(data.args.options(test))

			assert.Equal(test, data.want(test), got)
		})
	}
}

func TestBaseProtocol_InitialScannerBufferSize(test *testing.T) {
	for _, data := range []struct {
		name string
		want int
	}{
		{
			name: "success",
			want: 4 * 1024,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got := (BaseProtocol{}).InitialScannerBufferSize()

			assert.Equal(test, data.want, got)
		})
	}
}

func TestBaseProtocol_MaxTokenSize(test *testing.T) {
	for _, data := range []struct {
		name string
		want int
	}{
		{
			name: "success",
			want: 64 * 1024,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got := (BaseProtocol{}).MaxTokenSize()

			assert.Equal(test, data.want, got)
		})
	}
}

func TestBaseProtocol_ParseRequest(test *testing.T) {
	type fields struct {
		options func(test *testing.T) BaseProtocolOptions
	}
	type args struct {
		data []byte
	}

	for _, data := range []struct {
		name    string
		fields  fields
		args    args
		want    defaultProtocolModels.Request
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			fields: fields{
				options: func(test *testing.T) BaseProtocolOptions {
					introduction, err :=
						defaultProtocolModelValueTypes.NewIntroduction([]byte("introduction"))
					require.NoError(test, err)

					message, err := defaultProtocolModels.NewMessageBuilder().
						SetIntroduction(introduction).
						Build()
					require.NoError(test, err)

					messageFormatMock := defaultBaseProtocolMocks.NewMockMessageFormat(test)
					messageFormatMock.EXPECT().
						ParseMessage([]byte("dummy")).
						Return(message, nil)

					return BaseProtocolOptions{
						MessageFormat: messageFormatMock,
					}
				},
			},
			args: args{
				data: []byte("dummy"),
			},
			want: func() defaultProtocolModels.Request {
				action, err :=
					defaultProtocolModelValueTypes.NewAction([]byte("introduction"))
				require.NoError(test, err)

				request, err := defaultProtocolModels.NewRequestBuilder().
					SetAction(action).
					Build()
				require.NoError(test, err)

				return request
			}(),
			wantErr: assert.NoError,
		},
		{
			name: "error/unable to parse the message",
			fields: fields{
				options: func(test *testing.T) BaseProtocolOptions {
					messageFormatMock := defaultBaseProtocolMocks.NewMockMessageFormat(test)
					messageFormatMock.EXPECT().
						ParseMessage([]byte("dummy")).
						Return(defaultProtocolModels.Message{}, iotest.ErrTimeout)

					return BaseProtocolOptions{
						MessageFormat: messageFormatMock,
					}
				},
			},
			args: args{
				data: []byte("dummy"),
			},
			want:    defaultProtocolModels.Request{},
			wantErr: assert.Error,
		},
		{
			name: "error/unable to construct the request",
			fields: fields{
				options: func(test *testing.T) BaseProtocolOptions {
					messageFormatMock := defaultBaseProtocolMocks.NewMockMessageFormat(test)
					messageFormatMock.EXPECT().
						ParseMessage([]byte("dummy")).
						Return(defaultProtocolModels.Message{}, nil)

					return BaseProtocolOptions{
						MessageFormat: messageFormatMock,
					}
				},
			},
			args: args{
				data: []byte("dummy"),
			},
			want:    defaultProtocolModels.Request{},
			wantErr: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			protocol := BaseProtocol{
				options: data.fields.options(test),
			}
			got, err := protocol.ParseRequest(data.args.data)

			assert.Equal(test, data.want, got)
			data.wantErr(test, err)
		})
	}
}

func TestBaseProtocol_ParseResponse(test *testing.T) {
	type fields struct {
		options func(test *testing.T) BaseProtocolOptions
	}
	type args struct {
		data []byte
	}

	for _, data := range []struct {
		name    string
		fields  fields
		args    args
		want    defaultProtocolModels.Response
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			fields: fields{
				options: func(test *testing.T) BaseProtocolOptions {
					introduction, err :=
						defaultProtocolModelValueTypes.NewIntroduction([]byte("introduction"))
					require.NoError(test, err)

					message, err := defaultProtocolModels.NewMessageBuilder().
						SetIntroduction(introduction).
						Build()
					require.NoError(test, err)

					messageFormatMock := defaultBaseProtocolMocks.NewMockMessageFormat(test)
					messageFormatMock.EXPECT().
						ParseMessage([]byte("dummy")).
						Return(message, nil)

					return BaseProtocolOptions{
						MessageFormat: messageFormatMock,
					}
				},
			},
			args: args{
				data: []byte("dummy"),
			},
			want: func() defaultProtocolModels.Response {
				status, err :=
					defaultProtocolModelValueTypes.NewStatus([]byte("introduction"))
				require.NoError(test, err)

				response, err := defaultProtocolModels.NewResponseBuilder().
					SetStatus(status).
					Build()
				require.NoError(test, err)

				return response
			}(),
			wantErr: assert.NoError,
		},
		{
			name: "error/unable to parse the message",
			fields: fields{
				options: func(test *testing.T) BaseProtocolOptions {
					messageFormatMock := defaultBaseProtocolMocks.NewMockMessageFormat(test)
					messageFormatMock.EXPECT().
						ParseMessage([]byte("dummy")).
						Return(defaultProtocolModels.Message{}, iotest.ErrTimeout)

					return BaseProtocolOptions{
						MessageFormat: messageFormatMock,
					}
				},
			},
			args: args{
				data: []byte("dummy"),
			},
			want:    defaultProtocolModels.Response{},
			wantErr: assert.Error,
		},
		{
			name: "error/unable to construct the response",
			fields: fields{
				options: func(test *testing.T) BaseProtocolOptions {
					messageFormatMock := defaultBaseProtocolMocks.NewMockMessageFormat(test)
					messageFormatMock.EXPECT().
						ParseMessage([]byte("dummy")).
						Return(defaultProtocolModels.Message{}, nil)

					return BaseProtocolOptions{
						MessageFormat: messageFormatMock,
					}
				},
			},
			args: args{
				data: []byte("dummy"),
			},
			want:    defaultProtocolModels.Response{},
			wantErr: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			protocol := BaseProtocol{
				options: data.fields.options(test),
			}
			got, err := protocol.ParseResponse(data.args.data)

			assert.Equal(test, data.want, got)
			data.wantErr(test, err)
		})
	}
}

func TestBaseProtocol_MarshalRequest(test *testing.T) {
	type fields struct {
		options func(test *testing.T) BaseProtocolOptions
	}
	type args struct {
		request defaultProtocolModels.Request
	}

	for _, data := range []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			fields: fields{
				options: func(test *testing.T) BaseProtocolOptions {
					introduction, err :=
						defaultProtocolModelValueTypes.NewIntroduction([]byte("introduction"))
					require.NoError(test, err)

					message, err := defaultProtocolModels.NewMessageBuilder().
						SetIntroduction(introduction).
						Build()
					require.NoError(test, err)

					messageFormatMock := defaultBaseProtocolMocks.NewMockMessageFormat(test)
					messageFormatMock.EXPECT().
						MarshalMessage(message).
						Return([]byte("dummy"), nil)

					return BaseProtocolOptions{
						MessageFormat: messageFormatMock,
					}
				},
			},
			args: args{
				request: func() defaultProtocolModels.Request {
					action, err :=
						defaultProtocolModelValueTypes.NewAction([]byte("introduction"))
					require.NoError(test, err)

					request, err := defaultProtocolModels.NewRequestBuilder().
						SetAction(action).
						Build()
					require.NoError(test, err)

					return request
				}(),
			},
			want:    []byte("dummy"),
			wantErr: assert.NoError,
		},
		{
			name: "error/unable to convert the request to the message",
			fields: fields{
				options: func(test *testing.T) BaseProtocolOptions {
					return BaseProtocolOptions{
						MessageFormat: defaultBaseProtocolMocks.NewMockMessageFormat(test),
					}
				},
			},
			args: args{
				request: defaultProtocolModels.Request{},
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "error/unable to marshal the message",
			fields: fields{
				options: func(test *testing.T) BaseProtocolOptions {
					introduction, err :=
						defaultProtocolModelValueTypes.NewIntroduction([]byte("introduction"))
					require.NoError(test, err)

					message, err := defaultProtocolModels.NewMessageBuilder().
						SetIntroduction(introduction).
						Build()
					require.NoError(test, err)

					messageFormatMock := defaultBaseProtocolMocks.NewMockMessageFormat(test)
					messageFormatMock.EXPECT().
						MarshalMessage(message).
						Return(nil, iotest.ErrTimeout)

					return BaseProtocolOptions{
						MessageFormat: messageFormatMock,
					}
				},
			},
			args: args{
				request: func() defaultProtocolModels.Request {
					action, err :=
						defaultProtocolModelValueTypes.NewAction([]byte("introduction"))
					require.NoError(test, err)

					request, err := defaultProtocolModels.NewRequestBuilder().
						SetAction(action).
						Build()
					require.NoError(test, err)

					return request
				}(),
			},
			want:    nil,
			wantErr: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			protocol := BaseProtocol{
				options: data.fields.options(test),
			}
			got, err := protocol.MarshalRequest(data.args.request)

			assert.Equal(test, data.want, got)
			data.wantErr(test, err)
		})
	}
}

func TestBaseProtocol_MarshalResponse(test *testing.T) {
	type fields struct {
		options func(test *testing.T) BaseProtocolOptions
	}
	type args struct {
		response defaultProtocolModels.Response
	}

	for _, data := range []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			fields: fields{
				options: func(test *testing.T) BaseProtocolOptions {
					introduction, err :=
						defaultProtocolModelValueTypes.NewIntroduction([]byte("introduction"))
					require.NoError(test, err)

					message, err := defaultProtocolModels.NewMessageBuilder().
						SetIntroduction(introduction).
						Build()
					require.NoError(test, err)

					messageFormatMock := defaultBaseProtocolMocks.NewMockMessageFormat(test)
					messageFormatMock.EXPECT().
						MarshalMessage(message).
						Return([]byte("dummy"), nil)

					return BaseProtocolOptions{
						MessageFormat: messageFormatMock,
					}
				},
			},
			args: args{
				response: func() defaultProtocolModels.Response {
					status, err :=
						defaultProtocolModelValueTypes.NewStatus([]byte("introduction"))
					require.NoError(test, err)

					response, err := defaultProtocolModels.NewResponseBuilder().
						SetStatus(status).
						Build()
					require.NoError(test, err)

					return response
				}(),
			},
			want:    []byte("dummy"),
			wantErr: assert.NoError,
		},
		{
			name: "error/unable to convert the response to the message",
			fields: fields{
				options: func(test *testing.T) BaseProtocolOptions {
					return BaseProtocolOptions{
						MessageFormat: defaultBaseProtocolMocks.NewMockMessageFormat(test),
					}
				},
			},
			args: args{
				response: defaultProtocolModels.Response{},
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "error/unable to marshal the message",
			fields: fields{
				options: func(test *testing.T) BaseProtocolOptions {
					introduction, err :=
						defaultProtocolModelValueTypes.NewIntroduction([]byte("introduction"))
					require.NoError(test, err)

					message, err := defaultProtocolModels.NewMessageBuilder().
						SetIntroduction(introduction).
						Build()
					require.NoError(test, err)

					messageFormatMock := defaultBaseProtocolMocks.NewMockMessageFormat(test)
					messageFormatMock.EXPECT().
						MarshalMessage(message).
						Return(nil, iotest.ErrTimeout)

					return BaseProtocolOptions{
						MessageFormat: messageFormatMock,
					}
				},
			},
			args: args{
				response: func() defaultProtocolModels.Response {
					status, err :=
						defaultProtocolModelValueTypes.NewStatus([]byte("introduction"))
					require.NoError(test, err)

					response, err := defaultProtocolModels.NewResponseBuilder().
						SetStatus(status).
						Build()
					require.NoError(test, err)

					return response
				}(),
			},
			want:    nil,
			wantErr: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			protocol := BaseProtocol{
				options: data.fields.options(test),
			}
			got, err := protocol.MarshalResponse(data.args.response)

			assert.Equal(test, data.want, got)
			data.wantErr(test, err)
		})
	}
}
