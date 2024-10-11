package proposer

import (
	"bsquared.network/message-sharing-applications/internal/config"
	"bsquared.network/message-sharing-applications/internal/enums"
	"bsquared.network/message-sharing-applications/internal/models"
	"bsquared.network/message-sharing-applications/internal/utils/ethereum/message"
	"bsquared.network/message-sharing-applications/internal/utils/log"
	"bsquared.network/message-sharing-applications/internal/utils/tx"
	"bsquared.network/message-sharing-applications/internal/vo"
	"bufio"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/common"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"math/big"
	"time"
)

const (
	PROTOCOL = "/chat/1.0.0"
)

type Proposer struct {
	conf     config.Blockchain
	particle config.Particle
	host     host.Host
	db       *gorm.DB
	pk       *ecdsa.PrivateKey
	client   *vo.RpcClient
	logger   *log.Logger
	smap     map[string]common.Address
	rws      map[string]*bufio.ReadWriter
}

func NewProposer(pk *ecdsa.PrivateKey, host host.Host, db *gorm.DB, client *vo.RpcClient, logger *log.Logger, conf config.Blockchain, particle config.Particle) *Proposer {
	return &Proposer{
		conf:     conf,
		particle: particle,
		pk:       pk,
		host:     host,
		db:       db,
		client:   client,
		logger:   logger,
		smap:     make(map[string]common.Address, 0),
		rws:      make(map[string]*bufio.ReadWriter, 0),
	}
}

func (p *Proposer) Start() {
	if !p.conf.Status {
		p.logger.Infof("status: %t", p.conf.Status)
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go p.listen()
	go p.proposal()
	go p.submit()
	p.logger.Infof("proposer-node start success ,node-port: %d ,node-id: %s", p.conf.NodePort, p.host.ID())
	<-ctx.Done()
}

func (p *Proposer) listen() {
	p.host.SetStreamHandler(PROTOCOL, func(s network.Stream) {
		rw, ok := p.rws[s.ID()]
		if !ok {
			rw = bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
			p.rws[s.ID()] = rw
		}
		p.logger.Infof("connect id: %s", s.ID())
		go p.accept(s.ID(), rw)
	})
}

func (p *Proposer) accept(sid string, rw *bufio.ReadWriter) {
	for {
		p.logger.Infof("accept ....")
		if rw == nil {
			time.Sleep(time.Second * 1)
			return
		}
		msg, err := rw.ReadString('\n')
		if err != nil {
			p.logger.Errorf("validator read err: %s", err)
			rw = nil
			continue
		}
		p.logger.Infof("msg: %s", msg)
		var messageWrap vo.MessageWrap
		err = json.Unmarshal([]byte(msg), &messageWrap)
		if err != nil {
			p.logger.Errorf("json unmarshal err: %s", err)
			continue
		}
		p.logger.Infof("messageWrap: %v", messageWrap)

		if messageWrap.MessageType == enums.P2PMessageTypeLogin {
			var l vo.Login
			err = json.Unmarshal([]byte(messageWrap.Data), &l)
			if err != nil {
				p.logger.Errorf("err: %s", err)
				continue
			}
			go func() {
				err = p.handleLogin(sid, l)
				if err != nil {
					p.logger.Errorf("handle login err: %s", err)
				}
			}()
		} else if messageWrap.MessageType == enums.P2PMessageTypeSign {
			var messageSignature vo.MessageSignature
			err = json.Unmarshal([]byte(messageWrap.Data), &messageSignature)
			if err != nil {
				p.logger.Errorf("json unmarshal err: %s", err)
				continue
			}
			go func() {
				err = p.handleSignature(sid, messageSignature)
				if err != nil {
					p.logger.Errorf("handle message signature err: %s", err)
				}
			}()
		}
	}
}

func (p *Proposer) submit() {
	for {
		time.Sleep(time.Second * 3)
		result := p.db.Model(models.Message{}).
			Where("status=? AND signatures_count>=?", enums.MessageStatusValidating, p.conf.SignatureWeight).
			Update("status", enums.MessageStatusPending)
		if result.Error != nil {
			p.logger.Errorf("submit message err: %s", result.Error)
			continue
		}
		p.logger.Infof("submit message count: %d", result.RowsAffected)
	}
}

func (p *Proposer) proposal() {
	for {
		time.Sleep(time.Second * 3)
		list, err := p.getValidatingMessages(p.conf.SignatureWeight, 10)
		if err != nil {
			p.logger.Errorf("validating call message", err)
			continue
		}
		if len(list) == 0 {
			p.logger.Info("message length is 0")
			continue
		}
		for _, message := range list {
			err = p.send(message)
			if err != nil {
				p.logger.Errorf("send message err: %s", err)
			}
		}
	}
}

func (p *Proposer) handleSignature(sid string, messageSignature vo.MessageSignature) error {
	signer, ok := p.smap[sid]
	if !ok {
		return errors.New("no login")
	}
	fromId := big.NewInt(0).SetBytes(common.FromHex(messageSignature.FromId))
	verify, err := message.VerifyMessageSend(messageSignature.ChainId, messageSignature.ToMessageContract, messageSignature.FromChainId, fromId, messageSignature.FromSender, messageSignature.ToChainId, messageSignature.ToContractAddress, messageSignature.Data, signer.Hex(), messageSignature.Signature)
	if err != nil {
		return err
	}
	if !verify {
		return errors.New("invalid signature")
	}
	err = p.db.Transaction(func(tx *gorm.DB) error {
		var count int64
		err = tx.Model(models.MessageSignature{}).Where("`message_id`=? AND `signer`=?", messageSignature.MessageId, signer.Hex()).Count(&count).Error
		if err != nil {
			return err
		}
		if count > 0 {
			return nil
		}
		err = tx.Create(&models.MessageSignature{
			MessageId: messageSignature.MessageId,
			Signer:    signer.Hex(),
			Signature: messageSignature.Signature,
		}).Error
		if err != nil {
			return err
		}
		err = tx.Exec(fmt.Sprintf("UPDATE %s set signatures_count=signatures_count+1 WHERE id =?", models.Message{}.TableName()),
			messageSignature.MessageId).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (p *Proposer) handleLogin(sid string, l vo.Login) error {
	p.logger.Infof("%d#%s#%d login: ", l.ChainId, l.Account, l.Timestamp)
	if l.ChainId != p.conf.ChainId {
		return errors.New("invalid chain id")
	}
	t := time.Now().Unix() - l.Timestamp
	if t < -60 || t > 60 {
		return errors.New("invalid timestamp")
	}
	var valid bool
	for _, validator := range p.conf.Validators {
		if common.HexToAddress(validator) == common.HexToAddress(l.Account) {
			valid = true
			break
		}
	}
	if !valid {
		return errors.New("invalid validator account")
	}
	verify, err := message.VerifyLogin(l.ChainId, l.Account, l.Timestamp, l.Account, l.Signature)
	if err != nil {
		p.logger.Errorf("verify login err: %s", err)
		return err
	}
	p.logger.Infof("verify: %v", verify)
	if verify {
		p.smap[sid] = common.HexToAddress(l.Account)
	}
	return nil
}

func (p *Proposer) send(message models.Message) error {
	if p.client.EthRpc != nil {
		verify, err := tx.VerifyEthTx(p.client.EthRpc, message.TxHash, message.LogIndex, message.FromMessageBridge, message.FromChainId,
			message.FromId, message.FromSender, message.ToChainId, message.ToContractAddress, message.ToBytes)
		if err != nil {
			p.logger.Errorf("verify eth tx err: %s", err)
			return err
		}
		p.logger.Infof("verify eth tx: %t", verify)
		if !verify {
			err = p.db.Model(models.Message{}).
				Where("id=?", message.Id).
				Update("status", enums.MessageStatusInvalid).Error
			if err != nil {
				p.logger.Errorf("update message err: %s", err)
				return err
			}
			return errors.New("verify message failed")
		}
	} else if p.client.BtcRpc != nil {
		var chainParams *chaincfg.Params
		if p.conf.Mainnet {
			chainParams = &chaincfg.MainNetParams
		} else {
			chainParams = &chaincfg.TestNet3Params
		}
		verify, err := tx.VerifyBtcTx(p.client.BtcRpc, chainParams, p.particle, message.FromMessageBridge, message.TxHash, message.FromId, message.ToBytes)
		if err != nil {
			p.logger.Errorf("verify btc tx err: %s", err)
			return err
		}
		p.logger.Infof("verify btc tx: %t", verify)
		if !verify {
			err = p.db.Model(models.Message{}).
				Where("id=?", message.Id).
				Update("status", enums.MessageStatusInvalid).Error
			if err != nil {
				p.logger.Errorf("update message err: %s", err)
				return err
			}
			return errors.New("verify message failed")
		}
	} else {
		return errors.New("rpc invalid")
	}
	proposal := vo.Message{
		MessageId:           message.Id,
		ChainId:             message.ToChainId,
		FromChainId:         message.FromChainId,
		FromMessageContract: message.FromMessageBridge,
		FromId:              message.FromId,
		FromSender:          message.FromSender,
		ToChainId:           message.ToChainId,
		ToMessageContract:   message.ToMessageBridge,
		ToContractAddress:   message.ToContractAddress,
		Data:                message.ToBytes,
		TxHash:              message.TxHash,
		LogIndex:            message.LogIndex,
	}
	value, err := json.Marshal(&proposal)
	if err != nil {
		p.logger.Errorf("json marshal err: %s", err)
		return err
	}
	msg := vo.MessageWrap{
		MessageType: enums.P2PMessageTypeProposal,
		Data:        string(value),
	}
	msgValue, err := json.Marshal(&msg)
	if err != nil {
		p.logger.Errorf("json marshal err: %s", err)
		return err
	}

	for id, writer := range p.rws {
		_, ok := p.smap[id]
		if !ok {
			//p.logger.Errorf("writer error network reset")
			//delete(p.rws, id)
			continue
		}
		_, err = writer.WriteString(fmt.Sprintf("%s\n", string(msgValue)))
		if err != nil && err != network.ErrReset {
			p.logger.Errorf("writer err: %s", err)
			continue
		} else if err == network.ErrReset {
			p.logger.Errorf("writer error network reset")
			delete(p.rws, id)
			continue
		}
		err = writer.Flush()
		if err != nil {
			p.logger.Errorf("writer flush err: %s", err)
			continue
		}
	}
	return nil
}

func (p *Proposer) getValidatingMessages(weight int64, limit int) ([]models.Message, error) {
	var list []models.Message
	err := p.db.Where("`chain_id`=? AND `type`=? AND `status`=? AND signatures_count<?",
		p.conf.ChainId, enums.MessageTypeCall, enums.MessageStatusValidating, weight).Limit(limit).Order("signatures_count").Find(&list).Error
	if err != nil {
		p.logger.Errorf("get message err: %s", err)
		return nil, errors.WithStack(err)
	}
	return list, nil
}
