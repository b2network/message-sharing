package message

import (
	"bsquared.network/message-sharing-applications/internal/utils/ethereum/event"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/shopspring/decimal"
)

var (
	MessageCallName = "message#call"
	MessageCallHash = crypto.Keccak256([]byte("Call(uint256,uint256,address,uint256,address,bytes)"))
)

type MessageCall struct {
	FromChainId     int64           `json:"from_chain_id"`
	FromId          decimal.Decimal `json:"from_id"`
	FromSender      string          `json:"from_sender"`
	ToChainId       int64           `json:"to_chain_id"`
	ContractAddress string          `json:"contract_address"`
	Bytes           string          `json:"bytes"`
}

func (*MessageCall) Name() string {
	return MessageCallName
}

func (*MessageCall) EventHash() common.Hash {
	return common.BytesToHash(MessageCallHash)
}

func (t *MessageCall) ToObj(data string) error {
	err := json.Unmarshal([]byte(data), &t)
	if err != nil {
		return err
	}
	return nil
}

func (*MessageCall) Data(log types.Log) (string, error) {
	transfer := &MessageCall{
		FromChainId:     event.DataToInt64(log, 0),
		FromId:          event.DataToDecimal(log, 1, 0),
		FromSender:      event.DataToAddress(log, 2).Hex(),
		ToChainId:       event.DataToInt64(log, 3),
		ContractAddress: event.DataToAddress(log, 4).Hex(),
		Bytes:           event.DataToBytes(log, 5),
	}
	data, err := event.ToJson(transfer)
	if err != nil {
		return "", err
	}
	return data, nil
}
