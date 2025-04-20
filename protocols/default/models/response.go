package defaultProtocolModels

import (
	"fmt"

	"github.com/samber/mo"
	defaultProtocolModelValueTypes "github.com/thewizardplusplus/go-tcp-server/protocols/default/models/value-types"
)

type Response struct {
	status  defaultProtocolModelValueTypes.Status
	headers mo.Option[defaultProtocolModelValueTypes.Headers]
	body    mo.Option[defaultProtocolModelValueTypes.Body]
}

func NewResponseFromMessage(message Message) (Response, error) {
	status, err := defaultProtocolModelValueTypes.NewStatus(
		message.Introduction().ToBytes(),
	)
	if err != nil {
		return Response{}, fmt.Errorf("unable to construct the status: %w", err)
	}

	model := Response{
		status:  status,
		headers: message.Headers(),
		body:    message.Body(),
	}
	return model, nil
}

func (model Response) Status() defaultProtocolModelValueTypes.Status {
	return model.status
}

func (model Response) Headers() mo.Option[defaultProtocolModelValueTypes.Headers] { //nolint:lll
	return model.headers
}

func (model Response) Body() mo.Option[defaultProtocolModelValueTypes.Body] {
	return model.body
}

func (model Response) ToMessage() (Message, error) {
	introduction, err := defaultProtocolModelValueTypes.NewIntroduction(
		model.status.ToBytes(),
	)
	if err != nil {
		return Message{}, fmt.Errorf("unable to construct the introduction: %w", err)
	}

	messageBuilder := NewMessageBuilder().
		SetIntroduction(introduction)

	if headers, isPresent := model.headers.Get(); isPresent {
		messageBuilder.SetHeaders(headers)
	}

	if body, isPresent := model.body.Get(); isPresent {
		messageBuilder.SetBody(body)
	}

	message, err := messageBuilder.Build()
	if err != nil {
		return Message{}, fmt.Errorf("unable to build the message: %w", err)
	}

	return message, nil
}
