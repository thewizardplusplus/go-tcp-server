package tcpServer

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/samber/mo"
)

type RouteExtractor[Req Request] func(
	ctx context.Context,
	request Req,
) (route []byte, err error)

type RequestRouterOptions[Req Request] struct {
	RouteExtractor RouteExtractor[Req]
}

type RequestRouter[Req Request, Resp Response] struct {
	options          RequestRouterOptions[Req]
	handlersByRoutes map[string]RequestHandler[Req, Resp]
	notFoundHandler  mo.Option[RequestHandler[Req, Resp]]
}

func NewRequestRouter[Req Request, Resp Response](
	options RequestRouterOptions[Req],
) *RequestRouter[Req, Resp] {
	return &RequestRouter[Req, Resp]{
		options:          options,
		handlersByRoutes: make(map[string]RequestHandler[Req, Resp]),
	}
}

func (router RequestRouter[Req, Resp]) SetRouteHandler(
	route []byte,
	handler RequestHandler[Req, Resp],
) {
	router.handlersByRoutes[hex.EncodeToString(route)] = handler
}

func (router *RequestRouter[Req, Resp]) SetNotFoundHandler(
	handler RequestHandler[Req, Resp],
) {
	router.notFoundHandler = mo.Some(handler)
}

func (router RequestRouter[Req, Resp]) HandleRequest(
	ctx context.Context,
	request Req,
) (Resp, error) {
	var zeroResponse Resp

	route, err := router.options.RouteExtractor(ctx, request)
	if err != nil {
		return zeroResponse, fmt.Errorf(
			"unable to extract the request route: %w",
			err,
		)
	}

	handler, isFound := router.handlersByRoutes[hex.EncodeToString(route)]
	if !isFound {
		notFoundHandler, isPresent := router.notFoundHandler.Get()
		if !isPresent {
			return zeroResponse, errors.New("`Not Found` handler isn't set")
		}

		handler = notFoundHandler
	}

	return handler.HandleRequest(ctx, request)
}
