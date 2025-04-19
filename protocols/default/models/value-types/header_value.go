package defaultProtocolModelValueTypes

import (
	"errors"
	"fmt"
)

type HeaderValue struct {
	rawValue []byte
}

func NewHeaderValue(rawValue []byte) (HeaderValue, error) {
	if len(rawValue) == 0 {
		return HeaderValue{}, errors.New("header value cannot be empty")
	}

	value := HeaderValue{
		rawValue: rawValue,
	}
	return value, nil
}

func MustNewHeaderValue(rawValue []byte) HeaderValue {
	value, err := NewHeaderValue(rawValue)
	if err != nil {
		panic(fmt.Sprintf(
			"tcpServerProtocolModelValueTypes.MustNewHeaderValue(): %s",
			err,
		))
	}

	return value
}

func (value HeaderValue) ToBytes() []byte {
	return value.rawValue
}
