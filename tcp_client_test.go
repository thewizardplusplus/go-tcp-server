package tcpServer_test

import (
	"bufio"
	"context"
	"io"
	"net"
	"testing"
	"testing/iotest"
	"time"

	"github.com/samber/mo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	tcpServer "github.com/thewizardplusplus/go-tcp-server"
	tcpServerExternalMocks "github.com/thewizardplusplus/go-tcp-server/mocks/external/github.com/thewizardplusplus/go-tcp-server"
	tcpServerMocks "github.com/thewizardplusplus/go-tcp-server/mocks/github.com/thewizardplusplus/go-tcp-server"
)

func TestNewTCPClient(test *testing.T) {
	type request string
	type response string
	type args struct {
		ctx     context.Context
		address mo.Option[string]
		options func(test *testing.T) tcpServer.TCPClientOptions[request, response]
	}

	for _, data := range []struct {
		name    string
		args    args
		want    assert.ValueAssertionFunc
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			args: args{
				ctx:     context.Background(),
				address: mo.None[string](),
				options: func(
					test *testing.T,
				) tcpServer.TCPClientOptions[request, response] {
					clientProtocolMock :=
						tcpServerExternalMocks.NewMockClientProtocol[request, response](test)
					clientProtocolMock.EXPECT().
						InitialScannerBufferSize().
						Return(4096)
					clientProtocolMock.EXPECT().
						MaxTokenSize().
						Return(bufio.MaxScanTokenSize)

					return tcpServer.TCPClientOptions[request, response]{
						ReadTimeout:    mo.Some(5 * time.Minute),
						WriteTimeout:   mo.Some(12 * time.Minute),
						ClientProtocol: clientProtocolMock,
					}
				},
			},
			want:    assert.NotEmpty,
			wantErr: assert.NoError,
		},
		{
			name: "error/canceled context",
			args: args{
				ctx: func() context.Context {
					ctx, ctxCancel := context.WithCancel(context.Background())
					ctxCancel()

					return ctx
				}(),
				address: mo.None[string](),
				options: func(
					test *testing.T,
				) tcpServer.TCPClientOptions[request, response] {
					clientProtocolMock :=
						tcpServerExternalMocks.NewMockClientProtocol[request, response](test)
					return tcpServer.TCPClientOptions[request, response]{
						ReadTimeout:    mo.Some(5 * time.Minute),
						WriteTimeout:   mo.Some(12 * time.Minute),
						ClientProtocol: clientProtocolMock,
					}
				},
			},
			want:    assert.Empty,
			wantErr: assert.Error,
		},
		{
			name: "error/invalid address",
			args: args{
				ctx:     context.Background(),
				address: mo.Some("127.0.0.1"),
				options: func(
					test *testing.T,
				) tcpServer.TCPClientOptions[request, response] {
					clientProtocolMock :=
						tcpServerExternalMocks.NewMockClientProtocol[request, response](test)
					return tcpServer.TCPClientOptions[request, response]{
						ReadTimeout:    mo.Some(5 * time.Minute),
						WriteTimeout:   mo.Some(12 * time.Minute),
						ClientProtocol: clientProtocolMock,
					}
				},
			},
			want:    assert.Empty,
			wantErr: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			listener := runTestServer(test, "127.0.0.1:")
			defer listener.Close()

			got, err := tcpServer.NewTCPClient(
				data.args.ctx,
				data.args.address.OrElse(listener.Addr().String()),
				data.args.options(test),
			)
			if err == nil {
				defer got.Close()
			}

			data.want(test, got)
			data.wantErr(test, err)
		})
	}
}

func TestNewTCPClientFromConnection(test *testing.T) {
	type request string
	type response string
	type args struct {
		connection func(test *testing.T) net.Conn
		options    func(test *testing.T) tcpServer.TCPClientOptions[request, response]
	}

	for _, data := range []struct {
		name string
		args args
		want assert.ValueAssertionFunc
	}{
		{
			name: "success",
			args: args{
				connection: func(test *testing.T) net.Conn {
					return tcpServerMocks.NewMocknetConn(test)
				},
				options: func(
					test *testing.T,
				) tcpServer.TCPClientOptions[request, response] {
					clientProtocolMock :=
						tcpServerExternalMocks.NewMockClientProtocol[request, response](test)
					clientProtocolMock.EXPECT().
						InitialScannerBufferSize().
						Return(4096)
					clientProtocolMock.EXPECT().
						MaxTokenSize().
						Return(bufio.MaxScanTokenSize)

					return tcpServer.TCPClientOptions[request, response]{
						ReadTimeout:    mo.Some(5 * time.Minute),
						WriteTimeout:   mo.Some(12 * time.Minute),
						ClientProtocol: clientProtocolMock,
					}
				},
			},
			want: assert.NotEmpty,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got := tcpServer.NewTCPClientFromConnection(
				data.args.connection(test),
				data.args.options(test),
			)

			data.want(test, got)
		})
	}
}

func TestTCPClient_SendRequest(test *testing.T) {
	const acceptableDeadlineError = time.Minute

	type request string
	type response string
	type constructorArgs struct {
		connection func(test *testing.T) net.Conn
		options    func(test *testing.T) tcpServer.TCPClientOptions[request, response]
	}
	type args struct {
		request request
	}

	for _, data := range []struct {
		name            string
		constructorArgs constructorArgs
		args            args
		want            response
		wantErr         assert.ErrorAssertionFunc
	}{
		{
			name: "success/without all the timeouts",
			constructorArgs: constructorArgs{
				connection: func(test *testing.T) net.Conn {
					netConnMock := tcpServerMocks.NewMocknetConn(test)
					netConnMock.EXPECT().
						Read(mock.AnythingOfType("[]uint8")).
						RunAndReturn(func(buffer []byte) (int, error) {
							return copy(buffer, "raw-response"), io.EOF
						})
					netConnMock.EXPECT().
						Write([]byte("marshalled-request")).
						RunAndReturn(func(data []byte) (int, error) {
							return len(data), nil
						})

					return netConnMock
				},
				options: func(
					test *testing.T,
				) tcpServer.TCPClientOptions[request, response] {
					clientProtocolMock :=
						tcpServerExternalMocks.NewMockClientProtocol[request, response](test)
					clientProtocolMock.EXPECT().
						InitialScannerBufferSize().
						Return(4096)
					clientProtocolMock.EXPECT().
						MaxTokenSize().
						Return(bufio.MaxScanTokenSize)
					clientProtocolMock.EXPECT().
						ExtractToken([]byte("raw-response"), true).
						Return(len("raw-response"), []byte("response"), bufio.ErrFinalToken)
					clientProtocolMock.EXPECT().
						MarshalRequest(request("request")).
						Return([]byte("marshalled-request"), nil)
					clientProtocolMock.EXPECT().
						ParseResponse([]byte("response")).
						Return("parsed-response", nil)

					return tcpServer.TCPClientOptions[request, response]{
						ClientProtocol: clientProtocolMock,
					}
				},
			},
			args: args{
				request: "request",
			},
			want:    "parsed-response",
			wantErr: assert.NoError,
		},
		{
			name: "success/with all the timeouts",
			constructorArgs: constructorArgs{
				connection: func(test *testing.T) net.Conn {
					netConnMock := tcpServerMocks.NewMocknetConn(test)
					netConnMock.EXPECT().
						SetReadDeadline(mock.MatchedBy(func(deadline time.Time) bool {
							expectedDeadline := time.Now().Add(5 * time.Minute)
							return expectedDeadline.Sub(deadline).Abs() < acceptableDeadlineError
						})).
						Return(nil)
					netConnMock.EXPECT().
						SetWriteDeadline(mock.MatchedBy(func(deadline time.Time) bool {
							expectedDeadline := time.Now().Add(12 * time.Minute)
							return expectedDeadline.Sub(deadline).Abs() < acceptableDeadlineError
						})).
						Return(nil)
					netConnMock.EXPECT().
						Read(mock.AnythingOfType("[]uint8")).
						RunAndReturn(func(buffer []byte) (int, error) {
							return copy(buffer, "raw-response"), io.EOF
						})
					netConnMock.EXPECT().
						Write([]byte("marshalled-request")).
						RunAndReturn(func(data []byte) (int, error) {
							return len(data), nil
						})

					return netConnMock
				},
				options: func(
					test *testing.T,
				) tcpServer.TCPClientOptions[request, response] {
					clientProtocolMock :=
						tcpServerExternalMocks.NewMockClientProtocol[request, response](test)
					clientProtocolMock.EXPECT().
						InitialScannerBufferSize().
						Return(4096)
					clientProtocolMock.EXPECT().
						MaxTokenSize().
						Return(bufio.MaxScanTokenSize)
					clientProtocolMock.EXPECT().
						ExtractToken([]byte("raw-response"), true).
						Return(len("raw-response"), []byte("response"), bufio.ErrFinalToken)
					clientProtocolMock.EXPECT().
						MarshalRequest(request("request")).
						Return([]byte("marshalled-request"), nil)
					clientProtocolMock.EXPECT().
						ParseResponse([]byte("response")).
						Return("parsed-response", nil)

					return tcpServer.TCPClientOptions[request, response]{
						ReadTimeout:    mo.Some(5 * time.Minute),
						WriteTimeout:   mo.Some(12 * time.Minute),
						ClientProtocol: clientProtocolMock,
					}
				},
			},
			args: args{
				request: "request",
			},
			want:    "parsed-response",
			wantErr: assert.NoError,
		},
		{
			name: "error/unable to marshal the request",
			constructorArgs: constructorArgs{
				connection: func(test *testing.T) net.Conn {
					return tcpServerMocks.NewMocknetConn(test)
				},
				options: func(
					test *testing.T,
				) tcpServer.TCPClientOptions[request, response] {
					clientProtocolMock :=
						tcpServerExternalMocks.NewMockClientProtocol[request, response](test)
					clientProtocolMock.EXPECT().
						InitialScannerBufferSize().
						Return(4096)
					clientProtocolMock.EXPECT().
						MaxTokenSize().
						Return(bufio.MaxScanTokenSize)
					clientProtocolMock.EXPECT().
						MarshalRequest(request("request")).
						Return(nil, iotest.ErrTimeout)

					return tcpServer.TCPClientOptions[request, response]{
						ReadTimeout:    mo.Some(5 * time.Minute),
						WriteTimeout:   mo.Some(12 * time.Minute),
						ClientProtocol: clientProtocolMock,
					}
				},
			},
			args: args{
				request: "request",
			},
			want:    "",
			wantErr: assert.Error,
		},
		{
			name: "error/unable to set the write deadline",
			constructorArgs: constructorArgs{
				connection: func(test *testing.T) net.Conn {
					netConnMock := tcpServerMocks.NewMocknetConn(test)
					netConnMock.EXPECT().
						SetWriteDeadline(mock.MatchedBy(func(deadline time.Time) bool {
							expectedDeadline := time.Now().Add(12 * time.Minute)
							return expectedDeadline.Sub(deadline).Abs() < acceptableDeadlineError
						})).
						Return(iotest.ErrTimeout)

					return netConnMock
				},
				options: func(
					test *testing.T,
				) tcpServer.TCPClientOptions[request, response] {
					clientProtocolMock :=
						tcpServerExternalMocks.NewMockClientProtocol[request, response](test)
					clientProtocolMock.EXPECT().
						InitialScannerBufferSize().
						Return(4096)
					clientProtocolMock.EXPECT().
						MaxTokenSize().
						Return(bufio.MaxScanTokenSize)
					clientProtocolMock.EXPECT().
						MarshalRequest(request("request")).
						Return([]byte("marshalled-request"), nil)

					return tcpServer.TCPClientOptions[request, response]{
						ReadTimeout:    mo.Some(5 * time.Minute),
						WriteTimeout:   mo.Some(12 * time.Minute),
						ClientProtocol: clientProtocolMock,
					}
				},
			},
			args: args{
				request: "request",
			},
			want:    "",
			wantErr: assert.Error,
		},
		{
			name: "error/unable to write the request",
			constructorArgs: constructorArgs{
				connection: func(test *testing.T) net.Conn {
					netConnMock := tcpServerMocks.NewMocknetConn(test)
					netConnMock.EXPECT().
						SetWriteDeadline(mock.MatchedBy(func(deadline time.Time) bool {
							expectedDeadline := time.Now().Add(12 * time.Minute)
							return expectedDeadline.Sub(deadline).Abs() < acceptableDeadlineError
						})).
						Return(nil)
					netConnMock.EXPECT().
						Write([]byte("marshalled-request")).
						Return(0, iotest.ErrTimeout)

					return netConnMock
				},
				options: func(
					test *testing.T,
				) tcpServer.TCPClientOptions[request, response] {
					clientProtocolMock :=
						tcpServerExternalMocks.NewMockClientProtocol[request, response](test)
					clientProtocolMock.EXPECT().
						InitialScannerBufferSize().
						Return(4096)
					clientProtocolMock.EXPECT().
						MaxTokenSize().
						Return(bufio.MaxScanTokenSize)
					clientProtocolMock.EXPECT().
						MarshalRequest(request("request")).
						Return([]byte("marshalled-request"), nil)

					return tcpServer.TCPClientOptions[request, response]{
						ReadTimeout:    mo.Some(5 * time.Minute),
						WriteTimeout:   mo.Some(12 * time.Minute),
						ClientProtocol: clientProtocolMock,
					}
				},
			},
			args: args{
				request: "request",
			},
			want:    "",
			wantErr: assert.Error,
		},
		{
			name: "error/unable to set the read deadline",
			constructorArgs: constructorArgs{
				connection: func(test *testing.T) net.Conn {
					netConnMock := tcpServerMocks.NewMocknetConn(test)
					netConnMock.EXPECT().
						SetReadDeadline(mock.MatchedBy(func(deadline time.Time) bool {
							expectedDeadline := time.Now().Add(5 * time.Minute)
							return expectedDeadline.Sub(deadline).Abs() < acceptableDeadlineError
						})).
						Return(iotest.ErrTimeout)
					netConnMock.EXPECT().
						SetWriteDeadline(mock.MatchedBy(func(deadline time.Time) bool {
							expectedDeadline := time.Now().Add(12 * time.Minute)
							return expectedDeadline.Sub(deadline).Abs() < acceptableDeadlineError
						})).
						Return(nil)
					netConnMock.EXPECT().
						Write([]byte("marshalled-request")).
						RunAndReturn(func(data []byte) (int, error) {
							return len(data), nil
						})

					return netConnMock
				},
				options: func(
					test *testing.T,
				) tcpServer.TCPClientOptions[request, response] {
					clientProtocolMock :=
						tcpServerExternalMocks.NewMockClientProtocol[request, response](test)
					clientProtocolMock.EXPECT().
						InitialScannerBufferSize().
						Return(4096)
					clientProtocolMock.EXPECT().
						MaxTokenSize().
						Return(bufio.MaxScanTokenSize)
					clientProtocolMock.EXPECT().
						MarshalRequest(request("request")).
						Return([]byte("marshalled-request"), nil)

					return tcpServer.TCPClientOptions[request, response]{
						ReadTimeout:    mo.Some(5 * time.Minute),
						WriteTimeout:   mo.Some(12 * time.Minute),
						ClientProtocol: clientProtocolMock,
					}
				},
			},
			args: args{
				request: "request",
			},
			want:    "",
			wantErr: assert.Error,
		},
		{
			name: "error/unable to read the response/`net.Conn.Read()`",
			constructorArgs: constructorArgs{
				connection: func(test *testing.T) net.Conn {
					netConnMock := tcpServerMocks.NewMocknetConn(test)
					netConnMock.EXPECT().
						SetReadDeadline(mock.MatchedBy(func(deadline time.Time) bool {
							expectedDeadline := time.Now().Add(5 * time.Minute)
							return expectedDeadline.Sub(deadline).Abs() < acceptableDeadlineError
						})).
						Return(nil)
					netConnMock.EXPECT().
						SetWriteDeadline(mock.MatchedBy(func(deadline time.Time) bool {
							expectedDeadline := time.Now().Add(12 * time.Minute)
							return expectedDeadline.Sub(deadline).Abs() < acceptableDeadlineError
						})).
						Return(nil)
					netConnMock.EXPECT().
						Read(mock.AnythingOfType("[]uint8")).
						Return(0, iotest.ErrTimeout)
					netConnMock.EXPECT().
						Write([]byte("marshalled-request")).
						RunAndReturn(func(data []byte) (int, error) {
							return len(data), nil
						})

					return netConnMock
				},
				options: func(
					test *testing.T,
				) tcpServer.TCPClientOptions[request, response] {
					clientProtocolMock :=
						tcpServerExternalMocks.NewMockClientProtocol[request, response](test)
					clientProtocolMock.EXPECT().
						InitialScannerBufferSize().
						Return(4096)
					clientProtocolMock.EXPECT().
						MaxTokenSize().
						Return(bufio.MaxScanTokenSize)
					clientProtocolMock.EXPECT().
						ExtractToken([]byte{}, true).
						Return(0, nil, bufio.ErrFinalToken)
					clientProtocolMock.EXPECT().
						MarshalRequest(request("request")).
						Return([]byte("marshalled-request"), nil)

					return tcpServer.TCPClientOptions[request, response]{
						ReadTimeout:    mo.Some(5 * time.Minute),
						WriteTimeout:   mo.Some(12 * time.Minute),
						ClientProtocol: clientProtocolMock,
					}
				},
			},
			args: args{
				request: "request",
			},
			want:    "",
			wantErr: assert.Error,
		},
		{
			name: "error/unable to read the response/`Protocol.ExtractToken()`",
			constructorArgs: constructorArgs{
				connection: func(test *testing.T) net.Conn {
					netConnMock := tcpServerMocks.NewMocknetConn(test)
					netConnMock.EXPECT().
						SetReadDeadline(mock.MatchedBy(func(deadline time.Time) bool {
							expectedDeadline := time.Now().Add(5 * time.Minute)
							return expectedDeadline.Sub(deadline).Abs() < acceptableDeadlineError
						})).
						Return(nil)
					netConnMock.EXPECT().
						SetWriteDeadline(mock.MatchedBy(func(deadline time.Time) bool {
							expectedDeadline := time.Now().Add(12 * time.Minute)
							return expectedDeadline.Sub(deadline).Abs() < acceptableDeadlineError
						})).
						Return(nil)
					netConnMock.EXPECT().
						Read(mock.AnythingOfType("[]uint8")).
						RunAndReturn(func(buffer []byte) (int, error) {
							return copy(buffer, "raw-response"), io.EOF
						})
					netConnMock.EXPECT().
						Write([]byte("marshalled-request")).
						RunAndReturn(func(data []byte) (int, error) {
							return len(data), nil
						})

					return netConnMock
				},
				options: func(
					test *testing.T,
				) tcpServer.TCPClientOptions[request, response] {
					clientProtocolMock :=
						tcpServerExternalMocks.NewMockClientProtocol[request, response](test)
					clientProtocolMock.EXPECT().
						InitialScannerBufferSize().
						Return(4096)
					clientProtocolMock.EXPECT().
						MaxTokenSize().
						Return(bufio.MaxScanTokenSize)
					clientProtocolMock.EXPECT().
						ExtractToken([]byte("raw-response"), true).
						Return(0, nil, iotest.ErrTimeout)
					clientProtocolMock.EXPECT().
						MarshalRequest(request("request")).
						Return([]byte("marshalled-request"), nil)

					return tcpServer.TCPClientOptions[request, response]{
						ReadTimeout:    mo.Some(5 * time.Minute),
						WriteTimeout:   mo.Some(12 * time.Minute),
						ClientProtocol: clientProtocolMock,
					}
				},
			},
			args: args{
				request: "request",
			},
			want:    "",
			wantErr: assert.Error,
		},
		{
			name: "error/scanner has no more tokens",
			constructorArgs: constructorArgs{
				connection: func(test *testing.T) net.Conn {
					netConnMock := tcpServerMocks.NewMocknetConn(test)
					netConnMock.EXPECT().
						SetReadDeadline(mock.MatchedBy(func(deadline time.Time) bool {
							expectedDeadline := time.Now().Add(5 * time.Minute)
							return expectedDeadline.Sub(deadline).Abs() < acceptableDeadlineError
						})).
						Return(nil)
					netConnMock.EXPECT().
						SetWriteDeadline(mock.MatchedBy(func(deadline time.Time) bool {
							expectedDeadline := time.Now().Add(12 * time.Minute)
							return expectedDeadline.Sub(deadline).Abs() < acceptableDeadlineError
						})).
						Return(nil)
					netConnMock.EXPECT().
						Read(mock.AnythingOfType("[]uint8")).
						RunAndReturn(func(buffer []byte) (int, error) {
							return 0, io.EOF
						})
					netConnMock.EXPECT().
						Write([]byte("marshalled-request")).
						RunAndReturn(func(data []byte) (int, error) {
							return len(data), nil
						})

					return netConnMock
				},
				options: func(
					test *testing.T,
				) tcpServer.TCPClientOptions[request, response] {
					clientProtocolMock :=
						tcpServerExternalMocks.NewMockClientProtocol[request, response](test)
					clientProtocolMock.EXPECT().
						InitialScannerBufferSize().
						Return(4096)
					clientProtocolMock.EXPECT().
						MaxTokenSize().
						Return(bufio.MaxScanTokenSize)
					clientProtocolMock.EXPECT().
						ExtractToken([]byte{}, true).
						Return(0, nil, bufio.ErrFinalToken)
					clientProtocolMock.EXPECT().
						MarshalRequest(request("request")).
						Return([]byte("marshalled-request"), nil)

					return tcpServer.TCPClientOptions[request, response]{
						ReadTimeout:    mo.Some(5 * time.Minute),
						WriteTimeout:   mo.Some(12 * time.Minute),
						ClientProtocol: clientProtocolMock,
					}
				},
			},
			args: args{
				request: "request",
			},
			want:    "",
			wantErr: assert.Error,
		},
		{
			name: "error/unable to parse the response",
			constructorArgs: constructorArgs{
				connection: func(test *testing.T) net.Conn {
					netConnMock := tcpServerMocks.NewMocknetConn(test)
					netConnMock.EXPECT().
						SetReadDeadline(mock.MatchedBy(func(deadline time.Time) bool {
							expectedDeadline := time.Now().Add(5 * time.Minute)
							return expectedDeadline.Sub(deadline).Abs() < acceptableDeadlineError
						})).
						Return(nil)
					netConnMock.EXPECT().
						SetWriteDeadline(mock.MatchedBy(func(deadline time.Time) bool {
							expectedDeadline := time.Now().Add(12 * time.Minute)
							return expectedDeadline.Sub(deadline).Abs() < acceptableDeadlineError
						})).
						Return(nil)
					netConnMock.EXPECT().
						Read(mock.AnythingOfType("[]uint8")).
						RunAndReturn(func(buffer []byte) (int, error) {
							return copy(buffer, "raw-response"), io.EOF
						})
					netConnMock.EXPECT().
						Write([]byte("marshalled-request")).
						RunAndReturn(func(data []byte) (int, error) {
							return len(data), nil
						})

					return netConnMock
				},
				options: func(
					test *testing.T,
				) tcpServer.TCPClientOptions[request, response] {
					clientProtocolMock :=
						tcpServerExternalMocks.NewMockClientProtocol[request, response](test)
					clientProtocolMock.EXPECT().
						InitialScannerBufferSize().
						Return(4096)
					clientProtocolMock.EXPECT().
						MaxTokenSize().
						Return(bufio.MaxScanTokenSize)
					clientProtocolMock.EXPECT().
						ExtractToken([]byte("raw-response"), true).
						Return(len("raw-response"), []byte("response"), bufio.ErrFinalToken)
					clientProtocolMock.EXPECT().
						MarshalRequest(request("request")).
						Return([]byte("marshalled-request"), nil)
					clientProtocolMock.EXPECT().
						ParseResponse([]byte("response")).
						Return("", iotest.ErrTimeout)

					return tcpServer.TCPClientOptions[request, response]{
						ReadTimeout:    mo.Some(5 * time.Minute),
						WriteTimeout:   mo.Some(12 * time.Minute),
						ClientProtocol: clientProtocolMock,
					}
				},
			},
			args: args{
				request: "request",
			},
			want:    "",
			wantErr: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			client := tcpServer.NewTCPClientFromConnection(
				data.constructorArgs.connection(test),
				data.constructorArgs.options(test),
			)
			got, err := client.SendRequest(data.args.request)

			assert.Equal(test, data.want, got)
			data.wantErr(test, err)
		})
	}
}

func TestTCPClient_Close(test *testing.T) {
	type request string
	type response string
	type constructorArgs struct {
		connection func(test *testing.T) net.Conn
		options    func(test *testing.T) tcpServer.TCPClientOptions[request, response]
	}

	for _, data := range []struct {
		name            string
		constructorArgs constructorArgs
		wantErr         assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			constructorArgs: constructorArgs{
				connection: func(test *testing.T) net.Conn {
					netConnMock := tcpServerMocks.NewMocknetConn(test)
					netConnMock.EXPECT().
						Close().
						Return(nil)

					return netConnMock
				},
				options: func(
					test *testing.T,
				) tcpServer.TCPClientOptions[request, response] {
					clientProtocolMock :=
						tcpServerExternalMocks.NewMockClientProtocol[request, response](test)
					clientProtocolMock.EXPECT().
						InitialScannerBufferSize().
						Return(4096)
					clientProtocolMock.EXPECT().
						MaxTokenSize().
						Return(bufio.MaxScanTokenSize)

					return tcpServer.TCPClientOptions[request, response]{
						ReadTimeout:    mo.Some(5 * time.Minute),
						WriteTimeout:   mo.Some(12 * time.Minute),
						ClientProtocol: clientProtocolMock,
					}
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "error",
			constructorArgs: constructorArgs{
				connection: func(test *testing.T) net.Conn {
					netConnMock := tcpServerMocks.NewMocknetConn(test)
					netConnMock.EXPECT().
						Close().
						Return(iotest.ErrTimeout)

					return netConnMock
				},
				options: func(
					test *testing.T,
				) tcpServer.TCPClientOptions[request, response] {
					clientProtocolMock :=
						tcpServerExternalMocks.NewMockClientProtocol[request, response](test)
					clientProtocolMock.EXPECT().
						InitialScannerBufferSize().
						Return(4096)
					clientProtocolMock.EXPECT().
						MaxTokenSize().
						Return(bufio.MaxScanTokenSize)

					return tcpServer.TCPClientOptions[request, response]{
						ReadTimeout:    mo.Some(5 * time.Minute),
						WriteTimeout:   mo.Some(12 * time.Minute),
						ClientProtocol: clientProtocolMock,
					}
				},
			},
			wantErr: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			client := tcpServer.NewTCPClientFromConnection(
				data.constructorArgs.connection(test),
				data.constructorArgs.options(test),
			)
			err := client.Close()

			data.wantErr(test, err)
		})
	}
}

func runTestServer(test *testing.T, address string) net.Listener {
	listener, err := net.Listen(tcpServer.TCPServerNetwork, address)
	require.NoError(test, err)

	go func() {
		for {
			connection, err := listener.Accept()
			if err != nil {
				test.Logf("unable to accept the connection: %s", err)
				break
			}

			go func() {
				defer connection.Close()

				_, err := io.Copy(connection, connection)
				test.Logf("unable to echo all incoming data: %s", err)
			}()
		}
	}()

	return listener
}
