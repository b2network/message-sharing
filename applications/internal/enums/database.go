package enums

type ChainType int64

const (
	ChainTypeUnknown ChainType = iota
	ChainTypeEVM
	ChainTypeUTXO
)

type TaskStatus int64

const (
	TaskStatusUnknown TaskStatus = iota
	TaskStatusPending
	TaskStatusInvalid
	TaskStatusDone
)

type MessageType int64

const (
	MessageTypeUnknown MessageType = iota
	MessageTypeCall
	MessageTypeSend
)

type MessageStatus int64

const (
	MessageStatusUnknown MessageStatus = iota
	MessageStatusValidating
	MessageStatusPending
	MessageStatusBroadcast
	MessageStatusValid
	MessageStatusInvalid
)

type DepositStatus int64

const (
	DepositStatusUnknown DepositStatus = iota
	DepositStatusPending
	DepositStatusValid
	DepositStatusInvalid
)

type SignatureStatus int64

const (
	SignatureStatusPending SignatureStatus = iota
	SignatureStatusBroadcast
	SignatureStatusSuccess
	SignatureStatusFailed
	SignatureStatusInvalid
)
