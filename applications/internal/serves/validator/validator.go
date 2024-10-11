package validator

import (
	"bsquared.network/message-sharing-applications/internal/config"
	"bsquared.network/message-sharing-applications/internal/enums"
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
	crypto_ "github.com/ethereum/go-ethereum/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/multiformats/go-multiaddr"
	"github.com/pkg/errors"
	"time"
)

const (
	PROTOCOL = "/chat/1.0.0"
)

type Validator struct {
	particle config.Particle
	conf     config.Blockchain
	host     host.Host
	rw       *bufio.ReadWriter
	pk       *ecdsa.PrivateKey
	logger   *log.Logger
	client   *vo.RpcClient
}

func NewValidator(pk *ecdsa.PrivateKey, host host.Host, logger *log.Logger, client *vo.RpcClient, particle config.Particle, conf config.Blockchain) *Validator {
	return &Validator{
		conf:     conf,
		particle: particle,
		host:     host,
		pk:       pk,
		logger:   logger,
		client:   client,
	}
}

func (v *Validator) Start() {
	if !v.conf.Status {
		v.logger.Infof("status: %t", v.conf.Status)
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go v.connect()
	<-ctx.Done()
}

func (v *Validator) connect() {
	for {
		if v.rw != nil {
			time.Sleep(time.Second * 10)
			continue
		}
		time.Sleep(time.Second * 2)
		v.logger.Infof("connect ...")
		multiaddr, err := multiaddr.NewMultiaddr(v.conf.Endpoint)
		if err != nil {
			v.logger.Errorf("new multiaddr err: %s", err)
			continue
		}
		info, err := peer.AddrInfoFromP2pAddr(multiaddr)
		if err != nil {
			v.logger.Errorf("addr info from p2p addr err: %s", err)
			continue
		}
		v.host.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)
		v.host.RemoveStreamHandler(PROTOCOL)
		s, err := v.host.NewStream(context.Background(), info.ID, PROTOCOL)
		if err != nil {
			v.logger.Errorf("new stream err: %s", err)
			continue
		}
		v.rw = bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
		go v.accept()
		err = v.login()
		if err != nil {
			v.logger.Errorf("login err: %s", err)
		}
	}
}

func (v *Validator) login() error {
	v.logger.Infof("validator login ...")
	timestamp := time.Now().Unix()
	account := crypto_.PubkeyToAddress(v.pk.PublicKey).Hex()
	signature, err := message.SignLogin(v.conf.ChainId, account, timestamp, v.pk)
	if err != nil {
		v.logger.Errorf("sign login err: %s", err)
		return err
	}
	value, err := json.Marshal(&vo.Login{
		ChainId:   v.conf.ChainId,
		Account:   account,
		Timestamp: timestamp,
		Signature: signature,
	})
	if err != nil {
		v.logger.Errorf("value json marshal err: %s", err)
		return err
	}
	valueWrap, err := json.Marshal(&vo.MessageWrap{
		MessageType: enums.P2PMessageTypeLogin,
		Data:        string(value),
	})
	if err != nil {
		v.logger.Errorf("valueWrap json marshal err: %s", err)
		return err
	}
	_, err = v.rw.WriteString(fmt.Sprintf("%s\n", string(valueWrap)))
	if err != nil {
		v.logger.Errorf("send login message err: %s", err)
		return err
	}
	err = v.rw.Flush()
	if err != nil {
		v.logger.Errorf("flush login message err: %s", err)
		return err
	}
	v.logger.Infof("login success")
	return nil
}

func (v *Validator) accept() {
	for {
		v.logger.Infof("validator accept ...")
		msg, err := v.rw.ReadString('\n')
		if err != nil {
			v.logger.Errorf("validator read err: %s", err)
			v.rw = nil
			return
		}
		var messageWrap vo.MessageWrap
		err = json.Unmarshal([]byte(msg), &messageWrap)
		if err != nil {
			v.logger.Errorf("messageWrap json unmarshal err: %s", err)
			continue
		}
		if messageWrap.MessageType == enums.P2PMessageTypeLogin {

		} else if messageWrap.MessageType == enums.P2PMessageTypeProposal {
			var message vo.Message
			err = json.Unmarshal([]byte(messageWrap.Data), &message)
			if err != nil {
				v.logger.Errorf("message json unmarshal err: %s", err)
				continue
			}
			go func() {
				err = v.handleMessage(message)
				if err != nil {
					v.logger.Errorf("handle message signature err: %s", err)
				}
			}()
		}
	}
}

func (v *Validator) handleMessage(msg vo.Message) error {
	if v.client.EthRpc != nil {
		verify, err := tx.VerifyEthTx(v.client.EthRpc, msg.TxHash, msg.LogIndex, msg.FromMessageContract, msg.FromChainId, msg.FromId, msg.FromSender, msg.ToChainId, msg.ToContractAddress, msg.Data)
		if err != nil {
			v.logger.Errorf("verify eth tx err: %s", err)
			return err
		}
		if !verify {
			return errors.New("verify message failed")
		}
	} else if v.client.BtcRpc != nil {
		var chainParams *chaincfg.Params
		if v.conf.Mainnet {
			chainParams = &chaincfg.MainNetParams
		} else {
			chainParams = &chaincfg.TestNet3Params
		}
		verify, err := tx.VerifyBtcTx(v.client.BtcRpc, chainParams, v.particle, msg.FromMessageContract, msg.TxHash, msg.FromId, msg.Data)
		if err != nil {
			v.logger.Errorf("verify btc tx err: %s", err)
			return err
		}
		if !verify {
			return errors.New("verify message failed")
		}
	} else {
		return errors.New("rpc invalid")
	}
	fmt.Printf("data :%s\n", msg.Data)
	fmt.Printf("validator :%s\n", crypto_.PubkeyToAddress(v.pk.PublicKey))
	signature, err := message.SignMessageSend(msg.ChainId, msg.ToMessageContract, msg.FromChainId, common.HexToHash(msg.FromId).Big(), msg.FromSender, msg.ToChainId, msg.ToContractAddress, msg.Data, v.pk)
	if err != nil {
		v.logger.Errorf("validator sign err: %s", err)
		return err
	}
	value, err := json.Marshal(&vo.MessageSignature{
		MessageId:           msg.MessageId,
		ChainId:             msg.ChainId,
		FromMessageContract: msg.FromMessageContract,
		FromChainId:         msg.FromChainId,
		FromId:              msg.FromId,
		FromSender:          msg.FromSender,
		ToChainId:           msg.ToChainId,
		ToMessageContract:   msg.ToMessageContract,
		ToContractAddress:   msg.ToContractAddress,
		Data:                msg.Data,
		Signature:           signature,
	})
	if err != nil {
		v.logger.Errorf("value json marshal err: %s", err)
		return err
	}
	valueWrap, err := json.Marshal(&vo.MessageWrap{
		MessageType: enums.P2PMessageTypeSign,
		Data:        string(value),
	})
	if err != nil {
		v.logger.Errorf("valueWrap json marshal err: %s", err)
		return err
	}
	v.rw.WriteString(fmt.Sprintf("%s\n", string(valueWrap)))
	err = v.rw.Flush()
	if err != nil {
		v.logger.Errorf("validator flush err: %s", err)
		return err
	}
	return nil
}
