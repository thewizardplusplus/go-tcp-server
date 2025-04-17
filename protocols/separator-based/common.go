package separatorBasedTCPServerProtocol

type SeparationParams struct {
	MessageSeparator        []byte
	MessagePartSeparator    []byte
	HeaderSeparator         []byte
	HeaderKeyValueSeparator []byte
}
