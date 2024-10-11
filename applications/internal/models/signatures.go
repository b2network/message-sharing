package models

import (
	"bsquared.network/message-sharing-applications/internal/enums"
	"github.com/shopspring/decimal"
)

type Signature struct {
	Base
	ChainId   int64                 `json:"chain_id"`
	ReferId   string                `json:"refer_id"`
	Type      enums.MessageType     `json:"type"`
	Address   string                `json:"address"`
	Nonce     int64                 `json:"nonce"`
	Data      string                `json:"data"`
	Value     decimal.Decimal       `json:"value"`
	Signature string                `json:"signature"`
	Status    enums.SignatureStatus `json:"status"`
	Blockchain
}

func (Signature) TableName() string {
	return "`signatures`"
}
