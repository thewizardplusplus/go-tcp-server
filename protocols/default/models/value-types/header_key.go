package defaultProtocolModelValueTypes

import (
	"encoding/hex"
	"errors"
	"fmt"
)

type HeaderKey struct {
	encodedRawValue string
}

func NewHeaderKey(rawValue []byte) (HeaderKey, error) {
	if len(rawValue) == 0 {
		return HeaderKey{}, errors.New("header key cannot be empty")
	}

	value := HeaderKey{
		encodedRawValue: hex.EncodeToString(rawValue),
	}
	return value, nil
}

func MustNewHeaderKey(rawValue []byte) HeaderKey {
	value, err := NewHeaderKey(rawValue)
	if err != nil {
		panic(fmt.Sprintf(
			"tcpServerProtocolModelValueTypes.MustNewHeaderKey(): %s",
			err,
		))
	}

	return value
}

func (value HeaderKey) ToBytes() ([]byte, error) {
	rawValue, err := hex.DecodeString(value.encodedRawValue)
	if err != nil {
		return nil, fmt.Errorf("unable to decode the header key: %w", err)
	}

	return rawValue, nil
}
