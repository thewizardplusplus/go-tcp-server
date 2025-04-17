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

type SeparationParams struct {
	MessageSeparator        []byte
	MessagePartSeparator    []byte
	HeaderSeparator         []byte
	HeaderKeyValueSeparator []byte
}

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

	parsedHeaders := make(map[string][]byte)
	if headers := messageParts[1]; len(headers) != 0 {
		headerPairs := bytes.Split(headers, params.HeaderSeparator)
		for headerPairIndex, headerPair := range headerPairs {
			separatorIndex := bytes.Index(headerPair, params.HeaderKeyValueSeparator)
			if separatorIndex == -1 {
				return Message{}, fmt.Errorf(
					"header #%d has no key-value separator",
					headerPairIndex,
				)
			}

			escapedKey := headerPair[:separatorIndex]
			if len(escapedKey) == 0 {
				return Message{}, errors.New("header key cannot be empty")
			}

			key, err := UnescapeSeparators(escapedKey)
			if err != nil {
				return Message{}, fmt.Errorf(
					"unable to unescape separators in the header key: %w",
					err,
				)
			}

			escapedValue := headerPair[separatorIndex+1:]
			if len(escapedValue) == 0 {
				return Message{}, errors.New("header value cannot be empty")
			}

			value, err := UnescapeSeparators(escapedValue)
			if err != nil {
				return Message{}, fmt.Errorf(
					"unable to unescape separators in the header value: %w",
					err,
				)
			}

			parsedHeaders[hex.EncodeToString(key)] = value
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
		Headers:      parsedHeaders,
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
	for encodedKey, value := range message.Headers {
		key, err := hex.DecodeString(encodedKey)
		if err != nil {
			return nil, fmt.Errorf(
				"unable to decode the header key %q: %w",
				encodedKey,
				err,
			)
		}

		marshalledHeaders = append(marshalledHeaders, bytes.Join(
			[][]byte{
				EscapeSeparators(key, separators),
				EscapeSeparators(value, separators),
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
