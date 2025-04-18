package tcpServer_test

import (
	"context"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	tcpServer "github.com/thewizardplusplus/go-tcp-server"
	tcpServerExternalMocks "github.com/thewizardplusplus/go-tcp-server/mocks/external/github.com/thewizardplusplus/go-tcp-server"
)

func TestApplyRequestMiddlewares(test *testing.T) {
	const maxNumberCount = 10

	type request string
	type response string
	type args struct {
		handler     func(test *testing.T) tcpServer.RequestHandler[request, response]
		middlewares func(
			test *testing.T,
			numbers chan int,
		) []tcpServer.RequestMiddleware[request, response]
	}
	type handlerArgs struct {
		ctx     context.Context
		request request
	}

	for _, data := range []struct {
		name                string
		args                args
		handlerArgs         handlerArgs
		wantNumberSlice     []int
		wantHandlerResponse response
		wantHandlerErr      assert.ErrorAssertionFunc
	}{
		{
			name: "success/without middlewares",
			args: args{
				handler: func(test *testing.T) tcpServer.RequestHandler[request, response] {
					requestHandlerMock :=
						tcpServerExternalMocks.NewMockRequestHandler[request, response](test)
					requestHandlerMock.EXPECT().
						HandleRequest(context.Background(), request("request")).
						Return("response", nil)

					return requestHandlerMock
				},
				middlewares: func(
					test *testing.T,
					numbers chan int,
				) []tcpServer.RequestMiddleware[request, response] {
					return nil
				},
			},
			handlerArgs: handlerArgs{
				ctx:     context.Background(),
				request: "request",
			},
			wantNumberSlice:     []int{},
			wantHandlerResponse: "response",
			wantHandlerErr:      assert.NoError,
		},
		{
			name: "success/with middlewares",
			args: args{
				handler: func(test *testing.T) tcpServer.RequestHandler[request, response] {
					requestHandlerMock :=
						tcpServerExternalMocks.NewMockRequestHandler[request, response](test)
					requestHandlerMock.EXPECT().
						HandleRequest(context.Background(), request("request")).
						Return("response", nil)

					return requestHandlerMock
				},
				middlewares: func(
					test *testing.T,
					numbers chan int,
				) []tcpServer.RequestMiddleware[request, response] {
					return []tcpServer.RequestMiddleware[request, response]{
						func(
							handler tcpServer.RequestHandler[request, response],
						) tcpServer.RequestHandler[request, response] {
							return tcpServer.RequestHandlerFunc[request, response](func(
								ctx context.Context,
								request request,
							) (response, error) {
								numbers <- 23
								return handler.HandleRequest(ctx, request)
							})
						},
						func(
							handler tcpServer.RequestHandler[request, response],
						) tcpServer.RequestHandler[request, response] {
							return tcpServer.RequestHandlerFunc[request, response](func(
								ctx context.Context,
								request request,
							) (response, error) {
								numbers <- 42
								return handler.HandleRequest(ctx, request)
							})
						},
					}
				},
			},
			handlerArgs: handlerArgs{
				ctx:     context.Background(),
				request: "request",
			},
			wantNumberSlice:     []int{42, 23},
			wantHandlerResponse: "response",
			wantHandlerErr:      assert.NoError,
		},
		{
			name: "error/with middlewares",
			args: args{
				handler: func(test *testing.T) tcpServer.RequestHandler[request, response] {
					requestHandlerMock :=
						tcpServerExternalMocks.NewMockRequestHandler[request, response](test)
					requestHandlerMock.EXPECT().
						HandleRequest(context.Background(), request("request")).
						Return("", iotest.ErrTimeout)

					return requestHandlerMock
				},
				middlewares: func(
					test *testing.T,
					numbers chan int,
				) []tcpServer.RequestMiddleware[request, response] {
					return []tcpServer.RequestMiddleware[request, response]{
						func(
							handler tcpServer.RequestHandler[request, response],
						) tcpServer.RequestHandler[request, response] {
							return tcpServer.RequestHandlerFunc[request, response](func(
								ctx context.Context,
								request request,
							) (response, error) {
								numbers <- 23
								return handler.HandleRequest(ctx, request)
							})
						},
						func(
							handler tcpServer.RequestHandler[request, response],
						) tcpServer.RequestHandler[request, response] {
							return tcpServer.RequestHandlerFunc[request, response](func(
								ctx context.Context,
								request request,
							) (response, error) {
								numbers <- 42
								return handler.HandleRequest(ctx, request)
							})
						},
					}
				},
			},
			handlerArgs: handlerArgs{
				ctx:     context.Background(),
				request: "request",
			},
			wantNumberSlice:     []int{42, 23},
			wantHandlerResponse: "",
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
			gotHandler := tcpServer.ApplyRequestMiddlewares(
				data.args.handler(test),
				data.args.middlewares(test, numbers),
			)

			gotHandlerResponse, handlerErr := gotHandler.HandleRequest(
				data.handlerArgs.ctx,
				data.handlerArgs.request,
			)

			close(numbers)

			gotNumberSlice := make([]int, 0, maxNumberCount)
			for number := range numbers {
				gotNumberSlice = append(gotNumberSlice, number)
			}

			assert.Equal(test, data.wantNumberSlice, gotNumberSlice)
			assert.Equal(test, data.wantHandlerResponse, gotHandlerResponse)
			data.wantHandlerErr(test, handlerErr)
		})
	}
}
