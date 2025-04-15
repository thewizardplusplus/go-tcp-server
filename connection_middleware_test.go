package tcpServer

import (
	"context"
	"net"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	tcpServerMocks "github.com/thewizardplusplus/go-tcp-server/mocks/github.com/thewizardplusplus/go-tcp-server"
)

func TestApplyConnectionMiddlewares(test *testing.T) {
	const maxNumberCount = 10

	type args struct {
		handler     func(test *testing.T) ConnectionHandler
		middlewares func(test *testing.T, numbers chan int) []ConnectionMiddleware
	}
	type handlerArgs struct {
		ctx        context.Context
		connection func(test *testing.T) net.Conn
	}

	for _, data := range []struct {
		name            string
		args            args
		handlerArgs     handlerArgs
		wantNumberSlice []int
		wantHandlerErr  assert.ErrorAssertionFunc
	}{
		{
			name: "success/without middlewares",
			args: args{
				handler: func(test *testing.T) ConnectionHandler {
					connectionHandlerMock := tcpServerMocks.NewMockConnectionHandler(test)
					connectionHandlerMock.EXPECT().
						HandleConnection(
							context.Background(),
							tcpServerMocks.NewMocknetConn(test),
						).
						Return(nil)

					return connectionHandlerMock
				},
				middlewares: func(
					test *testing.T,
					numbers chan int,
				) []ConnectionMiddleware {
					return nil
				},
			},
			handlerArgs: handlerArgs{
				ctx: context.Background(),
				connection: func(test *testing.T) net.Conn {
					return tcpServerMocks.NewMocknetConn(test)
				},
			},
			wantNumberSlice: []int{},
			wantHandlerErr:  assert.NoError,
		},
		{
			name: "success/with middlewares",
			args: args{
				handler: func(test *testing.T) ConnectionHandler {
					connectionHandlerMock := tcpServerMocks.NewMockConnectionHandler(test)
					connectionHandlerMock.EXPECT().
						HandleConnection(
							context.Background(),
							tcpServerMocks.NewMocknetConn(test),
						).
						Return(nil)

					return connectionHandlerMock
				},
				middlewares: func(
					test *testing.T,
					numbers chan int,
				) []ConnectionMiddleware {
					return []ConnectionMiddleware{
						func(handler ConnectionHandler) ConnectionHandler {
							return ConnectionHandlerFunc(func(
								ctx context.Context,
								connection net.Conn,
							) error {
								numbers <- 23
								return handler.HandleConnection(ctx, connection)
							})
						},
						func(handler ConnectionHandler) ConnectionHandler {
							return ConnectionHandlerFunc(func(
								ctx context.Context,
								connection net.Conn,
							) error {
								numbers <- 42
								return handler.HandleConnection(ctx, connection)
							})
						},
					}
				},
			},
			handlerArgs: handlerArgs{
				ctx: context.Background(),
				connection: func(test *testing.T) net.Conn {
					return tcpServerMocks.NewMocknetConn(test)
				},
			},
			wantNumberSlice: []int{42, 23},
			wantHandlerErr:  assert.NoError,
		},
		{
			name: "error/with middlewares",
			args: args{
				handler: func(test *testing.T) ConnectionHandler {
					connectionHandlerMock := tcpServerMocks.NewMockConnectionHandler(test)
					connectionHandlerMock.EXPECT().
						HandleConnection(
							context.Background(),
							tcpServerMocks.NewMocknetConn(test),
						).
						Return(iotest.ErrTimeout)

					return connectionHandlerMock
				},
				middlewares: func(
					test *testing.T,
					numbers chan int,
				) []ConnectionMiddleware {
					return []ConnectionMiddleware{
						func(handler ConnectionHandler) ConnectionHandler {
							return ConnectionHandlerFunc(func(
								ctx context.Context,
								connection net.Conn,
							) error {
								numbers <- 23
								return handler.HandleConnection(ctx, connection)
							})
						},
						func(handler ConnectionHandler) ConnectionHandler {
							return ConnectionHandlerFunc(func(
								ctx context.Context,
								connection net.Conn,
							) error {
								numbers <- 42
								return handler.HandleConnection(ctx, connection)
							})
						},
					}
				},
			},
			handlerArgs: handlerArgs{
				ctx: context.Background(),
				connection: func(test *testing.T) net.Conn {
					return tcpServerMocks.NewMocknetConn(test)
				},
			},
			wantNumberSlice: []int{42, 23},
			wantHandlerErr: func(
				test assert.TestingT,
				err error,
				msgAndArgs ...any,
			) bool {
				return assert.ErrorIs(test, err, iotest.ErrTimeout)
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			numbers := make(chan int, maxNumberCount)
			got := ApplyConnectionMiddlewares(
				data.args.handler(test),
				data.args.middlewares(test, numbers),
			)

			handlerErr := got.HandleConnection(
				data.handlerArgs.ctx,
				data.handlerArgs.connection(test),
			)

			close(numbers)

			gotNumberSlice := make([]int, 0, maxNumberCount)
			for number := range numbers {
				gotNumberSlice = append(gotNumberSlice, number)
			}

			assert.Equal(test, data.wantNumberSlice, gotNumberSlice)
			data.wantHandlerErr(test, handlerErr)
		})
	}
}
