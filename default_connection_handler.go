package tcpServer

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net"
	"slices"
	"time"

	"github.com/samber/mo"
)

var (
	ErrHandlingStopIsRequired = errors.New("handling stop is required")
)

type Request any

type Response any

type ServerProtocol[Req Request, Resp Response] interface {
	InitialScannerBufferSize() int
	MaxTokenSize() int
	ExtractToken(
		data []byte,
		isLatestData bool,
	) (offsetToNextToken int, token []byte, err error)
	ParseRequest(token []byte) (Req, error)
	MarshalResponse(response Resp) ([]byte, error)
}

type DefaultConnectionHandlerOptions[Req Request, Resp Response] struct {
	ReadTimeout     mo.Option[time.Duration]
	WriteTimeout    mo.Option[time.Duration]
	HandlingTimeout mo.Option[time.Duration]
	ServerProtocol  ServerProtocol[Req, Resp]
	RequestHandler  RequestHandler[Req, Resp]
}

type DefaultConnectionHandler[Req Request, Resp Response] struct {
	options DefaultConnectionHandlerOptions[Req, Resp]
}

func NewDefaultConnectionHandler[Req Request, Resp Response](
	options DefaultConnectionHandlerOptions[Req, Resp],
) DefaultConnectionHandler[Req, Resp] {
	return DefaultConnectionHandler[Req, Resp]{
		options: options,
	}
}

func (handler DefaultConnectionHandler[Req, Resp]) HandleRequest(
	ctx context.Context,
	connection net.Conn,
	scanner *bufio.Scanner,
) error {
	if readTimeout, isPresent := handler.options.ReadTimeout.Get(); isPresent {
		readDeadline := time.Now().Add(readTimeout)
		if err := connection.SetReadDeadline(readDeadline); err != nil {
			return fmt.Errorf("unable to set the read deadline: %w", err)
		}
	}

	if isPossibleToContinue := scanner.Scan(); !isPossibleToContinue {
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("unable to read the request: %w", err)
		}

		return errors.Join(
			errors.New("scanner has no more tokens"),
			ErrHandlingStopIsRequired,
		)
	}

	request, err := handler.options.ServerProtocol.ParseRequest(
		slices.Clone(scanner.Bytes()),
	)
	if err != nil {
		return fmt.Errorf("unable to parse the request: %w", err)
	}

	// ignore the parent context cancellation, because even in this case
	// we need to finish handling the request
	handlingCtx := context.WithoutCancel(ctx)

	if handlingTimeout, isPresent :=
		handler.options.HandlingTimeout.Get(); isPresent {
		var handlingCtxCancel func()
		handlingCtx, handlingCtxCancel = context.WithTimeout(
			handlingCtx,
			handlingTimeout,
		)
		defer handlingCtxCancel()
	}

	response, handlingErr := handler.options.RequestHandler.HandleRequest(
		handlingCtx,
		request,
	)
	if handlingErr != nil && !errors.Is(handlingErr, ErrHandlingStopIsRequired) {
		return fmt.Errorf("unable to handle the request: %w", err)
	}

	marshalledResponse, err := handler.options.ServerProtocol.MarshalResponse(
		response,
	)
	if err != nil {
		return fmt.Errorf("unable to marshal the response: %w", err)
	}

	if writeTimeout, isPresent := handler.options.WriteTimeout.Get(); isPresent {
		writeDeadline := time.Now().Add(writeTimeout)
		if err := connection.SetWriteDeadline(writeDeadline); err != nil {
			return fmt.Errorf("unable to set the write deadline: %w", err)
		}
	}

	if _, err := connection.Write(marshalledResponse); err != nil {
		return fmt.Errorf("unable to write the response: %w", err)
	}

	if errors.Is(handlingErr, ErrHandlingStopIsRequired) {
		return fmt.Errorf(
			"request handler requested to stop handling: %w",
			handlingErr,
		)
	}

	return nil
}

func (handler DefaultConnectionHandler[Req, Resp]) HandleConnection(
	ctx context.Context,
	connection net.Conn,
) error {
	scanner := bufio.NewScanner(connection)
	scanner.Buffer(
		make([]byte, handler.options.ServerProtocol.InitialScannerBufferSize()),
		handler.options.ServerProtocol.MaxTokenSize(),
	)
	scanner.Split(handler.options.ServerProtocol.ExtractToken)

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context is done: %w", ctx.Err())

		default:
		}

		if err := handler.HandleRequest(ctx, connection, scanner); err != nil {
			if errors.Is(err, ErrHandlingStopIsRequired) {
				break
			}

			return fmt.Errorf("unable to handle the request: %w", err)
		}
	}

	return nil
}
