package defaultProtocolModels

import (
	"errors"

	"github.com/samber/mo"
	defaultProtocolModelValueTypes "github.com/thewizardplusplus/go-tcp-server/protocols/default/models/value-types"
)

type ResponseBuilder struct {
	status  mo.Option[defaultProtocolModelValueTypes.Status]
	headers mo.Option[defaultProtocolModelValueTypes.Headers]
	body    mo.Option[defaultProtocolModelValueTypes.Body]
}

func NewResponseBuilder() *ResponseBuilder {
	return &ResponseBuilder{}
}

func (builder *ResponseBuilder) SetStatus(
	value defaultProtocolModelValueTypes.Status,
) *ResponseBuilder {
	builder.status = mo.Some(value)
	return builder
}

func (builder *ResponseBuilder) SetHeaders(
	value defaultProtocolModelValueTypes.Headers,
) *ResponseBuilder {
	builder.headers = mo.Some(value)
	return builder
}

func (builder *ResponseBuilder) SetBody(
	value defaultProtocolModelValueTypes.Body,
) *ResponseBuilder {
	builder.body = mo.Some(value)
	return builder
}

func (builder ResponseBuilder) Build() (Response, error) {
	var errs []error

	status, isPresent := builder.status.Get()
	if !isPresent {
		errs = append(errs, errors.New("status is required"))
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
		return Response{}, errors.Join(errs...)
	}

	model := Response{
		status:  status,
		headers: builder.headers,
		body:    builder.body,
	}
	return model, nil
}
