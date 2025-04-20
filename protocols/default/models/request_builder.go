package defaultProtocolModels

import (
	"errors"

	"github.com/samber/mo"
	defaultProtocolModelValueTypes "github.com/thewizardplusplus/go-tcp-server/protocols/default/models/value-types"
)

type RequestBuilder struct {
	action  mo.Option[defaultProtocolModelValueTypes.Action]
	headers mo.Option[defaultProtocolModelValueTypes.Headers]
	body    mo.Option[defaultProtocolModelValueTypes.Body]
}

func NewRequestBuilder() *RequestBuilder {
	return &RequestBuilder{}
}

func (builder *RequestBuilder) SetAction(
	value defaultProtocolModelValueTypes.Action,
) *RequestBuilder {
	builder.action = mo.Some(value)
	return builder
}

func (builder *RequestBuilder) SetHeaders(
	value defaultProtocolModelValueTypes.Headers,
) *RequestBuilder {
	builder.headers = mo.Some(value)
	return builder
}

func (builder *RequestBuilder) SetBody(
	value defaultProtocolModelValueTypes.Body,
) *RequestBuilder {
	builder.body = mo.Some(value)
	return builder
}

func (builder RequestBuilder) Build() (Request, error) {
	var errs []error

	action, isPresent := builder.action.Get()
	if !isPresent {
		errs = append(errs, errors.New("action is required"))
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
		return Request{}, errors.Join(errs...)
	}

	model := Request{
		action:  action,
		headers: builder.headers,
		body:    builder.body,
	}
	return model, nil
}
