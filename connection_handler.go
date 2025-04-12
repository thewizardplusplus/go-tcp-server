package tcpServer

import (
	"context"
	"net"
)

type ConnectionHandler interface {
	HandleConnection(ctx context.Context, connection net.Conn) error
}

type ConnectionHandlerFunc func(ctx context.Context, connection net.Conn) error

func (f ConnectionHandlerFunc) HandleConnection(
	ctx context.Context,
	connection net.Conn,
) error {
	return f(ctx, connection)
}
