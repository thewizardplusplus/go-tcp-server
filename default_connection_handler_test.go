package tcpServer_test

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"strings"
	"testing"
	"testing/iotest"
	"time"

	"github.com/samber/mo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	tcpServer "github.com/thewizardplusplus/go-tcp-server"
	tcpServerExternalMocks "github.com/thewizardplusplus/go-tcp-server/mocks/external/github.com/thewizardplusplus/go-tcp-server"
	tcpServerMocks "github.com/thewizardplusplus/go-tcp-server/mocks/github.com/thewizardplusplus/go-tcp-server"
)

func TestDefaultConnectionHandler_interface(test *testing.T) {
	type request string
	type response string

	assert.Implements(
		test,
		(*tcpServer.ConnectionHandler)(nil),
		tcpServer.DefaultConnectionHandler[request, response]{},
	)
}

func TestDefaultConnectionHandler_HandleRequest(test *testing.T) {
	const acceptableDeadlineError = time.Minute

	type request string
	type response string
	type constructorArgs struct {
		options func(test *testing.T) tcpServer.DefaultConnectionHandlerOptions[
			request,
			response,
		]
	}
	type args struct {
		ctx        context.Context
		connection func(test *testing.T) net.Conn
		scanner    *bufio.Scanner
	}

	for _, data := range []struct {
		name            string
		constructorArgs constructorArgs
		args            args
		wantErr         assert.ErrorAssertionFunc
	}{
		{
			name: "success/without all the timeouts",
			constructorArgs: constructorArgs{
				options: func(test *testing.T) tcpServer.DefaultConnectionHandlerOptions[
					request,
					response,
				] {
					serverProtocolMock :=
						tcpServerExternalMocks.NewMockServerProtocol[request, response](test)
					serverProtocolMock.EXPECT().
						ParseRequest([]byte("request")).
						Return("parsed-request", nil)
					serverProtocolMock.EXPECT().
						MarshalResponse(response("response")).
						Return([]byte("marshalled-response"), nil)

					requestHandlerMock :=
						tcpServerExternalMocks.NewMockRequestHandler[request, response](test)
					requestHandlerMock.EXPECT().
						HandleRequest(
							mock.AnythingOfType("context.withoutCancelCtx"),
							request("parsed-request"),
						).
						Return("response", nil)

					return tcpServer.DefaultConnectionHandlerOptions[request, response]{
						ServerProtocol: serverProtocolMock,
						RequestHandler: requestHandlerMock,
					}
				},
			},
			args: args{
				ctx: context.Background(),
				connection: func(test *testing.T) net.Conn {
					netConnMock := tcpServerMocks.NewMocknetConn(test)
					netConnMock.EXPECT().
						Write([]byte("marshalled-response")).
						RunAndReturn(func(data []byte) (int, error) {
							return len(data), nil
						})

					return netConnMock
				},
				scanner: bufio.NewScanner(strings.NewReader("request")),
			},
			wantErr: assert.NoError,
		},
		{
			name: "success/with all the timeouts",
			constructorArgs: constructorArgs{
				options: func(test *testing.T) tcpServer.DefaultConnectionHandlerOptions[
					request,
					response,
				] {
					serverProtocolMock :=
						tcpServerExternalMocks.NewMockServerProtocol[request, response](test)
					serverProtocolMock.EXPECT().
						ParseRequest([]byte("request")).
						Return("parsed-request", nil)
					serverProtocolMock.EXPECT().
						MarshalResponse(response("response")).
						Return([]byte("marshalled-response"), nil)

					requestHandlerMock :=
						tcpServerExternalMocks.NewMockRequestHandler[request, response](test)
					requestHandlerMock.EXPECT().
						HandleRequest(
							mock.AnythingOfType("*context.timerCtx"),
							request("parsed-request"),
						).
						Return("response", nil)

					return tcpServer.DefaultConnectionHandlerOptions[request, response]{
						ReadTimeout:     mo.Some(5 * time.Minute),
						WriteTimeout:    mo.Some(12 * time.Minute),
						HandlingTimeout: mo.Some(23 * time.Minute),
						ServerProtocol:  serverProtocolMock,
						RequestHandler:  requestHandlerMock,
					}
				},
			},
			args: args{
				ctx: context.Background(),
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
						Write([]byte("marshalled-response")).
						RunAndReturn(func(data []byte) (int, error) {
							return len(data), nil
						})

					return netConnMock
				},
				scanner: bufio.NewScanner(strings.NewReader("request")),
			},
			wantErr: assert.NoError,
		},
		{
			name: "error/unable to set the read deadline",
			constructorArgs: constructorArgs{
				options: func(test *testing.T) tcpServer.DefaultConnectionHandlerOptions[
					request,
					response,
				] {
					serverProtocolMock :=
						tcpServerExternalMocks.NewMockServerProtocol[request, response](test)
					requestHandlerMock :=
						tcpServerExternalMocks.NewMockRequestHandler[request, response](test)
					return tcpServer.DefaultConnectionHandlerOptions[request, response]{
						ReadTimeout:    mo.Some(5 * time.Minute),
						ServerProtocol: serverProtocolMock,
						RequestHandler: requestHandlerMock,
					}
				},
			},
			args: args{
				ctx: context.Background(),
				connection: func(test *testing.T) net.Conn {
					netConnMock := tcpServerMocks.NewMocknetConn(test)
					netConnMock.EXPECT().
						SetReadDeadline(mock.MatchedBy(func(deadline time.Time) bool {
							expectedDeadline := time.Now().Add(5 * time.Minute)
							return expectedDeadline.Sub(deadline).Abs() < acceptableDeadlineError
						})).
						Return(iotest.ErrTimeout)

					return netConnMock
				},
				scanner: bufio.NewScanner(strings.NewReader("request")),
			},
			wantErr: assert.Error,
		},
		{
			name: "error/unable to read the request",
			constructorArgs: constructorArgs{
				options: func(test *testing.T) tcpServer.DefaultConnectionHandlerOptions[
					request,
					response,
				] {
					serverProtocolMock :=
						tcpServerExternalMocks.NewMockServerProtocol[request, response](test)
					requestHandlerMock :=
						tcpServerExternalMocks.NewMockRequestHandler[request, response](test)
					return tcpServer.DefaultConnectionHandlerOptions[request, response]{
						ReadTimeout:    mo.Some(5 * time.Minute),
						ServerProtocol: serverProtocolMock,
						RequestHandler: requestHandlerMock,
					}
				},
			},
			args: args{
				ctx: context.Background(),
				connection: func(test *testing.T) net.Conn {
					netConnMock := tcpServerMocks.NewMocknetConn(test)
					netConnMock.EXPECT().
						SetReadDeadline(mock.MatchedBy(func(deadline time.Time) bool {
							expectedDeadline := time.Now().Add(5 * time.Minute)
							return expectedDeadline.Sub(deadline).Abs() < acceptableDeadlineError
						})).
						Return(nil)

					return netConnMock
				},
				scanner: bufio.NewScanner(iotest.ErrReader(iotest.ErrTimeout)),
			},
			wantErr: assert.Error,
		},
		{
			name: "error/scanner has no more tokens",
			constructorArgs: constructorArgs{
				options: func(test *testing.T) tcpServer.DefaultConnectionHandlerOptions[
					request,
					response,
				] {
					serverProtocolMock :=
						tcpServerExternalMocks.NewMockServerProtocol[request, response](test)
					requestHandlerMock :=
						tcpServerExternalMocks.NewMockRequestHandler[request, response](test)
					return tcpServer.DefaultConnectionHandlerOptions[request, response]{
						ReadTimeout:    mo.Some(5 * time.Minute),
						ServerProtocol: serverProtocolMock,
						RequestHandler: requestHandlerMock,
					}
				},
			},
			args: args{
				ctx: context.Background(),
				connection: func(test *testing.T) net.Conn {
					netConnMock := tcpServerMocks.NewMocknetConn(test)
					netConnMock.EXPECT().
						SetReadDeadline(mock.MatchedBy(func(deadline time.Time) bool {
							expectedDeadline := time.Now().Add(5 * time.Minute)
							return expectedDeadline.Sub(deadline).Abs() < acceptableDeadlineError
						})).
						Return(nil)

					return netConnMock
				},
				scanner: bufio.NewScanner(strings.NewReader("")),
			},
			wantErr: func(test assert.TestingT, err error, msgAndArgs ...any) bool {
				return assert.ErrorIs(test, err, tcpServer.ErrHandlingStopIsRequired)
			},
		},
		{
			name: "error/unable to parse the request",
			constructorArgs: constructorArgs{
				options: func(test *testing.T) tcpServer.DefaultConnectionHandlerOptions[
					request,
					response,
				] {
					serverProtocolMock :=
						tcpServerExternalMocks.NewMockServerProtocol[request, response](test)
					serverProtocolMock.EXPECT().
						ParseRequest([]byte("request")).
						Return("", iotest.ErrTimeout)

					requestHandlerMock :=
						tcpServerExternalMocks.NewMockRequestHandler[request, response](test)
					return tcpServer.DefaultConnectionHandlerOptions[request, response]{
						ReadTimeout:    mo.Some(5 * time.Minute),
						ServerProtocol: serverProtocolMock,
						RequestHandler: requestHandlerMock,
					}
				},
			},
			args: args{
				ctx: context.Background(),
				connection: func(test *testing.T) net.Conn {
					netConnMock := tcpServerMocks.NewMocknetConn(test)
					netConnMock.EXPECT().
						SetReadDeadline(mock.MatchedBy(func(deadline time.Time) bool {
							expectedDeadline := time.Now().Add(5 * time.Minute)
							return expectedDeadline.Sub(deadline).Abs() < acceptableDeadlineError
						})).
						Return(nil)

					return netConnMock
				},
				scanner: bufio.NewScanner(strings.NewReader("request")),
			},
			wantErr: assert.Error,
		},
		{
			name: "error/unable to handle the request",
			constructorArgs: constructorArgs{
				options: func(test *testing.T) tcpServer.DefaultConnectionHandlerOptions[
					request,
					response,
				] {
					serverProtocolMock :=
						tcpServerExternalMocks.NewMockServerProtocol[request, response](test)
					serverProtocolMock.EXPECT().
						ParseRequest([]byte("request")).
						Return("parsed-request", nil)

					requestHandlerMock :=
						tcpServerExternalMocks.NewMockRequestHandler[request, response](test)
					requestHandlerMock.EXPECT().
						HandleRequest(
							mock.AnythingOfType("context.withoutCancelCtx"),
							request("parsed-request"),
						).
						Return("", iotest.ErrTimeout)

					return tcpServer.DefaultConnectionHandlerOptions[request, response]{
						ReadTimeout:    mo.Some(5 * time.Minute),
						ServerProtocol: serverProtocolMock,
						RequestHandler: requestHandlerMock,
					}
				},
			},
			args: args{
				ctx: context.Background(),
				connection: func(test *testing.T) net.Conn {
					netConnMock := tcpServerMocks.NewMocknetConn(test)
					netConnMock.EXPECT().
						SetReadDeadline(mock.MatchedBy(func(deadline time.Time) bool {
							expectedDeadline := time.Now().Add(5 * time.Minute)
							return expectedDeadline.Sub(deadline).Abs() < acceptableDeadlineError
						})).
						Return(nil)

					return netConnMock
				},
				scanner: bufio.NewScanner(strings.NewReader("request")),
			},
			wantErr: assert.Error,
		},
		{
			name: "error/unable to marshal the response",
			constructorArgs: constructorArgs{
				options: func(test *testing.T) tcpServer.DefaultConnectionHandlerOptions[
					request,
					response,
				] {
					serverProtocolMock :=
						tcpServerExternalMocks.NewMockServerProtocol[request, response](test)
					serverProtocolMock.EXPECT().
						ParseRequest([]byte("request")).
						Return("parsed-request", nil)
					serverProtocolMock.EXPECT().
						MarshalResponse(response("response")).
						Return(nil, iotest.ErrTimeout)

					requestHandlerMock :=
						tcpServerExternalMocks.NewMockRequestHandler[request, response](test)
					requestHandlerMock.EXPECT().
						HandleRequest(
							mock.AnythingOfType("context.withoutCancelCtx"),
							request("parsed-request"),
						).
						Return("response", nil)

					return tcpServer.DefaultConnectionHandlerOptions[request, response]{
						ReadTimeout:    mo.Some(5 * time.Minute),
						ServerProtocol: serverProtocolMock,
						RequestHandler: requestHandlerMock,
					}
				},
			},
			args: args{
				ctx: context.Background(),
				connection: func(test *testing.T) net.Conn {
					netConnMock := tcpServerMocks.NewMocknetConn(test)
					netConnMock.EXPECT().
						SetReadDeadline(mock.MatchedBy(func(deadline time.Time) bool {
							expectedDeadline := time.Now().Add(5 * time.Minute)
							return expectedDeadline.Sub(deadline).Abs() < acceptableDeadlineError
						})).
						Return(nil)

					return netConnMock
				},
				scanner: bufio.NewScanner(strings.NewReader("request")),
			},
			wantErr: assert.Error,
		},
		{
			name: "error/unable to set the write deadline",
			constructorArgs: constructorArgs{
				options: func(test *testing.T) tcpServer.DefaultConnectionHandlerOptions[
					request,
					response,
				] {
					serverProtocolMock :=
						tcpServerExternalMocks.NewMockServerProtocol[request, response](test)
					serverProtocolMock.EXPECT().
						ParseRequest([]byte("request")).
						Return("parsed-request", nil)
					serverProtocolMock.EXPECT().
						MarshalResponse(response("response")).
						Return([]byte("marshalled-response"), nil)

					requestHandlerMock :=
						tcpServerExternalMocks.NewMockRequestHandler[request, response](test)
					requestHandlerMock.EXPECT().
						HandleRequest(
							mock.AnythingOfType("context.withoutCancelCtx"),
							request("parsed-request"),
						).
						Return("response", nil)

					return tcpServer.DefaultConnectionHandlerOptions[request, response]{
						ReadTimeout:    mo.Some(5 * time.Minute),
						WriteTimeout:   mo.Some(12 * time.Minute),
						ServerProtocol: serverProtocolMock,
						RequestHandler: requestHandlerMock,
					}
				},
			},
			args: args{
				ctx: context.Background(),
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
						Return(iotest.ErrTimeout)

					return netConnMock
				},
				scanner: bufio.NewScanner(strings.NewReader("request")),
			},
			wantErr: assert.Error,
		},
		{
			name: "error/unable to write the response",
			constructorArgs: constructorArgs{
				options: func(test *testing.T) tcpServer.DefaultConnectionHandlerOptions[
					request,
					response,
				] {
					serverProtocolMock :=
						tcpServerExternalMocks.NewMockServerProtocol[request, response](test)
					serverProtocolMock.EXPECT().
						ParseRequest([]byte("request")).
						Return("parsed-request", nil)
					serverProtocolMock.EXPECT().
						MarshalResponse(response("response")).
						Return([]byte("marshalled-response"), nil)

					requestHandlerMock :=
						tcpServerExternalMocks.NewMockRequestHandler[request, response](test)
					requestHandlerMock.EXPECT().
						HandleRequest(
							mock.AnythingOfType("context.withoutCancelCtx"),
							request("parsed-request"),
						).
						Return("response", nil)

					return tcpServer.DefaultConnectionHandlerOptions[request, response]{
						ReadTimeout:    mo.Some(5 * time.Minute),
						WriteTimeout:   mo.Some(12 * time.Minute),
						ServerProtocol: serverProtocolMock,
						RequestHandler: requestHandlerMock,
					}
				},
			},
			args: args{
				ctx: context.Background(),
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
						Write([]byte("marshalled-response")).
						Return(0, iotest.ErrTimeout)

					return netConnMock
				},
				scanner: bufio.NewScanner(strings.NewReader("request")),
			},
			wantErr: assert.Error,
		},
		{
			name: "error/request handler requested to stop handling",
			constructorArgs: constructorArgs{
				options: func(test *testing.T) tcpServer.DefaultConnectionHandlerOptions[
					request,
					response,
				] {
					serverProtocolMock :=
						tcpServerExternalMocks.NewMockServerProtocol[request, response](test)
					serverProtocolMock.EXPECT().
						ParseRequest([]byte("request")).
						Return("parsed-request", nil)
					serverProtocolMock.EXPECT().
						MarshalResponse(response("response")).
						Return([]byte("marshalled-response"), nil)

					requestHandlerMock :=
						tcpServerExternalMocks.NewMockRequestHandler[request, response](test)
					requestHandlerMock.EXPECT().
						HandleRequest(
							mock.AnythingOfType("context.withoutCancelCtx"),
							request("parsed-request"),
						).
						Return(
							"response",
							fmt.Errorf("wrapped error: %w", tcpServer.ErrHandlingStopIsRequired),
						)

					return tcpServer.DefaultConnectionHandlerOptions[request, response]{
						ReadTimeout:    mo.Some(5 * time.Minute),
						WriteTimeout:   mo.Some(12 * time.Minute),
						ServerProtocol: serverProtocolMock,
						RequestHandler: requestHandlerMock,
					}
				},
			},
			args: args{
				ctx: context.Background(),
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
						Write([]byte("marshalled-response")).
						RunAndReturn(func(data []byte) (int, error) {
							return len(data), nil
						})

					return netConnMock
				},
				scanner: bufio.NewScanner(strings.NewReader("request")),
			},
			wantErr: func(test assert.TestingT, err error, msgAndArgs ...any) bool {
				return assert.ErrorIs(test, err, tcpServer.ErrHandlingStopIsRequired)
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			handler := tcpServer.NewDefaultConnectionHandler(
				data.constructorArgs.options(test),
			)
			err := handler.HandleRequest(
				data.args.ctx,
				data.args.connection(test),
				data.args.scanner,
			)

			data.wantErr(test, err)
		})
	}
}

func TestDefaultConnectionHandler_HandleConnection(test *testing.T) {
	const acceptableDeadlineError = time.Minute

	type request string
	type response string
	type constructorArgs struct {
		options func(test *testing.T) tcpServer.DefaultConnectionHandlerOptions[
			request,
			response,
		]
	}
	type args struct {
		ctx        context.Context
		connection func(test *testing.T) net.Conn
	}

	for _, data := range []struct {
		name            string
		constructorArgs constructorArgs
		args            args
		wantErr         assert.ErrorAssertionFunc
	}{
		{
			name: "success/scanner has no more tokens",
			constructorArgs: constructorArgs{
				options: func(test *testing.T) tcpServer.DefaultConnectionHandlerOptions[
					request,
					response,
				] {
					serverProtocolMock :=
						tcpServerExternalMocks.NewMockServerProtocol[request, response](test)
					serverProtocolMock.EXPECT().
						InitialScannerBufferSize().
						Return(4096)
					serverProtocolMock.EXPECT().
						MaxTokenSize().
						Return(bufio.MaxScanTokenSize)
					serverProtocolMock.EXPECT().
						ExtractToken([]byte("raw-request"), true).
						Return(len("raw-request"), []byte("request"), bufio.ErrFinalToken)
					serverProtocolMock.EXPECT().
						ParseRequest([]byte("request")).
						Return("parsed-request", nil)
					serverProtocolMock.EXPECT().
						MarshalResponse(response("response")).
						Return([]byte("marshalled-response"), nil)

					requestHandlerMock :=
						tcpServerExternalMocks.NewMockRequestHandler[request, response](test)
					requestHandlerMock.EXPECT().
						HandleRequest(
							mock.AnythingOfType("*context.timerCtx"),
							request("parsed-request"),
						).
						Return("response", nil)

					return tcpServer.DefaultConnectionHandlerOptions[request, response]{
						ReadTimeout:     mo.Some(5 * time.Minute),
						WriteTimeout:    mo.Some(12 * time.Minute),
						HandlingTimeout: mo.Some(23 * time.Minute),
						ServerProtocol:  serverProtocolMock,
						RequestHandler:  requestHandlerMock,
					}
				},
			},
			args: args{
				ctx: context.Background(),
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
							return copy(buffer, "raw-request"), io.EOF
						})
					netConnMock.EXPECT().
						Write([]byte("marshalled-response")).
						RunAndReturn(func(data []byte) (int, error) {
							return len(data), nil
						})

					return netConnMock
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "success/request handler requested to stop handling",
			constructorArgs: constructorArgs{
				options: func(test *testing.T) tcpServer.DefaultConnectionHandlerOptions[
					request,
					response,
				] {
					serverProtocolMock :=
						tcpServerExternalMocks.NewMockServerProtocol[request, response](test)
					serverProtocolMock.EXPECT().
						InitialScannerBufferSize().
						Return(4096)
					serverProtocolMock.EXPECT().
						MaxTokenSize().
						Return(bufio.MaxScanTokenSize)
					serverProtocolMock.EXPECT().
						ExtractToken([]byte("raw-request"), true).
						Return(len("raw-request"), []byte("request"), bufio.ErrFinalToken)
					serverProtocolMock.EXPECT().
						ParseRequest([]byte("request")).
						Return("parsed-request", nil)
					serverProtocolMock.EXPECT().
						MarshalResponse(response("response")).
						Return([]byte("marshalled-response"), nil)

					requestHandlerMock :=
						tcpServerExternalMocks.NewMockRequestHandler[request, response](test)
					requestHandlerMock.EXPECT().
						HandleRequest(
							mock.AnythingOfType("*context.timerCtx"),
							request("parsed-request"),
						).
						Return(
							"response",
							fmt.Errorf("wrapped error: %w", tcpServer.ErrHandlingStopIsRequired),
						)

					return tcpServer.DefaultConnectionHandlerOptions[request, response]{
						ReadTimeout:     mo.Some(5 * time.Minute),
						WriteTimeout:    mo.Some(12 * time.Minute),
						HandlingTimeout: mo.Some(23 * time.Minute),
						ServerProtocol:  serverProtocolMock,
						RequestHandler:  requestHandlerMock,
					}
				},
			},
			args: args{
				ctx: context.Background(),
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
							return copy(buffer, "raw-request"), io.EOF
						})
					netConnMock.EXPECT().
						Write([]byte("marshalled-response")).
						RunAndReturn(func(data []byte) (int, error) {
							return len(data), nil
						})

					return netConnMock
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "error/context is done",
			constructorArgs: constructorArgs{
				options: func(test *testing.T) tcpServer.DefaultConnectionHandlerOptions[
					request,
					response,
				] {
					serverProtocolMock :=
						tcpServerExternalMocks.NewMockServerProtocol[request, response](test)
					serverProtocolMock.EXPECT().
						InitialScannerBufferSize().
						Return(4096)
					serverProtocolMock.EXPECT().
						MaxTokenSize().
						Return(bufio.MaxScanTokenSize)

					requestHandlerMock :=
						tcpServerExternalMocks.NewMockRequestHandler[request, response](test)
					return tcpServer.DefaultConnectionHandlerOptions[request, response]{
						ReadTimeout:     mo.Some(5 * time.Minute),
						WriteTimeout:    mo.Some(12 * time.Minute),
						HandlingTimeout: mo.Some(23 * time.Minute),
						ServerProtocol:  serverProtocolMock,
						RequestHandler:  requestHandlerMock,
					}
				},
			},
			args: args{
				ctx: func() context.Context {
					ctx, ctxCancel := context.WithCancel(context.Background())
					ctxCancel()

					return ctx
				}(),
				connection: func(test *testing.T) net.Conn {
					return tcpServerMocks.NewMocknetConn(test)
				},
			},
			wantErr: assert.Error,
		},
		{
			name: "error/unable to handle the request/regular error",
			constructorArgs: constructorArgs{
				options: func(test *testing.T) tcpServer.DefaultConnectionHandlerOptions[
					request,
					response,
				] {
					serverProtocolMock :=
						tcpServerExternalMocks.NewMockServerProtocol[request, response](test)
					serverProtocolMock.EXPECT().
						InitialScannerBufferSize().
						Return(4096)
					serverProtocolMock.EXPECT().
						MaxTokenSize().
						Return(bufio.MaxScanTokenSize)
					serverProtocolMock.EXPECT().
						ExtractToken([]byte("raw-request"), true).
						Return(len("raw-request"), []byte("request"), bufio.ErrFinalToken)
					serverProtocolMock.EXPECT().
						ParseRequest([]byte("request")).
						Return("parsed-request", nil)

					requestHandlerMock :=
						tcpServerExternalMocks.NewMockRequestHandler[request, response](test)
					requestHandlerMock.EXPECT().
						HandleRequest(
							mock.AnythingOfType("*context.timerCtx"),
							request("parsed-request"),
						).
						Return("", iotest.ErrTimeout)

					return tcpServer.DefaultConnectionHandlerOptions[request, response]{
						ReadTimeout:     mo.Some(5 * time.Minute),
						WriteTimeout:    mo.Some(12 * time.Minute),
						HandlingTimeout: mo.Some(23 * time.Minute),
						ServerProtocol:  serverProtocolMock,
						RequestHandler:  requestHandlerMock,
					}
				},
			},
			args: args{
				ctx: context.Background(),
				connection: func(test *testing.T) net.Conn {
					netConnMock := tcpServerMocks.NewMocknetConn(test)
					netConnMock.EXPECT().
						SetReadDeadline(mock.MatchedBy(func(deadline time.Time) bool {
							expectedDeadline := time.Now().Add(5 * time.Minute)
							return expectedDeadline.Sub(deadline).Abs() < acceptableDeadlineError
						})).
						Return(nil)
					netConnMock.EXPECT().
						Read(mock.AnythingOfType("[]uint8")).
						RunAndReturn(func(buffer []byte) (int, error) {
							return copy(buffer, "raw-request"), io.EOF
						})

					return netConnMock
				},
			},
			wantErr: assert.Error,
		},
		{
			name: "error/unable to handle the request/`net.Conn.Read()`",
			constructorArgs: constructorArgs{
				options: func(test *testing.T) tcpServer.DefaultConnectionHandlerOptions[
					request,
					response,
				] {
					serverProtocolMock :=
						tcpServerExternalMocks.NewMockServerProtocol[request, response](test)
					serverProtocolMock.EXPECT().
						InitialScannerBufferSize().
						Return(4096)
					serverProtocolMock.EXPECT().
						MaxTokenSize().
						Return(bufio.MaxScanTokenSize)
					serverProtocolMock.EXPECT().
						ExtractToken([]byte{}, true).
						Return(0, nil, bufio.ErrFinalToken)

					requestHandlerMock :=
						tcpServerExternalMocks.NewMockRequestHandler[request, response](test)
					return tcpServer.DefaultConnectionHandlerOptions[request, response]{
						ReadTimeout:     mo.Some(5 * time.Minute),
						WriteTimeout:    mo.Some(12 * time.Minute),
						HandlingTimeout: mo.Some(23 * time.Minute),
						ServerProtocol:  serverProtocolMock,
						RequestHandler:  requestHandlerMock,
					}
				},
			},
			args: args{
				ctx: context.Background(),
				connection: func(test *testing.T) net.Conn {
					netConnMock := tcpServerMocks.NewMocknetConn(test)
					netConnMock.EXPECT().
						SetReadDeadline(mock.MatchedBy(func(deadline time.Time) bool {
							expectedDeadline := time.Now().Add(5 * time.Minute)
							return expectedDeadline.Sub(deadline).Abs() < acceptableDeadlineError
						})).
						Return(nil)
					netConnMock.EXPECT().
						Read(mock.AnythingOfType("[]uint8")).
						Return(0, iotest.ErrTimeout)

					return netConnMock
				},
			},
			wantErr: assert.Error,
		},
		{
			name: "error/unable to handle the request/`Protocol.ExtractToken()`",
			constructorArgs: constructorArgs{
				options: func(test *testing.T) tcpServer.DefaultConnectionHandlerOptions[
					request,
					response,
				] {
					serverProtocolMock :=
						tcpServerExternalMocks.NewMockServerProtocol[request, response](test)
					serverProtocolMock.EXPECT().
						InitialScannerBufferSize().
						Return(4096)
					serverProtocolMock.EXPECT().
						MaxTokenSize().
						Return(bufio.MaxScanTokenSize)
					serverProtocolMock.EXPECT().
						ExtractToken([]byte("raw-request"), true).
						Return(0, nil, iotest.ErrTimeout)

					requestHandlerMock :=
						tcpServerExternalMocks.NewMockRequestHandler[request, response](test)
					return tcpServer.DefaultConnectionHandlerOptions[request, response]{
						ReadTimeout:     mo.Some(5 * time.Minute),
						WriteTimeout:    mo.Some(12 * time.Minute),
						HandlingTimeout: mo.Some(23 * time.Minute),
						ServerProtocol:  serverProtocolMock,
						RequestHandler:  requestHandlerMock,
					}
				},
			},
			args: args{
				ctx: context.Background(),
				connection: func(test *testing.T) net.Conn {
					netConnMock := tcpServerMocks.NewMocknetConn(test)
					netConnMock.EXPECT().
						SetReadDeadline(mock.MatchedBy(func(deadline time.Time) bool {
							expectedDeadline := time.Now().Add(5 * time.Minute)
							return expectedDeadline.Sub(deadline).Abs() < acceptableDeadlineError
						})).
						Return(nil)
					netConnMock.EXPECT().
						Read(mock.AnythingOfType("[]uint8")).
						RunAndReturn(func(buffer []byte) (int, error) {
							return copy(buffer, "raw-request"), io.EOF
						})

					return netConnMock
				},
			},
			wantErr: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			handler := tcpServer.NewDefaultConnectionHandler(
				data.constructorArgs.options(test),
			)
			err := handler.HandleConnection(data.args.ctx, data.args.connection(test))

			data.wantErr(test, err)
		})
	}
}
