package defaultProtocolModels

import (
	"github.com/samber/mo"
	defaultProtocolModelValueTypes "github.com/thewizardplusplus/go-tcp-server/protocols/default/models/value-types"
)

type Message struct {
	introduction defaultProtocolModelValueTypes.Introduction
	headers      mo.Option[defaultProtocolModelValueTypes.Headers]
	body         mo.Option[defaultProtocolModelValueTypes.Body]
}

func (model Message) Introduction() defaultProtocolModelValueTypes.Introduction { //nolint:lll
	return model.introduction
}

func (model Message) Headers() mo.Option[defaultProtocolModelValueTypes.Headers] { //nolint:lll
	return model.headers
}

func (model Message) Body() mo.Option[defaultProtocolModelValueTypes.Body] {
	return model.body
}
