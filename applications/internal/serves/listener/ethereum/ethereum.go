package ethereum

import (
	"bsquared.network/message-sharing-applications/internal/config"
	"bsquared.network/message-sharing-applications/internal/enums"
	"bsquared.network/message-sharing-applications/internal/models"
	"bsquared.network/message-sharing-applications/internal/utils/ethereum/event"
	"bsquared.network/message-sharing-applications/internal/utils/ethereum/event/message"
	"bsquared.network/message-sharing-applications/internal/utils/log"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"math/big"
	"strings"
	"sync"
	"time"
)

type EthereumListener struct {
	conf        config.Blockchain
	rpc         *ethclient.Client
	db          *gorm.DB
	logger      *log.Logger
	latestBlock int64
	bridges     map[int64]string
}

func NewListener(bridges map[int64]string, conf config.Blockchain, rpc *ethclient.Client, db *gorm.DB, logger *log.Logger) *EthereumListener {
	return &EthereumListener{
		conf:    conf,
		rpc:     rpc,
		db:      db,
		logger:  logger,
		bridges: bridges,
	}
}

func (l *EthereumListener) Start() {
	if !l.conf.Status {
		l.logger.Infof("status: %t", l.conf.Status)
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go l.syncLastBlock()
	go l.syncTask()
	go l.handEvent()
	go l.confirm()
	<-ctx.Done()
}

func (l *EthereumListener) syncLastBlock() {
	for {
		duration := time.Millisecond * time.Duration(l.conf.BlockInterval)
		latest, err := l.rpc.BlockNumber(context.Background())
		if err != nil {
			l.logger.Errorf("sync latest block error: %s", err.Error())
			time.Sleep(duration)
			continue
		}
		l.latestBlock = int64(latest)
		l.logger.Infof("sync latest block: %d", l.latestBlock)
		time.Sleep(duration)
	}
}

func (l *EthereumListener) syncTask() {
	for {
		duration := time.Millisecond * time.Duration(l.conf.BlockInterval) * 5
		var tasks []models.SyncTask
		err := l.db.Where("`chain_type`=? AND chain_id=? AND status=?", enums.ChainTypeEVM, l.conf.ChainId, enums.TaskStatusPending).Limit(20).Find(&tasks).Error
		if err != nil {
			l.logger.Errorf("task list error: %s\n", err)
			time.Sleep(duration)
			continue
		}
		if len(tasks) == 0 {
			l.logger.Info("task list is empty")
			time.Sleep(duration)
			continue
		}
		wg := sync.WaitGroup{}
		for _, task := range tasks {
			wg.Add(1)
			go func(take models.SyncTask, wg *sync.WaitGroup) {
				defer wg.Done()
				err := l.handleTask(take)
				if err != nil {
					l.logger.Errorf("handle task error: %s", err.Error())
				}
			}(task, &wg)

		}
		wg.Wait()
		time.Sleep(duration)
	}
}

func (l *EthereumListener) handleTask(task models.SyncTask) error {
	start := task.LatestBlock
	if task.StartBlock > start {
		start = task.StartBlock
	}

	if task.EndBlock > 0 && start > task.EndBlock {
		task.Status = enums.TaskStatusDone
		l.db.Save(&task)
		return nil
	}
	end := start
	if task.HandleNum > 0 {
		end = start + task.HandleNum - 1
	}
	if task.EndBlock > 0 && end > task.EndBlock {
		end = task.EndBlock
	}
	if end > l.latestBlock {
		end = l.latestBlock
	}
	if start > end {
		return nil
	}

	//Contracts := strings.Split(task.Contracts, ",")
	//if len(Contracts) == 0 {
	//	//Contracts = l.GetContracts()
	//	//log.Infof("[Handler.SyncTask]  Contracts invalid")
	//	////task.UpdateTime = time.Now()
	//	//task.Status = models.SyncTaskInvalid
	//	//ctx.Db.Save(&task)
	//	//return nil
	//}

	logs, err := l.rpc.FilterLogs(context.Background(), ethereum.FilterQuery{
		FromBlock: big.NewInt(start),
		ToBlock:   big.NewInt(end),
		Topics: [][]common.Hash{
			{
				common.BytesToHash(message.MessageCallHash),
				common.BytesToHash(message.MessageSendHash),
			},
		},
		Addresses: []common.Address{
			common.HexToAddress(l.conf.ListenAddress),
		},
	})
	l.logger.Infof(" start: %d, end: %d\n", start, end)
	if err != nil {
		l.logger.Errorf("[Handler.SyncTask]  FilterLogs error: %v", err)
		return errors.WithStack(err)
	}

	events, err := l.LogsToEvents(logs)
	if err != nil {
		l.logger.Errorf("[Handler.SyncTask]  LogsToEvents error: %v", err)
		return errors.WithStack(err)
	}
	BatchCreateEvents := make([]*models.SyncEvent, 0)
	BatchUpdateEventIds := make([]int64, 0)
	for _, event := range events {
		var one models.SyncEvent
		err = l.db.Select("id").Where("block_number=? AND block_log_indexed=? AND tx_hash=? AND event_hash=?",
			event.BlockNumber, event.BlockLogIndexed, event.TxHash, event.EventHash).First(&one).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			l.logger.Errorf("[Handler.SyncTask]  Get event err: %s\n", err)
			return errors.WithStack(err)
		} else if err == gorm.ErrRecordNotFound {
			BatchCreateEvents = append(BatchCreateEvents, event)
		} else {
			BatchUpdateEventIds = append(BatchUpdateEventIds, one.Id)
		}
	}

	err = l.db.Transaction(func(tx *gorm.DB) error {
		if len(BatchCreateEvents) > 0 {
			err = tx.CreateInBatches(&BatchCreateEvents, 100).Error
			if err != nil {
				l.logger.Errorf("[Handler.SyncEvent]Create SyncEvent err: %s\n", err)
				return errors.WithStack(err)
			}
		}
		if len(BatchUpdateEventIds) > 0 {
			err = tx.Model(models.SyncEvent{}).
				Where("id in ?", BatchUpdateEventIds).
				Update("status", models.EventPending).Error
			if err != nil {
				l.logger.Errorf("[Handler.SyncEvent]Update SyncEvent err: %s\n,", err)
				return errors.WithStack(err)
			}
		}
		task.LatestBlock = end + 1
		err = tx.Save(&task).Error
		if err != nil {
			l.logger.Errorf("[Handler.SyncEvent]Update SyncTask err: %s\n,", err)
			return errors.WithStack(err)
		}
		return nil
	})
	if err != nil {
		l.logger.Errorf("[Handler.SyncEvent]Update SyncTask err: %s\n,", err)
		return errors.WithStack(err)
	}
	return nil
}

func (l *EthereumListener) handEvent() {
	duration := time.Millisecond * time.Duration(l.conf.BlockInterval) * 5
	for {
		var events []models.SyncEvent
		err := l.db.Model(models.SyncEvent{}).Where("`chain_id`=? AND `event_hash` in ? AND `status`=?",
			l.conf.ChainId, []string{
				common.BytesToHash(message.MessageCallHash).Hex(),
				common.BytesToHash(message.MessageSendHash).Hex(),
			}, models.EventPending).
			Limit(500).Find(&events).Error
		if err != nil {
			l.logger.Errorf("[Handler.SyncEvent] err: %s\n", err)
			time.Sleep(duration)
			continue
		}
		if len(events) == 0 {
			l.logger.Errorf("[Handler.SyncEvent] no event\n")
			time.Sleep(duration)
			continue
		}
		valids := make([]int64, 0)
		invalids := make([]int64, 0)
		messages := make([]models.Message, 0)
		handles := make(map[string]bool)

		var Type enums.MessageType
		var FromMessageBridge string
		var FromChainId int64
		var FromSender string
		var FromId string
		var ToChainId int64
		var ToMessageBridge string
		var ToContractAddress string
		var ToBytes string
		var status enums.MessageStatus

		for _, event := range events {
			key := fmt.Sprintf("%s#%d", event.TxHash, event.BlockLogIndexed)
			if handles[key] {
				invalids = append(invalids, event.Id)
				continue
			}

			if event.EventName == message.MessageCallName {
				var messageCall message.MessageCall
				err := (&messageCall).ToObj(event.Data)
				if err != nil {
					l.logger.Errorf("to obj err: %s\n", err.Error())
					time.Sleep(duration)
					continue
				}
				FromMessageBridge = event.ContractAddress
				FromChainId = messageCall.FromChainId
				FromSender = messageCall.FromSender
				FromId = common.BytesToHash(messageCall.FromId.BigInt().Bytes()).Hex()
				ToChainId = messageCall.ToChainId
				ToContractAddress = messageCall.ContractAddress
				ToBytes = messageCall.Bytes
				Type = enums.MessageTypeCall
				status = enums.MessageStatusValidating
			} else if event.EventName == message.MessageSendName {
				var messageSend message.MessageSend
				err := (&messageSend).ToObj(event.Data)
				if err != nil {
					l.logger.Errorf("event to data err: %v, data: %v\n", err, event)
					continue
				}
				FromChainId = messageSend.FromChainId
				FromSender = messageSend.FromSender
				FromId = common.BytesToHash(messageSend.FromId.BigInt().Bytes()).Hex()
				ToChainId = messageSend.ToChainId
				ToContractAddress = messageSend.ContractAddress
				ToMessageBridge = event.ContractAddress
				ToBytes = messageSend.Bytes
				Type = enums.MessageTypeSend
				status = enums.MessageStatusPending
			}

			if FromMessageBridge == "" {
				messageBridge, ok := l.bridges[FromChainId]
				if ok {
					FromMessageBridge = messageBridge
				} else {
					invalids = append(invalids, event.Id)
					continue
				}
			}

			if ToMessageBridge == "" {
				messageBridge, ok := l.bridges[ToChainId]
				if ok {
					ToMessageBridge = messageBridge
				} else {
					invalids = append(invalids, event.Id)
					continue
				}
			}

			var message models.Message
			err = l.db.Where("`tx_hash`=? AND `log_index`=?", event.TxHash, event.BlockLogIndexed).First(&message).Error
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				l.logger.Errorf("query message err: %v, data: %v\n", err, event)
				time.Sleep(duration)
				continue
			} else if errors.Is(err, gorm.ErrRecordNotFound) {
				handles[key] = true
				message = models.Message{
					ChainId:           event.ChainId,
					Type:              Type,
					FromChainId:       FromChainId,
					FromSender:        FromSender,
					FromMessageBridge: FromMessageBridge,
					FromId:            FromId,
					ToChainId:         ToChainId,
					ToMessageBridge:   ToMessageBridge,
					ToContractAddress: ToContractAddress,
					ToBytes:           ToBytes,
					Signatures:        "[]",
					Status:            status,
					Blockchain: models.Blockchain{
						EventId:     event.Id,
						BlockTime:   event.BlockTime,
						BlockNumber: event.BlockNumber,
						LogIndex:    event.BlockLogIndexed,
						TxHash:      event.TxHash,
					},
				}
				messages = append(messages, message)
				valids = append(valids, event.Id)
			} else {
				invalids = append(invalids, event.Id)
			}

		}
		err = l.db.Transaction(func(tx *gorm.DB) error {
			if len(valids) > 0 {
				err = tx.Model(models.SyncEvent{}).Where("id in ?", valids).Update("status", models.EventValid).Error
				if err != nil {
					l.logger.Errorf("update valid Event  err: %v, data: %v\n", err, valids)
					return err
				}
			}
			if len(invalids) > 0 {
				err = tx.Model(models.SyncEvent{}).Where("id in ?", invalids).Update("status", models.EventInvalid).Error
				if err != nil {
					l.logger.Errorf("update invalid Event  err: %v, data: %v\n", err, invalids)
					return err
				}
			}
			if len(messages) > 0 {
				err = tx.CreateInBatches(messages, 100).Error
				if err != nil {
					l.logger.Errorf("create messages err: %v, data: %v\n", err, messages)
					return err
				}
			}
			return nil
		})
		if err != nil {
			l.logger.Errorf("update valid Event  err: %v\n", err)
			time.Sleep(duration)
		}
	}
}

func (l *EthereumListener) LogsToEvents(logs []types.Log) ([]*models.SyncEvent, error) {
	var events []*models.SyncEvent
	blockTimes := make(map[int64]int64)
	for _, vlog := range logs {
		eventHash := event.TopicToHash(vlog, 0)
		contractAddress := vlog.Address

		var eventName string
		var data string
		if eventHash == common.BytesToHash(message.MessageCallHash) {
			eventName = message.MessageCallName
			e := &message.MessageCall{}
			_data, err := e.Data(vlog)
			if err != nil {
				l.logger.Errorf("parse message call err: %v, data: %v\n", err, vlog)
				return nil, err
			}
			data = _data
		} else if eventHash == common.BytesToHash(message.MessageSendHash) {
			eventName = message.MessageSendName
			e := &message.MessageSend{}
			_data, err := e.Data(vlog)
			if err != nil {
				l.logger.Errorf("parse message send err: %v, data: %v\n", err, vlog)
				return nil, err
			}
			data = _data
		}
		blockTime := blockTimes[int64(vlog.BlockNumber)]
		if blockTime == 0 {
			//	blockJson, err := rpc2.HttpPostJson("", l.Blockchain.RpcUrl, "{\"jsonrpc\":\"2.0\",\"method\":\"eth_getBlockByNumber\",\"params\":[\""+fmt.Sprintf("0x%X", vlog.BlockNumber)+"\", true],\"id\":1}")
			//	if err != nil {
			//		log.Errorf("[Handler.SyncBlock] Syncing block by number error: %s\n", errors.WithStack(err))
			//		time.Sleep(3 * time.Second)
			//		continue
			//	}
			//	block := rpc2.ParseJsonBlock(string(blockJson))
			//	//
			//block, err := l.rpc.BlockByNumber(context.Background(), big.NewInt(int64(vlog.BlockNumber)))
			//if err != nil {
			//	l.logger.Errorf("Syncing block by number error: %s\n", errors.WithStack(err))
			//	return nil, errors.WithStack(err)
			//}
			//blockTime = int64(block.Time())
			//blockTimes[int64(vlog.BlockNumber)] = blockTime
		}
		events = append(events, &models.SyncEvent{
			ChainId:         l.conf.ChainId,
			BlockTime:       blockTime,
			BlockNumber:     int64(vlog.BlockNumber),
			BlockHash:       vlog.BlockHash.Hex(),
			BlockLogIndexed: int64(vlog.Index),
			TxIndex:         int64(vlog.TxIndex),
			TxHash:          vlog.TxHash.Hex(),
			EventName:       eventName,
			EventHash:       eventHash.Hex(),
			ContractAddress: strings.ToLower(contractAddress.Hex()),
			Data:            data,
			Status:          models.EventPending,
		})
	}
	return events, nil
}

func (l *EthereumListener) confirm() {
	duration := time.Millisecond * time.Duration(l.conf.BlockInterval) * 5
	for {
		list, err := l.pendingSendMessage(10)
		if err != nil {
			l.logger.Errorf("Get pending send message error: %s\n", err)
			time.Sleep(duration)
			continue
		}
		if len(list) == 0 {
			time.Sleep(duration)
			continue
		}
		var wg sync.WaitGroup
		for _, message := range list {
			wg.Add(1)
			go func(_wg *sync.WaitGroup, message models.Message) {
				defer _wg.Done()
				err = l.confirmMessage(message)
				if err != nil {
					l.logger.Errorf("")
				}
			}(&wg, message)
		}
		wg.Wait()
	}
}

func (l *EthereumListener) confirmMessage(message models.Message) error {
	var callMessage models.Message
	err := l.db.Where("type=? AND from_chain_id=? AND from_id=?", enums.MessageTypeCall, message.FromChainId, message.FromId).First(&callMessage).Error
	if err != nil {
		return err
	}
	err = l.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(models.Message{}).Where("id=?",
			callMessage.Id).
			Update("status", enums.MessageStatusValid).Error
		if err != nil {
			return err
		}
		err = tx.Model(models.Message{}).Where("id=? AND status=?", message.Id, message.Status).
			Update("status", enums.MessageStatusValid).Error
		if err != nil {
			return err
		}
		err = tx.Model(models.Signature{}).Where("chain_id=? AND refer_id=?", message.ChainId, message.FromId).
			Updates(map[string]interface{}{
				"status":       enums.SignatureStatusSuccess,
				"event_id":     message.EventId,
				"block_time":   message.BlockTime,
				"block_number": message.BlockNumber,
				"log_index":    message.LogIndex,
				"tx_hash":      message.TxHash,
			}).Error
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

func (l *EthereumListener) pendingSendMessage(limit int) ([]models.Message, error) {
	var list []models.Message
	err := l.db.Where("`to_chain_id`=? AND `type`=? AND status=?", l.conf.ChainId, enums.MessageTypeSend, enums.MessageStatusPending).Limit(limit).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}
