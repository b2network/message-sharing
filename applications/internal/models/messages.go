package models

import "bsquared.network/message-sharing-applications/internal/enums"

type Message struct {
	Base
	ChainId           int64               `json:"chain_id"`
	Type              enums.MessageType   `json:"type"`
	FromChainId       int64               `json:"from_chain_id"`
	FromMessageBridge string              `json:"from_message_bridge"`
	FromSender        string              `json:"from_sender"`
	FromId            string              `json:"from_id"`
	ToChainId         int64               `json:"to_chain_id"`
	ToMessageBridge   string              `json:"to_message_bridge"`
	ToContractAddress string              `json:"to_contract_address"`
	ToBytes           string              `json:"to_bytes"`
	Signatures        string              `json:"signatures"`
	SignaturesCount   int64               `json:"signatures_count"`
	Status            enums.MessageStatus `json:"status"`
	Blockchain
}

func (Message) TableName() string {
	return "`messages`"
}
