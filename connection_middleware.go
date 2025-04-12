package tcpServer

type ConnectionMiddleware func(handler ConnectionHandler) ConnectionHandler

func ApplyConnectionMiddlewares(
	handler ConnectionHandler,
	middlewares []ConnectionMiddleware,
) ConnectionHandler {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}

	return handler
}
