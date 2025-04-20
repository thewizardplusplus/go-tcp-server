package separatorBasedProtocol

import (
	"bytes"
	"encoding/hex"
	"fmt"
)

func EscapeSeparators(data []byte, separators [][]byte) []byte {
	data = bytes.ReplaceAll(data, []byte("%"), []byte("%25"))

	for _, separator := range separators {
		escapedSeparator := make([]byte, 0, 3*len(separator))
		for _, separatorByte := range separator {
			escapedSeparatorByte := make([]byte, 2)
			hex.Encode(escapedSeparatorByte, []byte{separatorByte})

			escapedSeparator = append(escapedSeparator, '%')
			escapedSeparator = append(escapedSeparator, escapedSeparatorByte...)
		}

		data = bytes.ReplaceAll(data, separator, escapedSeparator)
	}

	return data
}

func UnescapeSeparators(data []byte) ([]byte, error) {
	unescapedData := make([]byte, 0, len(data))
	for index := 0; index < len(data); {
		if data[index] != '%' {
			unescapedData = append(unescapedData, data[index])
			index++

			continue
		}

		if index+2 >= len(data) {
			return nil, fmt.Errorf("not enough bytes to decode at position %d", index)
		}

		unescapedByte := make([]byte, 1)
		escapedByte := data[index+1 : index+3]
		if _, err := hex.Decode(unescapedByte, escapedByte); err != nil {
			return nil, fmt.Errorf(
				"unable to decode the byte sequence 0x%x at position %d: %w",
				escapedByte,
				index,
				err,
			)
		}

		unescapedData = append(unescapedData, unescapedByte...)
		index += 3
	}

	return unescapedData, nil
}
