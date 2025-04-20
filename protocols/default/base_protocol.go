package defaultProtocol

import (
	"fmt"

	defaultProtocolModels "github.com/thewizardplusplus/go-tcp-server/protocols/default/models"
)

type MessageFormat interface {
	ParseMessage(data []byte) (defaultProtocolModels.Message, error)
	MarshalMessage(message defaultProtocolModels.Message) ([]byte, error)
}

type BaseProtocolOptions struct {
	MessageFormat MessageFormat
}

type BaseProtocol struct {
	options BaseProtocolOptions
}

func NewBaseProtocol(options BaseProtocolOptions) BaseProtocol {
	return BaseProtocol{
		options: options,
	}
}

func (protocol BaseProtocol) InitialScannerBufferSize() int {
	return 4 * 1024 // based on the default values in package `bufio`
}

func (protocol BaseProtocol) MaxTokenSize() int {
	return 64 * 1024 // based on the default values in package `bufio`
}

func (protocol BaseProtocol) ParseRequest(
	data []byte,
) (defaultProtocolModels.Request, error) {
	message, err := protocol.options.MessageFormat.ParseMessage(data)
	if err != nil {
		return defaultProtocolModels.Request{}, fmt.Errorf(
			"unable to parse the message: %w",
			err,
		)
	}

	request, err := defaultProtocolModels.NewRequestFromMessage(message)
	if err != nil {
		return defaultProtocolModels.Request{}, fmt.Errorf(
			"unable to construct the request: %w",
			err,
		)
	}

	return request, nil
}

func (protocol BaseProtocol) ParseResponse(
	data []byte,
) (defaultProtocolModels.Response, error) {
	message, err := protocol.options.MessageFormat.ParseMessage(data)
	if err != nil {
		return defaultProtocolModels.Response{}, fmt.Errorf(
			"unable to parse the message: %w",
			err,
		)
	}

	response, err := defaultProtocolModels.NewResponseFromMessage(message)
	if err != nil {
		return defaultProtocolModels.Response{}, fmt.Errorf(
			"unable to construct the response: %w",
			err,
		)
	}

	return response, nil
}

func (protocol BaseProtocol) MarshalRequest(
	request defaultProtocolModels.Request,
) ([]byte, error) {
	message, err := request.ToMessage()
	if err != nil {
		return nil, fmt.Errorf(
			"unable to convert the request to the message: %w",
			err,
		)
	}

	marshalledMessage, err := protocol.options.MessageFormat.MarshalMessage(
		message,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal the message: %w", err)
	}

	return marshalledMessage, nil
}

func (protocol BaseProtocol) MarshalResponse(
	response defaultProtocolModels.Response,
) ([]byte, error) {
	message, err := response.ToMessage()
	if err != nil {
		return nil, fmt.Errorf(
			"unable to convert the response to the message: %w",
			err,
		)
	}

	marshalledMessage, err := protocol.options.MessageFormat.MarshalMessage(
		message,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal the message: %w", err)
	}

	return marshalledMessage, nil
}
