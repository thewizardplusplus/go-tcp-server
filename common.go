package tcpServer

import (
	"bufio"
	"errors"
	"io"
)

var (
	ErrHandlingStopIsRequired = errors.New("handling stop is required")
)

type Request any

type Response any

type BaseProtocol[Req Request, Resp Response] interface {
	InitialScannerBufferSize() int
	MaxTokenSize() int
	ExtractToken(
		data []byte,
		isLatestData bool,
	) (offsetToNextToken int, token []byte, err error)
}

type InitializeScannerParams[Req Request, Resp Response] struct {
	Reader       io.Reader
	BaseProtocol BaseProtocol[Req, Resp]
}

func InitializeScanner[Req Request, Resp Response](
	params InitializeScannerParams[Req, Resp],
) *bufio.Scanner {
	scanner := bufio.NewScanner(params.Reader)
	scanner.Buffer(
		make([]byte, params.BaseProtocol.InitialScannerBufferSize()),
		params.BaseProtocol.MaxTokenSize(),
	)
	scanner.Split(params.BaseProtocol.ExtractToken)

	return scanner
}
