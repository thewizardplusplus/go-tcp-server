package tcpServer

import (
	"context"
	"net"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	tcpServerMocks "github.com/thewizardplusplus/go-tcp-server/mocks/github.com/thewizardplusplus/go-tcp-server"
)

func TestConnectionHandlerFunc_interface(test *testing.T) {
	assert.Implements(test, (*ConnectionHandler)(nil), ConnectionHandlerFunc(nil))
}

func TestConnectionHandlerFunc_HandleConnection(test *testing.T) {
	type args struct {
		ctx        context.Context
		connection func(test *testing.T) net.Conn
	}

	for _, data := range []struct {
		name    string
		f       func(test *testing.T) ConnectionHandlerFunc
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			f: func(test *testing.T) ConnectionHandlerFunc {
				connectionHandlerMock := tcpServerMocks.NewMockConnectionHandler(test)
				connectionHandlerMock.EXPECT().
					HandleConnection(
						context.Background(),
						tcpServerMocks.NewMocknetConn(test),
					).
					Return(nil)

				return connectionHandlerMock.HandleConnection
			},
			args: args{
				ctx: context.Background(),
				connection: func(test *testing.T) net.Conn {
					return tcpServerMocks.NewMocknetConn(test)
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "error",
			f: func(test *testing.T) ConnectionHandlerFunc {
				connectionHandlerMock := tcpServerMocks.NewMockConnectionHandler(test)
				connectionHandlerMock.EXPECT().
					HandleConnection(
						context.Background(),
						tcpServerMocks.NewMocknetConn(test),
					).
					Return(iotest.ErrTimeout)

				return connectionHandlerMock.HandleConnection
			},
			args: args{
				ctx: context.Background(),
				connection: func(test *testing.T) net.Conn {
					return tcpServerMocks.NewMocknetConn(test)
				},
			},
			wantErr: func(test assert.TestingT, err error, msgAndArgs ...any) bool {
				return assert.ErrorIs(test, err, iotest.ErrTimeout)
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			err := data.f(test).HandleConnection(
				data.args.ctx,
				data.args.connection(test),
			)

			data.wantErr(test, err)
		})
	}
}
