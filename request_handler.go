package tcpServer

import (
	"context"
)

type RequestHandler[Req Request, Resp Response] interface {
	HandleRequest(ctx context.Context, request Req) (Resp, error)
}

type RequestHandlerFunc[Req Request, Resp Response] func(
	ctx context.Context,
	request Req,
) (Resp, error)

func (f RequestHandlerFunc[Req, Resp]) HandleRequest(
	ctx context.Context,
	request Req,
) (Resp, error) {
	return f(ctx, request)
}
