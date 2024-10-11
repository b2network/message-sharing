package vo

import (
	"bsquared.network/message-sharing-applications/internal/enums"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/ethereum/go-ethereum/ethclient"
)

type MessageWrap struct {
	MessageType enums.P2PMessageType
	Data        string
}

type Login struct {
	ChainId   int64
	Account   string
	Timestamp int64
	Signature string
}

type LoginResult struct {
	ChainId int64
	Account string
	Result  bool
}

type Message struct {
	MessageId           int64
	ChainId             int64
	FromMessageContract string
	FromChainId         int64
	FromId              string
	FromSender          string
	ToChainId           int64
	ToMessageContract   string
	ToContractAddress   string
	Data                string
	TxHash              string
	LogIndex            int64
}

type MessageSignature struct {
	MessageId           int64
	ChainId             int64
	FromMessageContract string
	FromChainId         int64
	FromId              string
	FromSender          string
	ToChainId           int64
	ToMessageContract   string
	ToContractAddress   string
	Data                string
	Signature           string
}

type RpcClient struct {
	EthRpc *ethclient.Client
	BtcRpc *rpcclient.Client
}
