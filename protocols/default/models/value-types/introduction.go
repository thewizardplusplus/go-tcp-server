package defaultProtocolModelValueTypes

import (
	"errors"
)

type Introduction struct {
	rawValue []byte
}

func NewIntroduction(rawValue []byte) (Introduction, error) {
	if len(rawValue) == 0 {
		return Introduction{}, errors.New("introduction cannot be empty")
	}

	value := Introduction{
		rawValue: rawValue,
	}
	return value, nil
}

func (value Introduction) ToBytes() []byte {
	return value.rawValue
}
