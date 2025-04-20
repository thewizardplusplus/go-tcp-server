package separatorBasedTCPServerProtocol

import (
	"bytes"
	"fmt"
	"slices"

	defaultProtocolModels "github.com/thewizardplusplus/go-tcp-server/protocols/default/models"
	defaultProtocolModelValueTypes "github.com/thewizardplusplus/go-tcp-server/protocols/default/models/value-types"
)

const (
	messagePartCount = 3
)

type MessageFormat struct {
	options SeparationParams
}

func NewMessageFormat(options SeparationParams) MessageFormat {
	return MessageFormat{
		options: options,
	}
}

func (format MessageFormat) ParseMessage(
	data []byte,
) (defaultProtocolModels.Message, error) {
	messageParts := bytes.SplitN(
		data,
		format.options.MessagePartSeparator,
		messagePartCount,
	)
	if len(messageParts) < messagePartCount {
		return defaultProtocolModels.Message{}, fmt.Errorf(
			"invalid message part count: %d",
			len(messageParts),
		)
	}

	rawIntroduction, err := UnescapeSeparators(messageParts[0])
	if err != nil {
		return defaultProtocolModels.Message{}, fmt.Errorf(
			"unable to unescape separators in the introduction: %w",
			err,
		)
	}

	introduction, err :=
		defaultProtocolModelValueTypes.NewIntroduction(rawIntroduction)
	if err != nil {
		return defaultProtocolModels.Message{}, fmt.Errorf(
			"unable to construct the introduction: %w",
			err,
		)
	}

	rawHeaders :=
		make(map[defaultProtocolModelValueTypes.HeaderKey]defaultProtocolModelValueTypes.HeaderValue) //nolint:lll
	if marshalledHeaders := messageParts[1]; len(marshalledHeaders) != 0 {
		for marshalledHeaderIndex, marshalledHeader := range bytes.Split(
			marshalledHeaders,
			format.options.HeaderSeparator,
		) {
			separatorIndex := bytes.Index(
				marshalledHeader,
				format.options.HeaderKeyValueSeparator,
			)
			if separatorIndex == -1 {
				return defaultProtocolModels.Message{}, fmt.Errorf(
					"header #%d has no key-value separator",
					marshalledHeaderIndex,
				)
			}

			rawHeaderKey, err := UnescapeSeparators(marshalledHeader[:separatorIndex])
			if err != nil {
				return defaultProtocolModels.Message{}, fmt.Errorf(
					"unable to unescape separators in the header key: %w",
					err,
				)
			}

			headerKey, err := defaultProtocolModelValueTypes.NewHeaderKey(rawHeaderKey)
			if err != nil {
				return defaultProtocolModels.Message{}, fmt.Errorf(
					"unable to construct the header key: %w",
					err,
				)
			}

			rawHeaderValue, err :=
				UnescapeSeparators(marshalledHeader[separatorIndex+1:])
			if err != nil {
				return defaultProtocolModels.Message{}, fmt.Errorf(
					"unable to unescape separators in the header value: %w",
					err,
				)
			}

			headerValue, err :=
				defaultProtocolModelValueTypes.NewHeaderValue(rawHeaderValue)
			if err != nil {
				return defaultProtocolModels.Message{}, fmt.Errorf(
					"unable to construct the header value: %w",
					err,
				)
			}

			rawHeaders[headerKey] = headerValue
		}
	}

	rawBody, err := UnescapeSeparators(messageParts[2])
	if err != nil {
		return defaultProtocolModels.Message{}, fmt.Errorf(
			"unable to unescape separators in the body: %w",
			err,
		)
	}

	message, err := defaultProtocolModels.NewMessageBuilder().
		SetIntroduction(introduction).
		SetHeaders(defaultProtocolModelValueTypes.NewHeaders(rawHeaders)).
		SetBody(defaultProtocolModelValueTypes.NewBody(rawBody)).
		Build()
	if err != nil {
		return defaultProtocolModels.Message{}, fmt.Errorf(
			"unable to build the message: %w",
			err,
		)
	}

	return message, nil
}

func (format MessageFormat) MarshalMessage(
	message defaultProtocolModels.Message,
) ([]byte, error) {
	separators := [][]byte{
		format.options.MessageSeparator,
		format.options.MessagePartSeparator,
		format.options.HeaderSeparator,
		format.options.HeaderKeyValueSeparator,
	}

	marshalledHeaders :=
		make([][]byte, 0, len(message.Headers().OrEmpty().ToMap()))
	for headerKey, headerValue := range message.Headers().OrEmpty().ToMap() {
		rawHeaderKey, err := headerKey.ToBytes()
		if err != nil {
			return nil, fmt.Errorf("unable to convert the header key to bytes: %w", err)
		}

		marshalledHeaders = append(marshalledHeaders, bytes.Join(
			[][]byte{
				EscapeSeparators(rawHeaderKey, separators),
				EscapeSeparators(headerValue.ToBytes(), separators),
			},
			format.options.HeaderKeyValueSeparator,
		))
	}
	slices.SortStableFunc(marshalledHeaders, func(a []byte, b []byte) int {
		return bytes.Compare(a, b)
	})

	marshalledMessage := bytes.Join(
		[][]byte{
			EscapeSeparators(message.Introduction().ToBytes(), separators),
			bytes.Join(marshalledHeaders, format.options.HeaderSeparator),
			EscapeSeparators(message.Body().OrEmpty().ToBytes(), separators),
		},
		format.options.MessagePartSeparator,
	)
	return marshalledMessage, nil
}
