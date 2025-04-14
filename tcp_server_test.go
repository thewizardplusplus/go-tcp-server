package tcpServer

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"testing"
	"testing/iotest"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	tcpServerMocks "github.com/thewizardplusplus/go-tcp-server/mocks/github.com/thewizardplusplus/go-tcp-server"
)

func TestNewTCPServer(test *testing.T) {
	type args struct {
		ctx     context.Context
		options func(test *testing.T) TCPServerOptions
	}

	for _, data := range []struct {
		name    string
		args    args
		want    func(test *testing.T, got *TCPServer)
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				options: func(test *testing.T) TCPServerOptions {
					connectionHandlerMock := tcpServerMocks.NewMockConnectionHandler(test)
					errorHandlerMock := tcpServerMocks.NewMockerrorHandler(test)
					return TCPServerOptions{
						Address:           "127.0.0.1:",
						ConnectionHandler: connectionHandlerMock,
						ErrorHandler:      errorHandlerMock.HandleError,
					}
				},
			},
			want: func(test *testing.T, got *TCPServer) {
				if !assert.NotNil(test, got) {
					return
				}

				assert.Equal(test, "127.0.0.1:", got.options.Address)
				assert.Equal(
					test,
					tcpServerMocks.NewMockConnectionHandler(test),
					got.options.ConnectionHandler,
				)
				assert.NotEmpty(test, got.options.ErrorHandler)
				assert.NotEmpty(test, got.listener)
				assert.Equal(test, false, got.isStopped.Load())
			},
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
				options: func(test *testing.T) TCPServerOptions {
					connectionHandlerMock := tcpServerMocks.NewMockConnectionHandler(test)
					errorHandlerMock := tcpServerMocks.NewMockerrorHandler(test)
					return TCPServerOptions{
						Address:           "https://example.com/",
						ConnectionHandler: connectionHandlerMock,
						ErrorHandler:      errorHandlerMock.HandleError,
					}
				},
			},
			want: func(test *testing.T, got *TCPServer) {
				assert.Nil(test, got)
			},
			wantErr: assert.Error,
		},
		{
			name: "error/invalid address",
			args: args{
				ctx: context.Background(),
				options: func(test *testing.T) TCPServerOptions {
					connectionHandlerMock := tcpServerMocks.NewMockConnectionHandler(test)
					errorHandlerMock := tcpServerMocks.NewMockerrorHandler(test)
					return TCPServerOptions{
						Address:           "127.0.0.1",
						ConnectionHandler: connectionHandlerMock,
						ErrorHandler:      errorHandlerMock.HandleError,
					}
				},
			},
			want: func(test *testing.T, got *TCPServer) {
				assert.Nil(test, got)
			},
			wantErr: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got, err := NewTCPServer(data.args.ctx, data.args.options(test))

			data.want(test, got)
			data.wantErr(test, err)
		})
	}
}

func TestTCPServer_Run(test *testing.T) {
	type constructorArgs struct {
		ctx     context.Context
		options func(test *testing.T, serverDone chan struct{}) TCPServerOptions
	}
	type args struct {
		ctx context.Context
	}
	type runTestClientsParams struct {
		address string
		timeout time.Duration
	}

	for _, data := range []struct {
		name            string
		constructorArgs constructorArgs
		args            args
		runTestClients  func(test *testing.T, params runTestClientsParams) int
		timeout         time.Duration
	}{
		{
			name: "success/single client",
			constructorArgs: constructorArgs{
				ctx: context.Background(),
				options: func(test *testing.T, serverDone chan struct{}) TCPServerOptions {
					connectionHandlerMock := tcpServerMocks.NewMockConnectionHandler(test)
					connectionHandlerMock.EXPECT().
						HandleConnection(
							mock.MatchedBy(func(context.Context) bool { return true }),
							mock.MatchedBy(func(net.Conn) bool { return true }),
						).
						RunAndReturn(func(ctx context.Context, connection net.Conn) error {
							defer func() { serverDone <- struct{}{} }()

							if err := runTestHandler(connection); err != nil {
								return fmt.Errorf("unable to run the test handler: %w", err)
							}

							return nil
						})

					errorHandlerMock := tcpServerMocks.NewMockerrorHandler(test)
					errorHandlerMock.EXPECT().
						HandleError(mock.MatchedBy(func(err error) bool {
							if err == nil {
								return false
							}

							_, isAsserted := errors.Unwrap(err).(*net.OpError)
							return isAsserted
						})).
						Return()

					return TCPServerOptions{
						Address:           "127.0.0.1:",
						ConnectionHandler: connectionHandlerMock,
						ErrorHandler:      errorHandlerMock.HandleError,
					}
				},
			},
			args: args{
				ctx: context.Background(),
			},
			runTestClients: func(test *testing.T, params runTestClientsParams) int {
				runTestClient(test, runTestClientParams{
					address: params.address,
					timeout: params.timeout,
					index:   0,
				})

				return 1
			},
			timeout: 5 * time.Second,
		},
		{
			name: "success/several clients",
			constructorArgs: constructorArgs{
				ctx: context.Background(),
				options: func(test *testing.T, serverDone chan struct{}) TCPServerOptions {
					connectionHandlerMock := tcpServerMocks.NewMockConnectionHandler(test)
					connectionHandlerMock.EXPECT().
						HandleConnection(
							mock.MatchedBy(func(context.Context) bool { return true }),
							mock.MatchedBy(func(net.Conn) bool { return true }),
						).
						RunAndReturn(func(ctx context.Context, connection net.Conn) error {
							defer func() { serverDone <- struct{}{} }()

							if err := runTestHandler(connection); err != nil {
								return fmt.Errorf("unable to run the test handler: %w", err)
							}

							return nil
						})

					errorHandlerMock := tcpServerMocks.NewMockerrorHandler(test)
					errorHandlerMock.EXPECT().
						HandleError(mock.MatchedBy(func(err error) bool {
							if err == nil {
								return false
							}

							_, isAsserted := errors.Unwrap(err).(*net.OpError)
							return isAsserted
						})).
						Return()

					return TCPServerOptions{
						Address:           "127.0.0.1:",
						ConnectionHandler: connectionHandlerMock,
						ErrorHandler:      errorHandlerMock.HandleError,
					}
				},
			},
			args: args{
				ctx: context.Background(),
			},
			runTestClients: func(test *testing.T, params runTestClientsParams) int {
				const clientCount = 5

				var waitGroup sync.WaitGroup
				waitGroup.Add(clientCount)

				for clientIndex := range clientCount {
					go func(clientIndex int) {
						defer waitGroup.Done()

						runTestClient(test, runTestClientParams{
							address: params.address,
							timeout: params.timeout,
							index:   clientIndex,
						})
					}(clientIndex)
				}

				waitGroup.Wait()

				return clientCount
			},
			timeout: 5 * time.Second,
		},
		{
			name: "error/unable to handle the connection",
			constructorArgs: constructorArgs{
				ctx: context.Background(),
				options: func(test *testing.T, serverDone chan struct{}) TCPServerOptions {
					connectionHandlerMock := tcpServerMocks.NewMockConnectionHandler(test)
					connectionHandlerMock.EXPECT().
						HandleConnection(
							mock.MatchedBy(func(context.Context) bool { return true }),
							mock.MatchedBy(func(net.Conn) bool { return true }),
						).
						RunAndReturn(func(ctx context.Context, connection net.Conn) error {
							defer func() { serverDone <- struct{}{} }()

							if err := runTestHandler(connection); err != nil {
								return fmt.Errorf("unable to run the test handler: %w", err)
							}

							return iotest.ErrTimeout
						})

					errorHandlerMock := tcpServerMocks.NewMockerrorHandler(test)
					errorHandlerMock.EXPECT().
						HandleError(mock.MatchedBy(func(err error) bool {
							if err == nil {
								return false
							}

							_, isAsserted := errors.Unwrap(err).(*net.OpError)
							return isAsserted || errors.Is(err, iotest.ErrTimeout)
						})).
						Return()

					return TCPServerOptions{
						Address:           "127.0.0.1:",
						ConnectionHandler: connectionHandlerMock,
						ErrorHandler:      errorHandlerMock.HandleError,
					}
				},
			},
			args: args{
				ctx: context.Background(),
			},
			runTestClients: func(test *testing.T, params runTestClientsParams) int {
				runTestClient(test, runTestClientParams{
					address: params.address,
					timeout: params.timeout,
					index:   0,
				})

				return 1
			},
			timeout: 5 * time.Second,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			serverRunningDone := make(chan struct{})
			server, err := NewTCPServer(
				data.constructorArgs.ctx,
				data.constructorArgs.options(test, serverRunningDone),
			)
			require.NoError(test, err)

			serverStoppingDone := make(chan struct{})
			go func() {
				defer func() { close(serverStoppingDone) }()

				server.Run(data.args.ctx)
			}()

			clientCount := data.runTestClients(test, runTestClientsParams{
				address: server.Address(),
				timeout: data.timeout,
			})
			for range clientCount {
				select {
				case <-serverRunningDone:
				case <-time.After(data.timeout):
				}
			}

			server.Stop()

			select {
			case <-serverStoppingDone:
			case <-time.After(data.timeout):
			}
		})
	}
}

type runTestClientParams struct {
	address string
	timeout time.Duration
	index   int
}

func runTestClient(test *testing.T, params runTestClientParams) {
	connection, err := net.Dial(TCPServerNetwork, params.address)
	require.NoError(test, err)
	defer connection.Close()

	err = connection.SetDeadline(time.Now().Add(params.timeout))
	require.NoError(test, err)

	sentContent := fmt.Sprintf("dummy #%d\n", params.index)
	_, err = connection.Write([]byte(sentContent))
	require.NoError(test, err)

	gotContent, err := bufio.NewReader(connection).ReadString('\n')
	require.NoError(test, err)

	assert.Equal(test, sentContent, gotContent)
}

func runTestHandler(connection net.Conn) error {
	content, err := bufio.NewReader(connection).ReadBytes('\n')
	if err != nil {
		return fmt.Errorf("unable to read the request: %w", err)
	}

	if _, err = connection.Write(content); err != nil {
		return fmt.Errorf("unable to write the response: %w", err)
	}

	return nil
}
