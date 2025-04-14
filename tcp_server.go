package tcpServer

import (
	"context"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
)

const (
	TCPServerNetwork = "tcp"
)

type ErrorHandler func(err error)

// is only used to generate the mock
type errorHandler interface { //nolint:unused
	HandleError(err error)
}

type TCPServerOptions struct {
	Address           string
	ConnectionHandler ConnectionHandler
	ErrorHandler      ErrorHandler
}

type TCPServer struct {
	options   TCPServerOptions
	listener  net.Listener
	isStopped atomic.Bool
}

func NewTCPServer(
	ctx context.Context,
	options TCPServerOptions,
) (*TCPServer, error) {
	var listenConfig net.ListenConfig
	listener, err := listenConfig.Listen(ctx, TCPServerNetwork, options.Address)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to listen on address %q: %w",
			options.Address,
			err,
		)
	}

	server := &TCPServer{
		options:  options,
		listener: listener,
	}
	return server, nil
}

func (server *TCPServer) Address() string {
	return server.listener.Addr().String()
}

func (server *TCPServer) Run(ctx context.Context) {
	var waitGroup sync.WaitGroup
	ctx, ctxCancel := context.WithCancel(ctx)
	for !server.isStopped.Load() {
		connection, err := server.listener.Accept()
		if err != nil {
			server.options.ErrorHandler(fmt.Errorf(
				"unable to accept the connection: %w",
				err,
			))
			continue
		}

		waitGroup.Add(1)

		go func() {
			defer waitGroup.Done()
			defer connection.Close()

			if err := server.options.ConnectionHandler.HandleConnection(
				ctx,
				connection,
			); err != nil {
				server.options.ErrorHandler(fmt.Errorf(
					"unable to handle the connection: %w",
					err,
				))
			}
		}()
	}

	ctxCancel()
	waitGroup.Wait()
}

func (server *TCPServer) Stop() {
	server.isStopped.Store(true)
	server.listener.Close()
}
