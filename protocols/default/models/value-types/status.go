package defaultProtocolModelValueTypes

import (
	"errors"
)

type Status struct {
	rawValue []byte
}

func NewStatus(rawValue []byte) (Status, error) {
	if len(rawValue) == 0 {
		return Status{}, errors.New("status cannot be empty")
	}

	value := Status{
		rawValue: rawValue,
	}
	return value, nil
}

func (value Status) ToBytes() []byte {
	return value.rawValue
}
