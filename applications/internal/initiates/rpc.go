package initiates

import (
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/ethereum/go-ethereum/ethclient"
)

func InitEthereumRpc(rpcUrl string) (*ethclient.Client, error) {
	rpc, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return nil, err
	}
	return rpc, nil
}

func InitBitcoinRpc(rpcUrl string, btcUser string, btcPass string, disableTLS bool) (*rpcclient.Client, error) {
	rpc, err := rpcclient.New(&rpcclient.ConnConfig{
		Host:         rpcUrl,
		User:         btcUser,
		Pass:         btcPass,
		HTTPPostMode: true,       // Bitcoin core only supports HTTP POST mode
		DisableTLS:   disableTLS, // Bitcoin core does not provide TLS by default
	}, nil)
	if err != nil {
		return nil, err
	}
	return rpc, nil
}
