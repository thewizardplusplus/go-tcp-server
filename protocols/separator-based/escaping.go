package separatorBasedTCPServerProtocol

import (
	"bytes"
	"encoding/hex"
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
