package message

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/storyicon/sigverify"
	"math/big"
	"strings"
)

const LoginTypedData = `{
    "types": {
        "EIP712Domain": [
            {
                "name": "name",
                "type": "string"
            },
            {
                "name": "version",
                "type": "string"
            },
            {
                "name": "chainId",
                "type": "uint256"
            }
        ],
        "Login": [
            {
                "name": "account",
                "type": "address"
            },
            {
                "name": "timestamp",
                "type": "uint256"
            }
        ]
    },
    "domain": {
        "name": "MessageSharingLogin",
        "version": "1",
        "chainId": "%d"
    },
    "primaryType": "Login",
    "message": {
        "account": "%s",
        "timestamp": "%d"
    }
}`

const MessageSendTypedData = `{
    "types": {
        "EIP712Domain": [
            {
                "name": "name",
                "type": "string"
            },
            {
                "name": "version",
                "type": "string"
            },
            {
                "name": "chainId",
                "type": "uint256"
            },
            {
                "name": "verifyingContract",
                "type": "address"
            }
        ],
        "Send": [
            {
                "name": "from_chain_id",
                "type": "uint256"
            },
            {
                "name": "from_id",
                "type": "uint256"
            },
            {
                "name": "from_sender",
                "type": "address"
            },
            {
                "name": "to_chain_id",
                "type": "uint256"
            },
            {
                "name": "to_business_contract",
                "type": "address"
            },
            {
                "name": "to_message",
                "type": "bytes"
            }
        ]
    },
    "domain": {
        "name": "B2MessageSharing",
        "version": "1",
        "chainId": "%d",
        "verifyingContract": "%s"
    },
    "primaryType": "Send",
    "message": {
        "from_chain_id": "%d",
        "from_id": "%s",
        "from_sender": "%s",
        "to_chain_id": "%d",
        "to_business_contract": "%s",
        "to_message": "%s"
    }
}`

func SignLogin(chainId int64, account string, timestamp int64, key *ecdsa.PrivateKey) (string, error) {
	_data := fmt.Sprintf(LoginTypedData, chainId, account, timestamp)
	//fmt.Println("_data", _data)
	var typedData apitypes.TypedData
	if err := json.Unmarshal([]byte(_data), &typedData); err != nil {
		return "", errors.WithStack(err)
	}
	_, originHash, err := sigverify.HashTypedData(typedData)
	//fmt.Println("msgHash", common.Bytes2Hex(msgHash))
	//fmt.Println("originHash", common.Bytes2Hex(originHash))
	if err != nil {
		return "", errors.WithStack(err)
	}
	sig, err := crypto.Sign(originHash, key)
	if err != nil {
		return "", errors.WithStack(err)
	}
	if sig[64] == 0 {
		sig[64] = 27
	} else if sig[64] == 1 {
		sig[64] = 28
	}
	return "0x" + common.Bytes2Hex(sig), nil
}

func VerifyLogin(chainId int64, account string, timestamp int64, signer, signature string) (bool, error) {
	if !strings.HasPrefix(signature, "0x") {
		signature = "0x" + signature
	}
	_signature, err := hexutil.Decode(signature)
	if err != nil {
		return false, errors.WithStack(err)
	}
	if _signature[64] == 27 || _signature[64] == 28 {
		_signature[64] = _signature[64] - 27
	}
	//fmt.Println("_signature:", hexutil.Encode(_signature))
	_data := fmt.Sprintf(LoginTypedData, chainId, account, timestamp)
	//fmt.Println("data:", data)
	var typedData apitypes.TypedData
	if err := json.Unmarshal([]byte(_data), &typedData); err != nil {
		return false, errors.WithStack(err)
	}
	verify, err := sigverify.VerifyTypedDataSignatureEx(
		common.HexToAddress(signer),
		typedData,
		_signature,
	)
	if err != nil || !verify {
		fmt.Errorf("Verify signature error: %v\n", err)
		return false, errors.New("Verify signature failed")
	}
	return true, nil
}

func SignMessageSend(chainId int64, messageContract string, fromChainId int64, fromId *big.Int, fromSender string, toChainId int64, contractAddress string, data string, key *ecdsa.PrivateKey) (string, error) {
	_data := fmt.Sprintf(MessageSendTypedData, chainId, messageContract, fromChainId, fromId.Text(10), fromSender, toChainId, contractAddress, data)
	//fmt.Println("_data", _data)
	var typedData apitypes.TypedData
	if err := json.Unmarshal([]byte(_data), &typedData); err != nil {
		return "", errors.WithStack(err)
	}
	_, originHash, err := sigverify.HashTypedData(typedData)
	//fmt.Println("msgHash", common.Bytes2Hex(msgHash))
	//fmt.Println("originHash", common.Bytes2Hex(originHash))
	if err != nil {
		return "", errors.WithStack(err)
	}
	sig, err := crypto.Sign(originHash, key)
	if err != nil {
		return "", errors.WithStack(err)
	}
	if sig[64] == 0 {
		sig[64] = 27
	} else if sig[64] == 1 {
		sig[64] = 28
	}
	return "0x" + common.Bytes2Hex(sig), nil
}

func VerifyMessageSend(chainId int64, messageContract string, fromChainId int64, fromId *big.Int, fromSender string, toChainId int64, contractAddress string, data string, signer, signature string) (bool, error) {
	if !strings.HasPrefix(signature, "0x") {
		signature = "0x" + signature
	}
	_signature, err := hexutil.Decode(signature)
	if err != nil {
		return false, errors.WithStack(err)
	}
	if _signature[64] == 27 || _signature[64] == 28 {
		_signature[64] = _signature[64] - 27
	}
	//fmt.Println("_signature:", hexutil.Encode(_signature))
	_data := fmt.Sprintf(MessageSendTypedData, chainId, messageContract, fromChainId, fromId.Text(10), fromSender, toChainId, contractAddress, data)
	//fmt.Println("data:", data)
	var typedData apitypes.TypedData
	if err := json.Unmarshal([]byte(_data), &typedData); err != nil {
		return false, errors.WithStack(err)
	}
	verify, err := sigverify.VerifyTypedDataSignatureEx(
		common.HexToAddress(signer),
		typedData,
		_signature,
	)
	if err != nil || !verify {
		fmt.Errorf("Verify signature error: %v\n", err)
		return false, errors.New("Verify signature failed")
	}
	return true, nil
}

func Send(fromChainId int64, fromId *big.Int, fromSender string, contractAddress string, toBytes string, signatures []string) []byte {
	// function send(uint256 from_chain_id, uint256 from_id, address from_sender, address contract_address, bytes calldata data, bytes[] calldata signatures) external
	Method := crypto.Keccak256([]byte("send(uint256,uint256,address,address,bytes,bytes[])"))[:4]
	FromChainId := common.BytesToHash(big.NewInt(fromChainId).Bytes()).Bytes()
	FromId := common.BytesToHash(fromId.Bytes()).Bytes()
	FromSender := common.BytesToHash(common.HexToAddress(fromSender).Bytes()).Bytes()
	ContractAddress := common.BytesToHash(common.HexToAddress(contractAddress).Bytes()).Bytes()

	ToBytes := common.FromHex(toBytes)
	ToBytesDataOffset := common.BytesToHash(big.NewInt(192).Bytes()).Bytes()
	ToBytesDataLength := common.BytesToHash(big.NewInt(int64(len(ToBytes))).Bytes()).Bytes()

	if len(ToBytes)%32 > 0 {
		ToBytes = append(ToBytes, make([]byte, 32-len(ToBytes)%32)...)
	}

	SignaturesDataOffset := common.BytesToHash(big.NewInt(int64(224 + len(ToBytes))).Bytes()).Bytes()
	SignaturesDataLength := common.BytesToHash(big.NewInt(int64(len(signatures))).Bytes()).Bytes()

	var streamOffsets []byte
	var streamData []byte
	streamIndex := int64(32 * len(signatures))
	for _, _signature := range signatures {
		signatureDataOffset := common.BytesToHash(big.NewInt(streamIndex).Bytes()).Bytes()
		streamOffsets = append(streamOffsets, signatureDataOffset...)

		signature := common.FromHex(_signature)
		signatureDataLength := common.BytesToHash(big.NewInt(int64(len(signature))).Bytes()).Bytes()
		streamData = append(streamData, signatureDataLength...)
		if len(signature)%32 > 0 {
			signature = append(signature, make([]byte, 32-len(signature)%32)...)
		}
		streamData = append(streamData, signature...)
		streamIndex = int64(32*len(signatures) + len(streamData))
	}

	var stream []byte
	stream = append(stream, Method...)
	stream = append(stream, FromChainId...)
	stream = append(stream, FromId...)
	stream = append(stream, FromSender...)
	stream = append(stream, ContractAddress...)
	stream = append(stream, ToBytesDataOffset...)
	stream = append(stream, SignaturesDataOffset...)
	stream = append(stream, ToBytesDataLength...)
	stream = append(stream, ToBytes...)
	stream = append(stream, SignaturesDataLength...)
	stream = append(stream, streamOffsets...)
	stream = append(stream, streamData...)
	return stream
}

func EncodeSendData(txId string, fromAddress string, toAddress string, amount decimal.Decimal) []byte {
	TxId := common.HexToHash(txId).Bytes()

	ToAddress := common.BytesToHash(common.HexToAddress(toAddress).Bytes()).Bytes()
	Amount := common.BytesToHash(amount.BigInt().Bytes()).Bytes()

	FromAddress := []byte(fromAddress)
	FromAddressOffset := common.BytesToHash(big.NewInt(128).Bytes()).Bytes()
	FromAddressLength := common.BytesToHash(big.NewInt(int64(len(FromAddress))).Bytes()).Bytes()
	if len(FromAddress)%32 > 0 {
		FromAddress = append(FromAddress, make([]byte, 32-len(FromAddress)%32)...)
	}

	var stream []byte
	stream = append(stream, TxId...)
	stream = append(stream, FromAddressOffset...)
	stream = append(stream, ToAddress...)
	stream = append(stream, Amount...)
	stream = append(stream, FromAddressLength...)
	stream = append(stream, FromAddress...)

	return stream
}
