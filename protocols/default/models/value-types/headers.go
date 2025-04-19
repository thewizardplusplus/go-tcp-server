package defaultProtocolModelValueTypes

type Headers struct {
	rawValue map[HeaderKey]HeaderValue
}

func NewHeaders(rawValue map[HeaderKey]HeaderValue) Headers {
	if rawValue == nil {
		rawValue = map[HeaderKey]HeaderValue{}
	}

	return Headers{
		rawValue: rawValue,
	}
}

func (value Headers) ToMap() map[HeaderKey]HeaderValue {
	return value.rawValue
}
