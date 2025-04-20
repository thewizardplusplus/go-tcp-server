package defaultProtocolModels

import (
	"fmt"

	"github.com/samber/mo"
	defaultProtocolModelValueTypes "github.com/thewizardplusplus/go-tcp-server/protocols/default/models/value-types"
)

type Request struct {
	action  defaultProtocolModelValueTypes.Action
	headers mo.Option[defaultProtocolModelValueTypes.Headers]
	body    mo.Option[defaultProtocolModelValueTypes.Body]
}

func NewRequestFromMessage(message Message) (Request, error) {
	action, err := defaultProtocolModelValueTypes.NewAction(
		message.Introduction().ToBytes(),
	)
	if err != nil {
		return Request{}, fmt.Errorf("unable to construct the action: %w", err)
	}

	model := Request{
		action:  action,
		headers: message.Headers(),
		body:    message.Body(),
	}
	return model, nil
}

func (model Request) Action() defaultProtocolModelValueTypes.Action {
	return model.action
}

func (model Request) Headers() mo.Option[defaultProtocolModelValueTypes.Headers] { //nolint:lll
	return model.headers
}

func (model Request) Body() mo.Option[defaultProtocolModelValueTypes.Body] {
	return model.body
}

func (model Request) ToMessage() (Message, error) {
	introduction, err := defaultProtocolModelValueTypes.NewIntroduction(
		model.action.ToBytes(),
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
