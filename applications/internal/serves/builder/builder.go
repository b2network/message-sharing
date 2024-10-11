package builder

import (
	"bsquared.network/message-sharing-applications/internal/config"
	"bsquared.network/message-sharing-applications/internal/enums"
	"bsquared.network/message-sharing-applications/internal/models"
	msg "bsquared.network/message-sharing-applications/internal/utils/ethereum/message"
	"bsquared.network/message-sharing-applications/internal/utils/log"
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	_types "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"math/big"
	"sync"
	"time"
)

type Builder struct {
	rpc    *ethclient.Client
	db     *gorm.DB
	conf   config.Blockchain
	mu     sync.Mutex
	keys   map[string]bool
	logger *log.Logger
}

func NewBuilder(keys []string, conf config.Blockchain, db *gorm.DB, rpc *ethclient.Client, logger *log.Logger) *Builder {
	_keys := make(map[string]bool, 0)
	for _, key := range keys {
		_keys[key] = true
	}
	return &Builder{
		db:     db,
		rpc:    rpc,
		conf:   conf,
		keys:   _keys,
		logger: logger,
	}
}

func (b *Builder) Start() {
	if !b.conf.Status {
		b.logger.Infof("status: %t", b.conf.Status)
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go b.build()
	go b.broadcast()
	<-ctx.Done()
}

func (b *Builder) build() {
	duration := time.Millisecond * time.Duration(b.conf.BlockInterval)
	for {
		list, err := b.pendingCallMessage(b.conf.SignatureWeight, 10)
		if err != nil {
			b.logger.Errorf("et pending call message err: %s", err)
			time.Sleep(duration)
			continue
		}
		if len(list) == 0 {
			b.logger.Infof("Get pending call message length is 0")
			time.Sleep(duration)
			continue
		}

		for _, message := range list {
			err = b.buildMessage(message)
			if err != nil {
				b.logger.Errorf("Handle err: %v, %v", err, message)
			}
		}
		//var wg sync.WaitGroup
		//for _, message := range list {
		//	wg.Add(1)
		//	go func(_wg *sync.WaitGroup, message models.Message) {
		//		defer _wg.Done()
		//		err = b.buildMessage(message)
		//		if err != nil {
		//			b.logger.Errorf("Handle err: %v, %v", err, message)
		//		}
		//	}(&wg, message)
		//}
		//wg.Wait()
	}
}

func (b *Builder) buildMessage(message models.Message) error {
	UserAddress, UserKey, err := b.BorrowAccount()
	if err != nil {
		b.logger.Errorf("borrow account err: %s\n", err)
		return errors.WithStack(err)
	}
	defer b.mu.Unlock()
	b.logger.Errorf("UserAddress: %s", UserAddress)

	//lock, err := b.LockUser(UserAddress, time.Minute*2)
	//if err != nil {
	//	b.logger.Errorf("lock err: %s\n", err)
	//	return errors.WithStack(err)
	//}
	//if !lock {
	//	b.logger.Infof("load result: %v\n", lock)
	//	return errors.WithStack(err)
	//}
	//defer b.UnlockUser(UserAddress)

	gasPrice, err := b.GasPrice()
	if err != nil {
		return errors.WithStack(err)
	}
	b.logger.Debugf("gasPrice: %v\n", gasPrice)

	toAddress := common.HexToAddress(message.ToMessageBridge)
	b.logger.Debugf("toAddress: %v\n", toAddress)

	var signatures []string
	err = b.db.Model(models.MessageSignature{}).Select([]string{"signature"}).Where("message_id=?", message.Id).Find(&signatures).Error
	if err != nil {
		return errors.WithStack(err)
	}
	//err = json.Unmarshal([]byte(message.Signatures), &signatures)
	//if err != nil {
	//	return errors.WithStack(err)
	//}
	b.logger.Infof("FromChainId: %d, FromId: %s, FromSender: %s, ToContractAddress: %s, ToBytes: %s", message.FromChainId, common.HexToHash(message.FromId).Big().Text(16), message.FromSender, message.ToContractAddress, message.ToBytes)
	data := msg.Send(message.FromChainId, common.HexToHash(message.FromId).Big(), message.FromSender, message.ToContractAddress, message.ToBytes, signatures)
	b.logger.Debugf("data: %x\n", data)
	gasLimit, err := b.rpc.EstimateGas(context.Background(), ethereum.CallMsg{
		From:     common.HexToAddress(UserAddress),
		To:       &toAddress,
		GasPrice: gasPrice,
		Value:    big.NewInt(0),
		Data:     data,
	})
	if err != nil {
		b.logger.Errorf("Get gasLimit err: %s\n", err)
		return errors.WithStack(err)
	}
	b.logger.Debugf("gasLimit: %v\n", gasLimit)
	err = b.db.Transaction(func(tx *gorm.DB) error {
		// nonce
		nonce, err := b.GetNonce(UserAddress)
		if err != nil {
			b.logger.Errorf("get nonce err: %s\n", err)
			return errors.WithStack(err)
		}
		b.logger.Debugf("nonce: %v\n", nonce)
		// signTx
		_signature, err := b.SignTx(UserAddress, UserKey, nonce, toAddress.Hex(), big.NewInt(0), gasLimit, gasPrice, data, b.conf.ChainId)
		if err != nil {
			b.logger.Errorf("sign tx err: %s\n", err)
			return errors.WithStack(err)
		}
		_txHash := crypto.Keccak256Hash(_signature)
		b.logger.Debugf("txHash: %s, signature: %s\n", _txHash, hex.EncodeToString(_signature))

		// create signature
		err = b.CreateSignature(tx, message.ToChainId, message.FromId, UserAddress, int64(nonce), enums.MessageTypeSend, hex.EncodeToString(data), decimal.Zero, hex.EncodeToString(_signature), _txHash.Hex())
		if err != nil {
			b.logger.Errorf("create signature err: %s\n", err)
			return errors.WithStack(err)
		}
		// broadcast
		_signatures, err := json.Marshal(signatures)
		if err != nil {
			return errors.WithStack(err)
		}
		message.Signatures = string(_signatures)
		message.Status = enums.MessageStatusBroadcast
		err = tx.Save(&message).Error
		if err != nil {
			b.logger.Errorf("build message err: %s", err)
			return errors.WithStack(err)
		}
		b.logger.Infof("build message success")
		return nil
	})
	if err != nil {
		b.logger.Errorf("build message err: %s", err)
		return errors.WithStack(err)
	}
	return nil
}

func (b *Builder) pendingCallMessage(weight int64, limit int) ([]models.Message, error) {
	var list []models.Message
	err := b.db.Where("`to_chain_id`=? AND `type`=? AND `signatures_count`>=? AND `status`=?",
		b.conf.ChainId, enums.MessageTypeCall, weight, enums.MessageStatusPending).Limit(limit).Find(&list).Error
	if err != nil {
		b.logger.Errorf("get message err: %s", err)
		return nil, err
	}
	return list, nil
}

func (b *Builder) SignTx(accountAddress string, accountKey *ecdsa.PrivateKey, nonce uint64, toAddress string, value *big.Int, gasLimit uint64, gasPrice *big.Int, bytecode []byte, chainID int64) ([]byte, error) {
	_signature, err := b._signTx(accountAddress, accountKey, nonce, toAddress, value, gasLimit, gasPrice, bytecode, chainID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return _signature, nil
}

func (b *Builder) _signTx(accountAddress string, accountKey *ecdsa.PrivateKey, nonce uint64, toAddress string, value *big.Int, gasLimit uint64, gasPrice *big.Int, bytecode []byte, chainID int64) ([]byte, error) {
	if crypto.PubkeyToAddress(accountKey.PublicKey) != common.HexToAddress(accountAddress) {
		return nil, errors.New(" address and index do not match ")
	}
	tx := types.NewTransaction(
		nonce,
		common.HexToAddress(toAddress),
		value,
		gasLimit,
		gasPrice,
		bytecode,
	)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(chainID)), accountKey)
	if err != nil {
		return nil, err
	}
	ts := types.Transactions{signedTx}
	var rawTxBytes bytes.Buffer
	ts.EncodeIndex(0, &rawTxBytes)
	return rawTxBytes.Bytes(), nil
}

func (b *Builder) GetNonce(userAddress string) (uint64, error) {
	var signature models.Signature
	err := b.db.Raw("SELECT * FROM signatures WHERE `status`!=? AND `address`=? ORDER BY nonce DESC FOR UPDATE ",
		enums.SignatureStatusInvalid, userAddress).First(&signature).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	} else if err == gorm.ErrRecordNotFound {
		nonce, err := b.rpc.NonceAt(context.Background(), common.HexToAddress(userAddress), nil)
		if err != nil {
			return 0, err
		}
		return nonce, nil
	} else {
		if signature.Status != enums.SignatureStatusSuccess && signature.Status != enums.SignatureStatusFailed {
			return 0, errors.New("The current user has pending transactions ")
		}
		nonce, err := b.rpc.NonceAt(context.Background(), common.HexToAddress(userAddress), nil)
		if err != nil {
			return 0, err
		}
		return nonce, nil
	}
}

func (b *Builder) CreateSignature(tx *gorm.DB, chainId int64, referId string, address string, nonce int64, signatureType enums.MessageType, data string, value decimal.Decimal, signature string, txHash string) error {
	err := tx.Create(&models.Signature{
		ChainId:   chainId,
		ReferId:   referId,
		Address:   address,
		Nonce:     nonce,
		Type:      signatureType,
		Data:      data,
		Value:     value,
		Signature: signature,
		Status:    enums.SignatureStatusPending,
		Blockchain: models.Blockchain{
			EventId:     0,
			BlockTime:   0,
			BlockNumber: 0,
			LogIndex:    0,
			TxHash:      txHash,
		},
	}).Error
	if err != nil {
		return err
	}
	return nil
}

func (b *Builder) GasPrice() (*big.Int, error) {
	gasPrice, err := b.rpc.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return gasPrice, nil
}

func (b *Builder) LockUser(key string, duration time.Duration) (bool, error) {
	//return l.Cache.Client.SetNX(context.Background(), key, true, duration).Result()
	return false, nil
}

func (b *Builder) UnlockUser(key string) error {
	//_, err := l.Cache.Client.Del(context.Background(), key).Result()
	//return err
	return nil
}

func (b *Builder) BorrowAccount() (string, *ecdsa.PrivateKey, error) {
	if b.mu.TryLock() {
		var pk string
		for _pk, _ := range b.keys {
			pk = _pk
			break
		}
		_key, err := crypto.ToECDSA(common.FromHex(pk))
		if err != nil {
			return "", nil, err
		}
		return crypto.PubkeyToAddress(_key.PublicKey).Hex(), _key, nil
	} else {
		return "", nil, errors.New("account is locked")
	}
}

//func (b *Builder) getKeyByAddress(accountAddress string) (string, error) {
//	if key, ok := l.DataMap.SenderMap[accountAddress]; ok {
//		return key, nil
//	} else {
//		return "", errors.New("account not found")
//	}
//	return "", nil
//}

func (b *Builder) broadcast() {
	duration := time.Millisecond * time.Duration(b.conf.BlockInterval)
	for {
		var signatures []models.Signature
		err := b.db.Where("`chain_id`=? AND `status`=?", b.conf.ChainId, enums.SignatureStatusPending).Order("id").Limit(100).Find(&signatures).Error
		if err != nil {
			time.Sleep(duration)
			continue
		}
		if len(signatures) == 0 {
			time.Sleep(duration)
			continue
		}

		var wg sync.WaitGroup
		for _, signature := range signatures {
			wg.Add(1)
			go func(_wg *sync.WaitGroup, signature models.Signature) {
				defer _wg.Done()
				err = b.broadcastSignature(signature)
				if err != nil {
					b.logger.Errorf("Broadcast signature err[%d]: %s\n", signature.Id, err)
				}
			}(&wg, signature)
		}
		wg.Wait()
	}
}

func (b *Builder) broadcastSignature(signature models.Signature) error {
	err := b._broadcast(signature.Signature)
	if err != nil && err.Error() != "nonce too low" && err.Error() != "already known" {
		return errors.WithStack(err)
	} else if err != nil && err.Error() == "nonce too low" {
		// _, err := ctx.RPC.TransactionReceipt(context.Background(), common.HexToHash(signature.TxHash))
		// if err != nil {
		//	log.Errorf("Get TransactionReceipt err[%d]: %s\n", signature.Id, err)
		//	return err
		// }
		err = b.db.Model(&models.Signature{}).Where("id = ?", signature.Id).Update("status", enums.SignatureStatusBroadcast).Error
		if err != nil {
			return errors.WithStack(err)
		}
	} else if err != nil && err.Error() == "already known" {
		err = b.db.Model(&models.Signature{}).Where("id = ?", signature.Id).Update("status", enums.SignatureStatusBroadcast).Error
		if err != nil {
			return errors.WithStack(err)
		}
	} else {
		err = b.db.Model(&models.Signature{}).Where("id = ?", signature.Id).Update("status", enums.SignatureStatusBroadcast).Error
		if err != nil {
			return errors.WithStack(err)
		}
		return nil
	}
	return nil
}

func (b *Builder) _broadcast(signature string) error {
	rawTxBytes, err := hex.DecodeString(signature)
	if err != nil {
		b.logger.Errorf("Decode signature err[%d]: %s\n", signature, err)
		return err
	}
	tx := new(_types.Transaction)
	rlp.DecodeBytes(rawTxBytes, &tx)
	err = b.rpc.SendTransaction(context.Background(), tx)
	if err != nil {
		return err
	}
	return nil
}
