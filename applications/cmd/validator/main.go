package main

import (
	"bsquared.network/message-sharing-applications/internal/config"
	"bsquared.network/message-sharing-applications/internal/initiates"
	"bsquared.network/message-sharing-applications/internal/serves/validator"
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
	flag.StringVar(&fileName, "f", "validator", "-f config filename, default: validator")
	flag.Parse()
	cfg := config.LoadConfig(fileName)
	logger := log.NewLogger(fmt.Sprintf("validator-common"), cfg.Log.Level)
	value, err := json.Marshal(cfg)
	if err != nil {
		logger.Panicf("json marshal err: %s", err)
	}
	logger.Infof("config: %s", value)
	logger.Info("------------------------------------------------------")
	go func() {
		logger := log.NewLogger(fmt.Sprintf("validator-%s", cfg.Bitcoin.Name), uint32(cfg.Log.Level))
		pk, host, err := initiates.InitHost(cfg.Bitcoin.NodeKey)
		if err != nil {
			logger.Panicf("init host err: %s", err)
		}
		rpc, err := initiates.InitBitcoinRpc(cfg.Bitcoin.RpcUrl, cfg.Bitcoin.BtcUser, cfg.Bitcoin.BtcPass, cfg.Bitcoin.DisableTLS)
		if err != nil {
			logger.Panicf("init ethereum rpc err: %s", err)
		}
		validator.NewValidator(pk, host, logger, &vo.RpcClient{BtcRpc: rpc}, cfg.Particle, cfg.Bitcoin).Start()
	}()
	go func() {
		logger := log.NewLogger(fmt.Sprintf("validator-%s", cfg.Bsquared.Name), uint32(cfg.Log.Level))
		pk, host, err := initiates.InitHost(cfg.Bsquared.NodeKey)
		if err != nil {
			logger.Panicf("init host err: %s", err)
		}
		rpc, err := initiates.InitEthereumRpc(cfg.Bsquared.RpcUrl)
		if err != nil {
			logger.Panicf("init ethereum rpc err: %s", err)
		}
		validator.NewValidator(pk, host, logger, &vo.RpcClient{EthRpc: rpc}, config.Particle{}, cfg.Bsquared).Start()
	}()
	go func() {
		logger := log.NewLogger(fmt.Sprintf("validator-%s", cfg.Arbitrum.Name), uint32(cfg.Log.Level))
		pk, host, err := initiates.InitHost(cfg.Arbitrum.NodeKey)
		if err != nil {
			logger.Panicf("init host err: %s", err)
		}
		rpc, err := initiates.InitEthereumRpc(cfg.Arbitrum.RpcUrl)
		if err != nil {
			logger.Panicf("init ethereum rpc err: %s", err)
		}
		validator.NewValidator(pk, host, logger, &vo.RpcClient{EthRpc: rpc}, config.Particle{}, cfg.Arbitrum).Start()
	}()
	logger.Info("======================================================")
	select {}
}
