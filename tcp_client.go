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

type ClientProtocol[Req Request, Resp Response] interface {
	InitialScannerBufferSize() int
	MaxTokenSize() int
	ExtractToken(
		data []byte,
		isLatestData bool,
	) (offsetToNextToken int, token []byte, err error)
	MarshalRequest(request Req) ([]byte, error)
	ParseResponse(data []byte) (Resp, error)
}

type TCPClientOptions[Req Request, Resp Response] struct {
	ReadTimeout    mo.Option[time.Duration]
	WriteTimeout   mo.Option[time.Duration]
	ClientProtocol ClientProtocol[Req, Resp]
}

type TCPClient[Req Request, Resp Response] struct {
	options    TCPClientOptions[Req, Resp]
	connection net.Conn
	scanner    *bufio.Scanner
}

func NewTCPClient[Req Request, Resp Response](
	ctx context.Context,
	address string,
	options TCPClientOptions[Req, Resp],
) (TCPClient[Req, Resp], error) {
	var dialer net.Dialer
	connection, err := dialer.DialContext(ctx, TCPServerNetwork, address)
	if err != nil {
		return TCPClient[Req, Resp]{}, fmt.Errorf(
			"unable to connect to address %q: %w",
			address,
			err,
		)
	}

	return NewTCPClientFromConnection(connection, options), nil
}

func NewTCPClientFromConnection[Req Request, Resp Response](
	connection net.Conn,
	options TCPClientOptions[Req, Resp],
) TCPClient[Req, Resp] {
	scanner := bufio.NewScanner(connection)
	scanner.Buffer(
		make([]byte, options.ClientProtocol.InitialScannerBufferSize()),
		options.ClientProtocol.MaxTokenSize(),
	)
	scanner.Split(options.ClientProtocol.ExtractToken)

	return TCPClient[Req, Resp]{
		options:    options,
		connection: connection,
		scanner:    scanner,
	}
}

func (client TCPClient[Req, Resp]) SendRequest(request Req) (Resp, error) {
	var zeroResponse Resp

	marshalledRequest, err := client.options.ClientProtocol.MarshalRequest(request)
	if err != nil {
		return zeroResponse, fmt.Errorf("unable to marshal the request: %w", err)
	}

	if writeTimeout, isPresent := client.options.WriteTimeout.Get(); isPresent {
		writeDeadline := time.Now().Add(writeTimeout)
		if err := client.connection.SetWriteDeadline(writeDeadline); err != nil {
			return zeroResponse, fmt.Errorf("unable to set the write deadline: %w", err)
		}
	}

	if _, err := client.connection.Write(marshalledRequest); err != nil {
		return zeroResponse, fmt.Errorf("unable to write the request: %w", err)
	}

	if readTimeout, isPresent := client.options.ReadTimeout.Get(); isPresent {
		readDeadline := time.Now().Add(readTimeout)
		if err := client.connection.SetReadDeadline(readDeadline); err != nil {
			return zeroResponse, fmt.Errorf("unable to set the read deadline: %w", err)
		}
	}

	if isPossibleToContinue := client.scanner.Scan(); !isPossibleToContinue {
		if err := client.scanner.Err(); err != nil {
			return zeroResponse, fmt.Errorf("unable to read the response: %w", err)
		}

		return zeroResponse, errors.Join(
			errors.New("scanner has no more tokens"),
			ErrHandlingStopIsRequired,
		)
	}

	response, err := client.options.ClientProtocol.ParseResponse(
		slices.Clone(client.scanner.Bytes()),
	)
	if err != nil {
		return zeroResponse, fmt.Errorf("unable to parse the response: %w", err)
	}

	return response, nil
}

func (client TCPClient[Req, Resp]) Close() error {
	if err := client.connection.Close(); err != nil {
		return fmt.Errorf("unable to close the connection: %w", err)
	}

	return nil
}
