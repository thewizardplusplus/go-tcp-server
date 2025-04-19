package defaultProtocolModelValueTypes

type Body struct {
	rawValue []byte
}

func NewBody(rawValue []byte) Body {
	if rawValue == nil {
		rawValue = []byte{}
	}

	return Body{
		rawValue: rawValue,
	}
}

func (value Body) ToBytes() []byte {
	return value.rawValue
}
