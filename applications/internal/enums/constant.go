package enums

type P2PMessageType int64

const (
	P2PMessageTypeUnknown P2PMessageType = iota
	P2PMessageTypeLogin
	P2PMessageTypeProposal
	P2PMessageTypeSign
)
