package separatorBasedTCPServerProtocol

import (
	"bufio"
	"bytes"

	defaultProtocol "github.com/thewizardplusplus/go-tcp-server/protocols/default"
)

type Protocol struct {
	defaultProtocol.BaseProtocol

	options SeparationParams
}

func NewProtocol(options SeparationParams) Protocol {
	return Protocol{
		BaseProtocol: defaultProtocol.NewBaseProtocol(
			defaultProtocol.BaseProtocolOptions{
				MessageFormat: NewMessageFormat(options),
			},
		),

		options: options,
	}
}

func (protocol Protocol) ExtractToken(
	data []byte,
	isLatestData bool,
) (offsetToNextToken int, token []byte, err error) {
	separator := protocol.options.MessageSeparator

	separatorIndex := bytes.Index(data, separator)
	if separatorIndex == -1 {
		if !isLatestData {
			return 0, nil, nil // request more data
		}

		if len(data) == 0 {
			data = nil // don't generate a new token if the data is empty
		}
		return 0, data, bufio.ErrFinalToken
	}

	return separatorIndex + len(separator), data[:separatorIndex], nil
}
