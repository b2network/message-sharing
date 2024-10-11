package tx

import (
	"bsquared.network/message-sharing-applications/internal/config"
	"bsquared.network/message-sharing-applications/internal/types"
	"bsquared.network/message-sharing-applications/internal/utils/aa"
	"bsquared.network/message-sharing-applications/internal/utils/ethereum/event"
	"bsquared.network/message-sharing-applications/internal/utils/ethereum/message"
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

var (
	ErrParsePkScript            = errors.New("parse pkscript err")
	ErrDecodeListenAddress      = errors.New("decode listen address err")
	ErrTargetConfirmations      = errors.New("target confirmation number was not reached")
	ErrParsePubKey              = errors.New("parse pubkey failed, not found pubkey or nonsupport ")
	ErrParsePkScriptNullData    = errors.New("parse pkscript null data err")
	ErrParsePkScriptNotNullData = errors.New("parse pkscript not null data err")
)

func VerifyEthTx(rpc *ethclient.Client, txHash string, logIndex int64, fromMessageAddress string,
	fromChainId int64, fromId string, fromSender string, toChainId int64, toContractAddress string, toBytes string) (bool, error) {
	tx, err := rpc.TransactionReceipt(context.Background(), common.HexToHash(txHash))
	if err != nil {
		return false, err
	}
	for _, log := range tx.Logs {
		if log.Index == uint(logIndex) {
			if common.HexToAddress(fromMessageAddress) == log.Address &&
				fromChainId == event.DataToInt64(*log, 0) &&
				common.HexToHash(fromId).Big().Text(16) == event.DataToDecimal(*log, 1, 0).BigInt().Text(16) &&
				common.HexToAddress(fromSender) == event.DataToAddress(*log, 2) &&
				toChainId == event.DataToInt64(*log, 3) &&
				common.HexToAddress(toContractAddress) == event.DataToAddress(*log, 4) &&
				toBytes == event.DataToBytes(*log, 5) {
				return true, nil
			}
		}
	}
	return false, nil
}

func VerifyBtcTx(rpc *rpcclient.Client, chainParams *chaincfg.Params, particle config.Particle, listenAddress string, txHash string, fromId string, data string) (bool, error) {
	_txHash, err := chainhash.NewHashFromStr(txHash[2:])
	if err != nil {
		return false, err
	}
	tx, err := rpc.GetRawTransaction(_txHash)
	if err != nil {
		return false, err
	}
	txResult := tx.MsgTx()
	_listenAddress, err := btcutil.DecodeAddress(listenAddress, chainParams)
	if err != nil {
		return false, err
	}
	var totalValue int64
	var depositAddress string
	for _, v := range txResult.TxOut {
		pkAddress, err := parseAddress(chainParams, v.PkScript)
		if err != nil {
			if errors.Is(err, ErrParsePkScript) {
				continue
			}
			if errors.Is(err, ErrParsePkScriptNullData) {
				nullData, err := parseNullData(v.PkScript)
				if err != nil {
					continue
				}
				evmAddress, err := parseEvmAddress(nullData)
				if err != nil {
					continue
				}
				fmt.Println("evmAddress:", evmAddress)
				if depositAddress == "" {
					depositAddress = evmAddress
				}
			} else {
				return false, err
			}
		}
		if pkAddress == _listenAddress.EncodeAddress() {
			totalValue += v.Value
		}
	}
	fromAddress, err := parseFromAddress(rpc, chainParams, txResult)
	if err != nil {
		return false, err
	}
	if len(fromAddress) == 0 {
		return false, errors.New("fromAddress invalid")
	}
	if depositAddress == "" {
		_depositAddress, err := getAADepositAddress(particle, fromAddress[0].Address)
		if err != nil {
			return false, err
		}
		depositAddress = _depositAddress
	}
	_data := message.EncodeSendData(txResult.TxHash().String(), fromAddress[0].Address, depositAddress, decimal.New(totalValue, 0))
	if common.HexToHash(fromId) == common.HexToHash(txHash) && data == "0x"+hex.EncodeToString(_data) {
		return true, nil
	}
	return false, nil
}

func parseAddress(chainParams *chaincfg.Params, pkScript []byte) (string, error) {
	pk, err := txscript.ParsePkScript(pkScript)
	if err != nil {
		scriptClass := txscript.GetScriptClass(pkScript)
		if scriptClass == txscript.NullDataTy {
			return "", ErrParsePkScriptNullData
		}
		return "", fmt.Errorf("%w:%s", ErrParsePkScript, err.Error())

	}
	if pk.Class() == txscript.NullDataTy {
		return "", ErrParsePkScriptNullData
	}

	//  encodes the script into an address for the given chain.
	//var chainParams *chaincfg.Params
	//if l.config.Mainnet {
	//	chainParams = &chaincfg.MainNetParams
	//} else {
	//	chainParams = &chaincfg.TestNet3Params
	//}
	pkAddress, err := pk.Address(chainParams)
	if err != nil {
		return "", fmt.Errorf("PKScript to address err:%w", err)
	}
	return pkAddress.EncodeAddress(), nil
}

func parseNullData(pkScript []byte) (string, error) {
	if !txscript.IsNullData(pkScript) {
		return "", ErrParsePkScriptNotNullData
	}
	return hex.EncodeToString(pkScript[1:]), nil
}

func parseEvmAddress(nullData string) (string, error) {
	decodeNullData, err := hex.DecodeString(nullData)
	if err != nil {
		return "", err
	}
	evmAddress := bytes.TrimSpace(decodeNullData[1:])
	if common.IsHexAddress(string(evmAddress)) {
		return string(evmAddress), nil
	}
	return "", nil
}

func parseFromAddress(rpc *rpcclient.Client, chainParams *chaincfg.Params, txResult *wire.MsgTx) (fromAddress []types.BitcoinFrom, err error) {
	for _, vin := range txResult.TxIn {
		// get prev tx hash
		prevTxID := vin.PreviousOutPoint.Hash
		vinResult, err := rpc.GetRawTransaction(&prevTxID)
		if err != nil {
			return nil, fmt.Errorf("vin get raw transaction err:%w", err)
		}
		if len(vinResult.MsgTx().TxOut) == 0 {
			return nil, fmt.Errorf("vin txOut is null")
		}
		vinPKScript := vinResult.MsgTx().TxOut[vin.PreviousOutPoint.Index].PkScript
		//  script to address
		vinPkAddress, err := parseAddress(chainParams, vinPKScript)
		if err != nil {
			//b.logger.Errorw("vin parse address", "error", err)
			if errors.Is(err, ErrParsePkScript) || errors.Is(err, ErrParsePkScriptNullData) {
				continue
			}
			return nil, err
		}

		fromAddress = append(fromAddress, types.BitcoinFrom{
			Address: vinPkAddress,
			Type:    types.BitcoinFromTypeBtc,
		})
	}
	return fromAddress, nil
}

func getAADepositAddress(particle config.Particle, btcFrom string) (string, error) {
	evmAddress, err := aa.BitcoinAddressToEthAddress(particle.AAPubKeyAPI, btcFrom,
		particle.Url, particle.ChainId, particle.ProjectUuid, particle.ProjectKey)
	if err != nil {
		return "", err
	}
	return evmAddress, nil
}
