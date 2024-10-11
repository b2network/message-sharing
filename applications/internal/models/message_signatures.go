package models

type MessageSignature struct {
	Base
	MessageId int64  `json:"message_id"`
	Signer    string `json:"signer"`
	Signature string `json:"signature"`
}

func (MessageSignature) TableName() string {
	return "`message_signatures`"
}
