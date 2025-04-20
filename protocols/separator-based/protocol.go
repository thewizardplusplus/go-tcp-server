//go:build ignore

package separatorBasedTCPServerProtocol

import (
	"bufio"
	"bytes"
	"fmt"
)

type Request struct {
	Action  []byte
	Headers map[string][]byte
	Body    []byte
}

type Response struct {
	Status  []byte
	Headers map[string][]byte
	Body    []byte
}

type ProtocolOptions struct {
	SeparationParams SeparationParams
}

type Protocol struct {
	options ProtocolOptions
}

func NewProtocol(options ProtocolOptions) Protocol {
	return Protocol{
		options: options,
	}
}

func (protocol Protocol) InitialScannerBufferSize() int {
	return 4 * 1024 // based on the default values in package `bufio`
}

func (protocol Protocol) MaxTokenSize() int {
	return 64 * 1024 // based on the default values in package `bufio`
}

func (protocol Protocol) ExtractToken(
	data []byte,
	isLatestData bool,
) (offsetToNextToken int, token []byte, err error) {
	separator := protocol.options.SeparationParams.MessageSeparator

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

func (protocol Protocol) ParseRequest(data []byte) (Request, error) {
	message, err := ParseMessage(data, protocol.options.SeparationParams)
	if err != nil {
		return Request{}, fmt.Errorf("unable to parse the message: %w", err)
	}

	request := Request{
		Action:  message.Introduction,
		Headers: message.Headers,
		Body:    message.Body,
	}
	return request, nil
}

func (protocol Protocol) ParseResponse(data []byte) (Response, error) {
	message, err := ParseMessage(data, protocol.options.SeparationParams)
	if err != nil {
		return Response{}, fmt.Errorf("unable to parse the message: %w", err)
	}

	response := Response{
		Status:  message.Introduction,
		Headers: message.Headers,
		Body:    message.Body,
	}
	return response, nil
}

func (protocol Protocol) MarshalRequest(request Request) ([]byte, error) {
	marshalledMessage, err := MarshalMessage(
		Message{
			Introduction: request.Action,
			Headers:      request.Headers,
			Body:         request.Body,
		},
		protocol.options.SeparationParams,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal the message: %w", err)
	}

	return marshalledMessage, nil
}

func (protocol Protocol) MarshalResponse(response Response) ([]byte, error) {
	marshalledMessage, err := MarshalMessage(
		Message{
			Introduction: response.Status,
			Headers:      response.Headers,
			Body:         response.Body,
		},
		protocol.options.SeparationParams,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal the message: %w", err)
	}

	return marshalledMessage, nil
}
