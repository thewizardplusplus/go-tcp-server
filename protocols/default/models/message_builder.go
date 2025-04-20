package defaultProtocolModels

import (
	"errors"

	"github.com/samber/mo"
	defaultProtocolModelValueTypes "github.com/thewizardplusplus/go-tcp-server/protocols/default/models/value-types"
)

type MessageBuilder struct {
	introduction mo.Option[defaultProtocolModelValueTypes.Introduction]
	headers      mo.Option[defaultProtocolModelValueTypes.Headers]
	body         mo.Option[defaultProtocolModelValueTypes.Body]
}

func NewMessageBuilder() *MessageBuilder {
	return &MessageBuilder{}
}

func (builder *MessageBuilder) SetIntroduction(
	value defaultProtocolModelValueTypes.Introduction,
) *MessageBuilder {
	builder.introduction = mo.Some(value)
	return builder
}

func (builder *MessageBuilder) SetHeaders(
	value defaultProtocolModelValueTypes.Headers,
) *MessageBuilder {
	builder.headers = mo.Some(value)
	return builder
}

func (builder *MessageBuilder) SetBody(
	value defaultProtocolModelValueTypes.Body,
) *MessageBuilder {
	builder.body = mo.Some(value)
	return builder
}

func (builder MessageBuilder) Build() (Message, error) {
	var errs []error

	introduction, isPresent := builder.introduction.Get()
	if !isPresent {
		errs = append(errs, errors.New("introduction is required"))
	}

	if headers, isPresent :=
		builder.headers.Get(); isPresent && len(headers.ToMap()) == 0 {
		builder.headers = mo.None[defaultProtocolModelValueTypes.Headers]()
	}

	if body, isPresent :=
		builder.body.Get(); isPresent && len(body.ToBytes()) == 0 {
		builder.body = mo.None[defaultProtocolModelValueTypes.Body]()
	}

	if len(errs) > 0 {
		return Message{}, errors.Join(errs...)
	}

	model := Message{
		introduction: introduction,
		headers:      builder.headers,
		body:         builder.body,
	}
	return model, nil
}
