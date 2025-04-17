package separatorBasedTCPServerProtocol

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"slices"
)

const (
	MessagePartCount = 3
)

type Message struct {
	Introduction []byte
	Headers      map[string][]byte
	Body         []byte
}

func ParseMessage(data []byte, params SeparationParams) (Message, error) {
	messageParts := bytes.SplitN(
		data,
		params.MessagePartSeparator,
		MessagePartCount,
	)
	if len(messageParts) < MessagePartCount {
		return Message{}, fmt.Errorf(
			"invalid message part count: %d",
			len(messageParts),
		)
	}

	escapedIntroduction := messageParts[0]
	if len(escapedIntroduction) == 0 {
		return Message{}, errors.New("introduction cannot be empty")
	}

	introduction, err := UnescapeSeparators(escapedIntroduction)
	if err != nil {
		return Message{}, fmt.Errorf(
			"unable to unescape separators in the introduction: %w",
			err,
		)
	}

	headers := make(map[string][]byte)
	if marshalledHeaders := messageParts[1]; len(marshalledHeaders) != 0 {
		for marshalledHeaderIndex, marshalledHeader := range bytes.Split(
			marshalledHeaders,
			params.HeaderSeparator,
		) {
			separatorIndex := bytes.Index(
				marshalledHeader,
				params.HeaderKeyValueSeparator,
			)
			if separatorIndex == -1 {
				return Message{}, fmt.Errorf(
					"header #%d has no key-value separator",
					marshalledHeaderIndex,
				)
			}

			escapedHeaderKey := marshalledHeader[:separatorIndex]
			if len(escapedHeaderKey) == 0 {
				return Message{}, errors.New("header key cannot be empty")
			}

			headerKey, err := UnescapeSeparators(escapedHeaderKey)
			if err != nil {
				return Message{}, fmt.Errorf(
					"unable to unescape separators in the header key: %w",
					err,
				)
			}

			escapedHeaderValue := marshalledHeader[separatorIndex+1:]
			if len(escapedHeaderValue) == 0 {
				return Message{}, errors.New("header value cannot be empty")
			}

			headerValue, err := UnescapeSeparators(escapedHeaderValue)
			if err != nil {
				return Message{}, fmt.Errorf(
					"unable to unescape separators in the header value: %w",
					err,
				)
			}

			headers[hex.EncodeToString(headerKey)] = headerValue
		}
	}

	body, err := UnescapeSeparators(messageParts[2])
	if err != nil {
		return Message{}, fmt.Errorf(
			"unable to unescape separators in the body: %w",
			err,
		)
	}

	message := Message{
		Introduction: introduction,
		Headers:      headers,
		Body:         body,
	}
	return message, nil
}

func MarshalMessage(message Message, params SeparationParams) ([]byte, error) {
	separators := [][]byte{
		params.MessageSeparator,
		params.MessagePartSeparator,
		params.HeaderSeparator,
		params.HeaderKeyValueSeparator,
	}

	marshalledHeaders := make([][]byte, 0, len(message.Headers))
	for encodedHeaderKey, headerValue := range message.Headers {
		headerKey, err := hex.DecodeString(encodedHeaderKey)
		if err != nil {
			return nil, fmt.Errorf("unable to decode the header key: %w", err)
		}

		marshalledHeaders = append(marshalledHeaders, bytes.Join(
			[][]byte{
				EscapeSeparators(headerKey, separators),
				EscapeSeparators(headerValue, separators),
			},
			params.HeaderKeyValueSeparator,
		))
	}
	slices.SortStableFunc(marshalledHeaders, func(a []byte, b []byte) int {
		return bytes.Compare(a, b)
	})

	marshalledMessage := bytes.Join(
		[][]byte{
			EscapeSeparators(message.Introduction, separators),
			bytes.Join(marshalledHeaders, params.HeaderSeparator),
			EscapeSeparators(message.Body, separators),
		},
		params.MessagePartSeparator,
	)
	return marshalledMessage, nil
}
