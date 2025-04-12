package tcpServer

type RequestMiddleware[Req Request, Resp Response] func(
	handler RequestHandler[Req, Resp],
) RequestHandler[Req, Resp]

func ApplyRequestMiddlewares[Req Request, Resp Response](
	handler RequestHandler[Req, Resp],
	middlewares []RequestMiddleware[Req, Resp],
) RequestHandler[Req, Resp] {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}

	return handler
}
