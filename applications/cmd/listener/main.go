package main

import (
	"bsquared.network/message-sharing-applications/internal/config"
	"bsquared.network/message-sharing-applications/internal/initiates"
	"bsquared.network/message-sharing-applications/internal/serves/listener/bitcoin"
	"bsquared.network/message-sharing-applications/internal/serves/listener/ethereum"
	"bsquared.network/message-sharing-applications/internal/utils/log"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/shopspring/decimal"
)

func main() {
	decimal.DivisionPrecision = 18
	var fileName string
	flag.StringVar(&fileName, "f", "listener", "-f config filename, default: listener")
	flag.Parse()
	cfg := config.LoadConfig(fileName)
	logger := log.NewLogger(fmt.Sprintf("listener-common"), cfg.Log.Level)
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
	bridges, err := config.ParseBridges(cfg.Bridges)
	if err != nil {
		logger.Panicf("parse bridges err: %s", err)
	}

	go func() {
		logger := log.NewLogger(fmt.Sprintf("listener-%s", cfg.Bsquared.Name), cfg.Log.Level)
		rpc, err := initiates.InitEthereumRpc(cfg.Bsquared.RpcUrl)
		if err != nil {
			logger.Panicf("init ethereum rpc err: %s", err)
		}
		ethereum.NewListener(bridges, cfg.Bsquared, rpc, db, logger).Start()
	}()

	go func() {
		logger := log.NewLogger(fmt.Sprintf("listener-%s", cfg.Arbitrum.Name), cfg.Log.Level)
		rpc, err := initiates.InitEthereumRpc(cfg.Arbitrum.RpcUrl)
		if err != nil {
			logger.Panicf("init ethereum rpc err: %s", err)
		}
		ethereum.NewListener(bridges, cfg.Arbitrum, rpc, db, logger).Start()
	}()

	go func() {
		logger := log.NewLogger(fmt.Sprintf("listener-%s", cfg.Bitcoin.Name), cfg.Log.Level)
		rpc, err := initiates.InitBitcoinRpc(cfg.Bitcoin.RpcUrl, cfg.Bitcoin.BtcUser, cfg.Bitcoin.BtcPass, cfg.Bitcoin.DisableTLS)
		if err != nil {
			logger.Panicf("init bitcoin rpc err: %s", err)
		}
		bitcoin.NewListener(bridges, cfg.Bitcoin, cfg.Particle, rpc, db, logger).Start()
	}()
	logger.Info("======================================================")
	select {}
}
