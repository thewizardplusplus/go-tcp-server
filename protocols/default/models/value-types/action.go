package defaultProtocolModelValueTypes

import (
	"errors"
)

type Action struct {
	rawValue []byte
}

func NewAction(rawValue []byte) (Action, error) {
	if len(rawValue) == 0 {
		return Action{}, errors.New("action cannot be empty")
	}

	value := Action{
		rawValue: rawValue,
	}
	return value, nil
}

func (value Action) ToBytes() []byte {
	return value.rawValue
}
