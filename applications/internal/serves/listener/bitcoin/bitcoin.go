package bitcoin

import (
	"bsquared.network/message-sharing-applications/internal/config"
	"bsquared.network/message-sharing-applications/internal/enums"
	"bsquared.network/message-sharing-applications/internal/models"
	"bsquared.network/message-sharing-applications/internal/types"
	"bsquared.network/message-sharing-applications/internal/utils/aa"
	"bsquared.network/message-sharing-applications/internal/utils/ethereum/message"
	"bsquared.network/message-sharing-applications/internal/utils/log"
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"sync"
	"time"
)

var (
	ErrParsePkScript            = errors.New("parse pkscript err")
	ErrDecodeListenAddress      = errors.New("decode listen address err")
	ErrTargetConfirmations      = errors.New("target confirmation number was not reached")
	ErrParsePubKey              = errors.New("parse pubkey failed, not found pubkey or nonsupport ")
	ErrParsePkScriptNullData    = errors.New("parse pkscript null data err")
	ErrParsePkScriptNotNullData = errors.New("parse pkscript not null data err")
)

const (
	// tx type
	TxTypeTransfer = "transfer" // btc transfer
	TxTypeWithdraw = "withdraw" // btc withdraw
)

type BitcoinListener struct {
	conf        config.Blockchain
	particle    config.Particle
	rpc         *rpcclient.Client
	db          *gorm.DB
	logger      *log.Logger
	latestBlock int64
	bridges     map[int64]string
}

func NewListener(bridges map[int64]string, conf config.Blockchain, particle config.Particle, rpc *rpcclient.Client, db *gorm.DB, logger *log.Logger) *BitcoinListener {
	return &BitcoinListener{
		conf:     conf,
		particle: particle,
		rpc:      rpc,
		db:       db,
		logger:   logger,
		bridges:  bridges,
	}
}

func (l *BitcoinListener) Start() {
	if !l.conf.Status {
		l.logger.Infof("status: %t", l.conf.Status)
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go l.syncLatestBlock()
	go l.syncTask()
	go l.handDeposit()
	<-ctx.Done()
}

func (l *BitcoinListener) syncLatestBlock() {
	for {
		duration := time.Millisecond * time.Duration(l.conf.BlockInterval) * 10
		latest, err := l.rpc.GetBlockCount()
		if err != nil {
			l.logger.Errorf("sync latest block failed, err: %v", err)
			time.Sleep(duration)
			continue
		}
		l.latestBlock = int64(latest) - l.conf.SafeBlockNumber
		l.logger.Infof("sync latest block success, latest block: %d", l.latestBlock)
		time.Sleep(duration)
	}
}

func (l *BitcoinListener) syncTask() {
	for {
		duration := time.Millisecond * time.Duration(l.conf.BlockInterval)
		var tasks []models.SyncTask
		err := l.db.Where("`chain_type`=? AND `chain_id`=? AND `status`=?", enums.ChainTypeUTXO, l.conf.ChainId, enums.TaskStatusPending).Limit(20).Find(&tasks).Error
		if err != nil {
			l.logger.Errorf("get task list err: %s", err)
			time.Sleep(duration)
			continue
		}
		if len(tasks) == 0 {
			l.logger.Infof("no task to handle")
			time.Sleep(duration)
			continue
		}
		wg := sync.WaitGroup{}
		for _, task := range tasks {
			wg.Add(1)
			go func(take models.SyncTask, wg *sync.WaitGroup) {
				defer wg.Done()
				err := l.handleTask(task)
				if err != nil {
					l.logger.Errorf("handle task err: %s", err)
				}
			}(task, &wg)
		}
		wg.Wait()
	}
}

func (l *BitcoinListener) handleTask(task models.SyncTask) error {
	l.logger.Infof("start handle task, task id: %d", task.Id)
	var (
		currentBlock   int64 // index current block number
		currentTxIndex int64 // index current block tx index
	)

	currentBlock = task.LatestBlock
	if currentBlock < task.StartBlock {
		currentBlock = task.StartBlock
	}
	currentTxIndex = task.LatestTx

	for {
		if l.latestBlock <= currentBlock {
			//time.Sleep(time.Second * time.Duration(l.config.BlockInterval))
			continue
		}
		if currentTxIndex == 0 {
			currentBlock++
		} else {
			currentTxIndex++
		}

		for i := currentBlock; i <= l.latestBlock; i++ {
			l.logger.Infof("start sync task, task id: %d, current block: %d, current tx index: %d", task.Id, currentBlock, currentTxIndex)
			txResults, blockHeader, err := l.ParseBlock(i, currentTxIndex)
			if err != nil {
				if errors.Is(err, ErrTargetConfirmations) {
					l.logger.Errorf("parse block err: %s", err)
					//time.Sleep(time.Second * time.Duration(l.config.BlockInterval))
				} else {
					l.logger.Errorf("parse block unknown err: %s", err)
				}
				if currentTxIndex == 0 {
					currentBlock = i - 1
				} else {
					currentBlock = i
					currentTxIndex--
				}
				break
			}
			if len(txResults) > 0 {
				currentBlock, currentTxIndex, err = l.HandleResults(txResults, task, blockHeader.Timestamp, i)
				if err != nil {
					l.logger.Errorf("handle results err: %s", err)
					//bis.log.Errorw("failed to handle results", "error", err,
					//	"currentBlock", currentBlock, "currentTxIndex", currentTxIndex, "latestBlock", latestBlock)
					rollback := true
					// not duplicated key, rollback index
					if pgErr, ok := err.(*pgconn.PgError); ok {
						// 23505 duplicate key value violates unique constraint , continue
						if pgErr.Code == "23505" {
							rollback = false
						}
					}

					if rollback {
						if currentTxIndex == 0 {
							currentBlock = i - 1
						} else {
							currentBlock = i
							currentTxIndex--
						}
						break
					}
				}
			}
			currentBlock = i
			currentTxIndex = 0
			task.LatestBlock = currentBlock
			task.LatestTx = currentTxIndex
			if err := l.db.Save(&task).Error; err != nil {
				l.logger.Errorf("save task err: %s", err)
				//bis.log.Errorw("failed to save bitcoin index block", "error", err, "currentBlock", i,
				//	"currentTxIndex", currentTxIndex, "latestBlock", latestBlock)
				// rollback
				currentBlock = i - 1
				break
			}
			//l.logger.Infof("bitcoin indexer parsed currentBlock, i)
			//bis.log.Infow("bitcoin indexer parsed", "currentBlock", i,
			//	"currentTxIndex", currentTxIndex, "latestBlock", latestBlock)
			//time.Sleep(time.Millisecond * time.Duration(l.config.BlockInterval))
		}
	}
	return nil
}

func (l *BitcoinListener) ParseBlock(height int64, txIndex int64) ([]*types.BitcoinTxParseResult, *wire.BlockHeader, error) {
	blockResult, err := l.getBlockByHeight(height)
	if err != nil {
		return nil, nil, err
	}

	blockParsedResult := make([]*types.BitcoinTxParseResult, 0)
	for k, v := range blockResult.Transactions {
		if int64(k) < txIndex {
			continue
		}

		//b.logger.Debugw("parse block", "k", k, "height", height, "txIndex", txIndex, "tx", v.TxHash().String())

		parseTxs, err := l.parseTx(v, k)
		if err != nil {
			return nil, nil, err
		}
		if parseTxs != nil {
			blockParsedResult = append(blockParsedResult, parseTxs)
		}
	}

	return blockParsedResult, &blockResult.Header, nil
}

// getBlockByHeight returns a raw block from the server given its height
func (l *BitcoinListener) getBlockByHeight(height int64) (*wire.MsgBlock, error) {
	blockhash, err := l.rpc.GetBlockHash(height)
	if err != nil {
		return nil, err
	}
	msgBlock, err := l.rpc.GetBlock(blockhash)
	if err != nil {
		return nil, err
	}
	return msgBlock, nil
}
func (l *BitcoinListener) parseTx(txResult *wire.MsgTx, index int) (*types.BitcoinTxParseResult, error) {
	listenAddress := false
	var totalValue int64
	tos := make([]types.BitcoinTo, 0)
	for _, v := range txResult.TxOut {
		pkAddress, err := l.parseAddress(v.PkScript)
		if err != nil {
			if errors.Is(err, ErrParsePkScript) {
				continue
			}
			// parse null data
			if errors.Is(err, ErrParsePkScriptNullData) {
				nullData, err := l.parseNullData(v.PkScript)
				if err != nil {
					continue
				}
				tos = append(tos, types.BitcoinTo{
					Type:     types.BitcoinToTypeNullData,
					NullData: nullData,
				})
			} else {
				return nil, err
			}
		} else {
			parseTo := types.BitcoinTo{
				Address: pkAddress,
				Value:   v.Value,
				Type:    types.BitcoinToTypeNormal,
			}
			tos = append(tos, parseTo)
		}

		var chainParams *chaincfg.Params
		if l.conf.Mainnet {
			chainParams = &chaincfg.MainNetParams
		} else {
			chainParams = &chaincfg.TestNet3Params
		}
		_listenAddress, err := btcutil.DecodeAddress(l.conf.ListenAddress, chainParams)

		// if pk address eq dest listened address, after parse from address by vin prev tx
		if pkAddress == _listenAddress.EncodeAddress() {
			listenAddress = true
			totalValue += v.Value
		}
	}
	if listenAddress {
		fromAddress, err := l.parseFromAddress(txResult)
		if err != nil {
			return nil, fmt.Errorf("vin parse err:%w", err)
		}

		// TODO: temp fix, if from is listened address, continue
		if len(fromAddress) == 0 {
			//b.logger.Warnw("parse from address empty or nonsupport tx type",
			//	"txId", txResult.TxHash().String(),
			//	"listenAddress", b.listenAddress.EncodeAddress())
			return nil, nil
		}

		var chainParams *chaincfg.Params
		if l.conf.Mainnet {
			chainParams = &chaincfg.MainNetParams
		} else {
			chainParams = &chaincfg.TestNet3Params
		}
		_listenAddress, err := btcutil.DecodeAddress(l.conf.ListenAddress, chainParams)
		if err != nil {
			return nil, err
		}
		return &types.BitcoinTxParseResult{
			TxID:   txResult.TxHash().String(),
			TxType: TxTypeTransfer,
			Index:  int64(index),
			Value:  totalValue,
			From:   fromAddress,
			To:     _listenAddress.EncodeAddress(),
			Tos:    tos,
		}, nil
	}
	return nil, nil
}

func (l *BitcoinListener) parseAddress(pkScript []byte) (string, error) {
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

	var chainParams *chaincfg.Params
	if l.conf.Mainnet {
		chainParams = &chaincfg.MainNetParams
	} else {
		chainParams = &chaincfg.TestNet3Params
	}
	pkAddress, err := pk.Address(chainParams)
	if err != nil {
		return "", fmt.Errorf("PKScript to address err:%w", err)
	}
	return pkAddress.EncodeAddress(), nil
}

// parseNullData from pkscript parse null data
func (l *BitcoinListener) parseNullData(pkScript []byte) (string, error) {
	if !txscript.IsNullData(pkScript) {
		return "", ErrParsePkScriptNotNullData
	}
	return hex.EncodeToString(pkScript[1:]), nil
}

func (l *BitcoinListener) parseFromAddress(txResult *wire.MsgTx) (fromAddress []types.BitcoinFrom, err error) {
	for _, vin := range txResult.TxIn {
		// get prev tx hash
		prevTxID := vin.PreviousOutPoint.Hash
		vinResult, err := l.rpc.GetRawTransaction(&prevTxID)
		if err != nil {
			return nil, fmt.Errorf("vin get raw transaction err:%w", err)
		}
		if len(vinResult.MsgTx().TxOut) == 0 {
			return nil, fmt.Errorf("vin txOut is null")
		}
		vinPKScript := vinResult.MsgTx().TxOut[vin.PreviousOutPoint.Index].PkScript
		//  script to address
		vinPkAddress, err := l.parseAddress(vinPKScript)
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

func (l *BitcoinListener) HandleResults(
	txResults []*types.BitcoinTxParseResult,
	syncTask models.SyncTask,
	btcBlockTime time.Time,
	currentBlock int64,
) (int64, int64, error) {
	for _, v := range txResults {
		// if from is listen address, skip
		if l.ToInFroms(v.From, v.To) {
			//bis.log.Infow("current transaction from is listen address", "currentBlock", currentBlock, "currentTxIndex", v.Index, "data", v)
			continue
		}

		syncTask.LatestBlock = currentBlock
		syncTask.LatestTx = v.Index
		// write db
		err := l.SaveParsedResult(
			v,
			currentBlock,
			models.DepositB2TxStatusPending,
			btcBlockTime,
			syncTask,
		)
		if err != nil {
			//bis.log.Errorw("failed to save bitcoin index tx", "error", err,
			//	"data", v)
			return currentBlock, v.Index, err
		}
		//bis.log.Infow("save bitcoin index tx success", "currentBlock", currentBlock, "currentTxIndex", v.Index, "data", v)
		time.Sleep(time.Second * 2)
	}
	return currentBlock, 0, nil
}

func (l *BitcoinListener) ToInFroms(a []types.BitcoinFrom, s string) bool {
	for _, i := range a {
		if i.Address == s {
			return true
		}
	}
	return false
}

func (l *BitcoinListener) SaveParsedResult(
	parseResult *types.BitcoinTxParseResult,
	btcBlockNumber int64,
	b2TxStatus int,
	btcBlockTime time.Time,
	syncTask models.SyncTask,
) error {
	// write db
	err := l.db.Transaction(func(tx *gorm.DB) error {
		if len(parseResult.From) == 0 {
			return fmt.Errorf("parse result from empty")
		}

		if len(parseResult.To) == 0 {
			return fmt.Errorf("parse result to empty")
		}

		if len(parseResult.Tos) == 0 {
			return fmt.Errorf("parse result to empty")
		}

		//bis.log.Infow("parseResult:", "result", parseResult)
		existsEvmAddressData := false // The evm address is processed only if it exists. Otherwise, aa is used
		parsedEvmAddress := ""        // evm address
		for _, v := range parseResult.Tos {
			// only handle first null data
			if existsEvmAddressData {
				continue
			}
			if v.Type == types.BitcoinToTypeNullData {
				decodeNullData, err := hex.DecodeString(v.NullData)
				if err != nil {
					//bis.log.Errorw("decode null data err", "error", err, "nullData", v.NullData)
					continue
				}
				evmAddress := bytes.TrimSpace(decodeNullData[1:])
				if common.IsHexAddress(string(evmAddress)) {
					existsEvmAddressData = true
					parsedEvmAddress = string(evmAddress)
					for k := range parseResult.From {
						parseResult.From[k].Type = types.BitcoinFromTypeEvm
						parseResult.From[k].EvmAddress = parsedEvmAddress
					}
				}
			}
		}
		froms, err := json.Marshal(parseResult.From)
		if err != nil {
			return err
		}
		tos, err := json.Marshal(parseResult.Tos)
		if err != nil {
			return err
		}
		// if existed, update deposit record
		var deposit models.Deposit
		err = tx.
			Set("gorm:query_option", "FOR UPDATE").
			First(&deposit,
				fmt.Sprintf("%s = ?", models.Deposit{}.Column().BtcTxHash),
				parseResult.TxID).Error
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			deposit := models.Deposit{
				BtcBlockNumber: btcBlockNumber,
				BtcTxIndex:     parseResult.Index,
				BtcTxHash:      parseResult.TxID,
				BtcFrom:        parseResult.From[0].Address,
				BtcTos:         string(tos),
				BtcTo:          parseResult.To,
				BtcValue:       parseResult.Value,
				BtcFroms:       string(froms),
				B2TxStatus:     b2TxStatus,
				BtcBlockTime:   btcBlockTime,
				B2TxRetry:      0,
				ListenerStatus: models.ListenerStatusSuccess,
				CallbackStatus: models.CallbackStatusPending,
			}
			if existsEvmAddressData {
				deposit.BtcFromEvmAddress = parsedEvmAddress
			}
			err = tx.Create(&deposit).Error
			if err != nil {
				//bis.log.Errorw("failed to save tx parsed result", "error", err)
				return err
			}
		} else if deposit.CallbackStatus == models.CallbackStatusSuccess &&
			deposit.ListenerStatus == models.ListenerStatusPending {
			if deposit.BtcValue != parseResult.Value || deposit.BtcFrom != parseResult.From[0].Address {
				return fmt.Errorf("invalid parameter")
			}
			// if existed, update deposit record
			updateFields := map[string]interface{}{
				models.Deposit{}.Column().BtcBlockNumber: btcBlockNumber,
				models.Deposit{}.Column().BtcTxIndex:     parseResult.Index,
				models.Deposit{}.Column().BtcFroms:       string(froms),
				models.Deposit{}.Column().BtcTos:         string(tos),
				models.Deposit{}.Column().BtcBlockTime:   btcBlockTime,
				models.Deposit{}.Column().ListenerStatus: models.ListenerStatusSuccess,
			}
			if existsEvmAddressData {
				updateFields[models.Deposit{}.Column().BtcFromEvmAddress] = parsedEvmAddress
			}
			err = tx.Model(&models.Deposit{}).Where("id = ?", deposit.Id).Updates(updateFields).Error
			if err != nil {
				//bis.log.Errorw("failed to update tx parsed result", "error", err)
				return err
			}
		}

		if err := tx.Save(&syncTask).Error; err != nil {
			//bis.log.Errorw("failed to save bitcoin tx index", "error", err)
			return err
		}
		return nil
	})
	return err
}

func (l *BitcoinListener) handDeposit() {
	duration := time.Millisecond * time.Duration(l.conf.BlockInterval) * 10
	for {
		var list []models.Deposit
		err := l.db.Where("status=?", enums.DepositStatusPending).Find(&list).Error
		if err != nil {
			l.logger.Errorf("[Handler.handDeposit] err: %s", err)
			time.Sleep(duration)
			continue
		}
		if len(list) == 0 {
			l.logger.Info("[Handler.handDeposit] deposit list is empty")
			time.Sleep(duration)
			continue
		}

		wg := sync.WaitGroup{}
		for _, one := range list {
			wg.Add(1)
			go func(one models.Deposit, wg *sync.WaitGroup) {
				defer wg.Done()
				err := l.handleMessage(one)
				if err != nil {
					l.logger.Errorf("[Handler.handDeposit] Deposit ID: %d , err: %s \n", one.Id, err)
				}
			}(one, &wg)
		}
		wg.Wait()
	}
}

func (l *BitcoinListener) GetDepositAddress(deposit models.Deposit) (string, error) {
	if deposit.BtcFromEvmAddress != "" && common.IsHexAddress(deposit.BtcFromEvmAddress) {
		return deposit.BtcFromEvmAddress, nil
	} else {
		evmAddress, err := aa.BitcoinAddressToEthAddress(l.particle.AAPubKeyAPI, deposit.BtcFrom,
			l.particle.Url, l.particle.ChainId, l.particle.ProjectUuid, l.particle.ProjectKey)
		if err != nil {
			l.logger.Errorf("[Handler.GetDepositAddress] Deposit ID: %d , err: %s \n", deposit.Id, err)
			return "", err
		}
		return evmAddress, nil
	}
}

func (l *BitcoinListener) handleMessage(deposit models.Deposit) error {
	toChainId := l.conf.ToChainId
	toContractAddress := l.conf.ToContractAddress

	depositAddress, err := l.GetDepositAddress(deposit)
	if err != nil && err.Error() != "AAGetBTCAccount not found" {
		l.logger.Errorf("[Handler.handleMessage] GetDepositAddress err: %s \n", err)
		return err
	} else if err != nil && err.Error() == "AAGetBTCAccount not found" {
		err = l.db.Model(models.Deposit{}).
			Where("id=?", deposit.Id).
			Update("status", enums.DepositStatusInvalid).Error
		return nil
	}

	data := message.EncodeSendData(deposit.BtcTxHash, deposit.BtcFrom, depositAddress, decimal.New(deposit.BtcValue, 0))
	var ToMessageBridge string
	messageBridge, ok := l.bridges[toChainId]
	if ok {
		ToMessageBridge = messageBridge
	} else {
		err = l.db.Model(models.Deposit{}).
			Where("id=?", deposit.Id).
			Update("status", enums.DepositStatusInvalid).Error
		return nil
	}
	msg := models.Message{
		ChainId:           l.conf.ChainId,
		Type:              enums.MessageTypeCall,
		FromChainId:       l.conf.ChainId,
		FromSender:        common.HexToAddress("0x0").Hex(),
		FromMessageBridge: deposit.BtcTo,
		FromId:            common.HexToHash(deposit.BtcTxHash).Hex(),
		ToChainId:         toChainId,
		ToMessageBridge:   ToMessageBridge,
		ToContractAddress: toContractAddress,
		ToBytes:           hexutil.Encode(data),
		Signatures:        "{}",
		Status:            enums.MessageStatusValidating,
		Blockchain: models.Blockchain{
			EventId:     deposit.Id,
			BlockTime:   deposit.BtcBlockTime.Unix(),
			BlockNumber: deposit.BtcBlockNumber,
			LogIndex:    deposit.BtcTxIndex,
			TxHash:      common.HexToHash(deposit.BtcTxHash).Hex(),
		},
	}
	err = l.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&msg).Error
		if err != nil {
			l.logger.Errorf("create message failed: %s", err.Error())
			return err
		}

		if deposit.BtcFromEvmAddress != "" {
			err = tx.Model(models.Deposit{}).
				Where("id=?", deposit.Id).
				Update("status", enums.DepositStatusValid).Error
		} else {
			err = tx.Model(models.Deposit{}).
				Where("id=?", deposit.Id).
				Update("btc_from_aa_address", depositAddress).
				Update("status", enums.DepositStatusValid).Error
		}
		if err != nil {
			l.logger.Errorf("update deposit failed: %s", err.Error())
			return err
		}
		return nil
	})
	if err != nil {
		l.logger.Errorf("error: %s\n", err)
		return err
	}
	return nil
}
