package tcpServer_test

import (
	"context"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	tcpServer "github.com/thewizardplusplus/go-tcp-server"
	tcpServerExternalMocks "github.com/thewizardplusplus/go-tcp-server/mocks/external/github.com/thewizardplusplus/go-tcp-server"
)

func TestRequestHandlerFunc_interface(test *testing.T) {
	type request string
	type response string

	assert.Implements(
		test,
		(*tcpServer.RequestHandler[request, response])(nil),
		tcpServer.RequestHandlerFunc[request, response](nil),
	)
}

func TestRequestHandlerFunc_HandleRequest(test *testing.T) {
	type request string
	type response string
	type args struct {
		ctx     context.Context
		request request
	}

	for _, data := range []struct {
		name    string
		f       func(test *testing.T) tcpServer.RequestHandlerFunc[request, response]
		args    args
		want    response
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			f: func(test *testing.T) tcpServer.RequestHandlerFunc[request, response] {
				requestHandlerMock :=
					tcpServerExternalMocks.NewMockRequestHandler[request, response](test)
				requestHandlerMock.EXPECT().
					HandleRequest(context.Background(), request("request")).
					Return("response", nil)

				return requestHandlerMock.HandleRequest
			},
			args: args{
				ctx:     context.Background(),
				request: "request",
			},
			want:    "response",
			wantErr: assert.NoError,
		},
		{
			name: "error",
			f: func(test *testing.T) tcpServer.RequestHandlerFunc[request, response] {
				requestHandlerMock :=
					tcpServerExternalMocks.NewMockRequestHandler[request, response](test)
				requestHandlerMock.EXPECT().
					HandleRequest(context.Background(), request("request")).
					Return("", iotest.ErrTimeout)

				return requestHandlerMock.HandleRequest
			},
			args: args{
				ctx:     context.Background(),
				request: "request",
			},
			want: "",
			wantErr: func(test assert.TestingT, err error, msgAndArgs ...any) bool {
				return assert.ErrorIs(test, err, iotest.ErrTimeout)
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got, err := data.f(test).HandleRequest(
				data.args.ctx,
				data.args.request,
			)

			assert.Equal(test, data.want, got)
			data.wantErr(test, err)
		})
	}
}
