package main

import (
	"bsquared.network/message-sharing-applications/internal/config"
	"bsquared.network/message-sharing-applications/internal/initiates"
	"bsquared.network/message-sharing-applications/internal/serves/proposer"
	"bsquared.network/message-sharing-applications/internal/utils/log"
	"bsquared.network/message-sharing-applications/internal/vo"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/shopspring/decimal"
)

func main() {
	decimal.DivisionPrecision = 18
	var fileName string
	flag.StringVar(&fileName, "f", "proposer", "-f config filename, default: proposer")
	flag.Parse()
	cfg := config.LoadConfig(fileName)
	logger := log.NewLogger(fmt.Sprintf("proposer-common"), cfg.Log.Level)
	value, err := json.Marshal(cfg)
	if err != nil {
		logger.Panicf("json marshal err: %s", err)
	}
	logger.Infof("config: %s", value)
	logger.Info("------------------------------------------------------")

	db, err := initiates.InitDB(cfg.Database)
	if err != nil {
		logger.Panicf("init db err: %s", err)
	}
	go func() {
		logger := log.NewLogger(fmt.Sprintf("proposer-%s", cfg.Bitcoin.Name), cfg.Log.Level)
		pk, host, err := initiates.InitListenHost(cfg.Bitcoin.NodePort, cfg.Bitcoin.NodeKey)
		if err != nil {
			logger.Panicf("init host err: %s", err)
		}
		rpc, err := initiates.InitBitcoinRpc(cfg.Bitcoin.RpcUrl, cfg.Bitcoin.BtcUser, cfg.Bitcoin.BtcPass, cfg.Bitcoin.DisableTLS)
		if err != nil {
			logger.Panicf("init ethereum rpc err: %s", err)
		}
		proposer.NewProposer(pk, host, db, &vo.RpcClient{BtcRpc: rpc}, logger, cfg.Bitcoin, cfg.Particle).Start()
	}()
	go func() {
		logger := log.NewLogger(fmt.Sprintf("proposer-%s", cfg.Bsquared.Name), cfg.Log.Level)

		pk, host, err := initiates.InitListenHost(cfg.Bsquared.NodePort, cfg.Bsquared.NodeKey)
		if err != nil {
			logger.Panicf("init host err: %s", err)
		}
		rpc, err := initiates.InitEthereumRpc(cfg.Bsquared.RpcUrl)
		if err != nil {
			logger.Panicf("init ethereum rpc err: %s", err)
		}
		proposer.NewProposer(pk, host, db, &vo.RpcClient{EthRpc: rpc}, logger, cfg.Bsquared, cfg.Particle).Start()
	}()
	go func() {
		logger := log.NewLogger(fmt.Sprintf("proposer-%s", cfg.Arbitrum.Name), cfg.Log.Level)
		pk, host, err := initiates.InitListenHost(cfg.Arbitrum.NodePort, cfg.Arbitrum.NodeKey)
		if err != nil {
			logger.Panicf("init host err: %s", err)
		}
		rpc, err := initiates.InitEthereumRpc(cfg.Arbitrum.RpcUrl)
		if err != nil {
			logger.Panicf("init ethereum rpc err: %s", err)
		}
		proposer.NewProposer(pk, host, db, &vo.RpcClient{EthRpc: rpc}, logger, cfg.Arbitrum, cfg.Particle).Start()
	}()
	logger.Info("======================================================")
	select {}
}
