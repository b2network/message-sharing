# Dapp Deployment

<!-- TOC -->
* [Dapp Deployment](#dapp-deployment)
  * [Message-Sharing Mechanism Overview](#message-sharing-mechanism-overview)
  * [Message-Sharing Services](#message-sharing-services)
    * [Listener](#listener)
    * [Proposer](#proposer)
    * [Validator](#validator)
    * [Builder](#builder)
  * [Deployment](#deployment)
    * [Build](#build)
    * [Database](#database)
    * [Config](#config)
      * [Yaml config](#yaml-config)
      * [Env config](#env-config)
    * [Quick start](#quick-start)
<!-- TOC -->

## Message-Sharing Mechanism Overview

The system is composed of four main services: Listener, Proposer, Validator, and Builder. Each service has a specific
role in ensuring secure and efficient cross-chain message transmission.

## Message-Sharing Services

### Listener

Role:

> The Listener service is responsible for monitoring the message-sharing smart contract for specific events.
>

Functionality:

> It listens for Call and Send events emitted by the contract. \
> Serves as a data source for the Proposer and Builder services by capturing and relaying event data.
>

### Proposer

Role:
> The Proposer service is tasked with transmitting transaction information to Validators and gathering valid signatures.
>
Functionality:
> Utilizes a peer-to-peer (p2p) protocol to communicate transaction details that require validation. \
> Collects signatures from Validators to confirm the legitimacy of the transaction data.
>

### Validator

Role:
> The Validator service is responsible for validating transaction information received from the Proposer.
>
Functionality:
> Performs on-chain data verification to ensure the authenticity of the transaction data. \
> Provides signatures for legitimate transaction data back to the Proposer.
>

### Builder

Role:
> The Builder service processes messages and their corresponding signatures to complete transaction construction and
> broadcasting.
>
Functionality:
> Constructs transactions based on validated messages and signatures. \
> Broadcasts these transactions onto the blockchain.
>

## Deployment

### Build

```
// Navigate to the Working Directory:
$ cd message-sharing/applications/
// Build Listener
$ go build -o listener cmd/listener/main.go
// Build Proposer
$ go build -o proposer cmd/proposer/main.go
// Build Validator
$ go build -o validator cmd/validator/main.go
// Build Builder
$ go build -o builder cmd/builder/main.go
```

### Database

The Message Sharing service depends on a MySQL service, and a service instance needs to be created first.

1. Create database b2_message

```
CREATE DATABASE `b2_message` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci */ /*!80016 DEFAULT ENCRYPTION='N' */
```

2. Create tables

2.1 messages

```
CREATE TABLE `messages` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `type` int NOT NULL COMMENT 'type',
  `chain_id` bigint NOT NULL COMMENT 'chain id',
  `from_chain_id` bigint NOT NULL COMMENT 'from_chain_id',
  `from_sender` varchar(66) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT 'from_sender',
  `from_message_bridge` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT 'from_message_bridge',
  `from_id` varchar(128) NOT NULL COMMENT 'from_id',
  `to_chain_id` bigint NOT NULL COMMENT 'to_chain_id',
  `to_message_bridge` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT 'to_message_bridge',
  `to_contract_address` varchar(66) NOT NULL COMMENT 'to_contract_address',
  `to_bytes` text NOT NULL COMMENT 'to_bytes',
  `event_id` bigint NOT NULL COMMENT 'event_id',
  `block_time` bigint NOT NULL COMMENT 'block_time',
  `block_number` bigint NOT NULL COMMENT 'block_number',
  `log_index` bigint NOT NULL COMMENT 'log_index',
  `tx_hash` varchar(128) NOT NULL COMMENT 'tx_hash',
  `status` tinyint NOT NULL DEFAULT '0' COMMENT 'status',
  `signatures` json NOT NULL COMMENT 'signatures',
  `signatures_count` int NOT NULL DEFAULT '0' COMMENT 'signatures count',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1000000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
```

2.2 message_signatures

```
CREATE TABLE `message_signatures` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `message_id` bigint NOT NULL COMMENT 'message id',
  `signer` varchar(66) NOT NULL COMMENT 'from_contract_address',
  `signature` varchar(256) NOT NULL COMMENT 'signatures',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1000000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
```

2.3 signatures

```
CREATE TABLE `signatures` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `chain_id` bigint NOT NULL COMMENT 'chain id',
  `refer_id` varchar(128) NOT NULL COMMENT ' refer id',
  `nonce` bigint NOT NULL COMMENT ' nonce',
  `type` bigint NOT NULL COMMENT ' type',
  `data` text NOT NULL COMMENT ' data',
  `value` decimal(64,18) NOT NULL DEFAULT '0.000000000000000000' COMMENT ' value',
  `address` varchar(42) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT ' sender address',
  `status` varchar(32) NOT NULL COMMENT ' status',
  `event_id` bigint NOT NULL COMMENT ' event_id',
  `block_time` bigint NOT NULL COMMENT ' block_time',
  `block_number` bigint NOT NULL COMMENT ' block_number',
  `log_index` bigint NOT NULL COMMENT ' log_index',
  `tx_hash` varchar(66) NOT NULL COMMENT ' Tx Hash',
  `signature` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT ' signature',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1000037 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
```

2.4 sync_events

```
CREATE TABLE `sync_events` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `sync_block_id` bigint NOT NULL COMMENT ' sync block id',
  `block_time` bigint NOT NULL COMMENT ' block time',
  `block_number` bigint NOT NULL COMMENT ' block height',
  `block_hash` varchar(66) NOT NULL COMMENT ' block hash',
  `block_log_indexed` bigint NOT NULL COMMENT ' block log index',
  `tx_index` bigint NOT NULL COMMENT ' tx index',
  `tx_hash` varchar(66) NOT NULL COMMENT ' tx hash',
  `event_name` varchar(32) NOT NULL COMMENT ' event name',
  `event_hash` varchar(66) NOT NULL COMMENT ' event hash',
  `contract_address` varchar(42) NOT NULL COMMENT ' contract address',
  `data` json NOT NULL COMMENT ' data content',
  `status` varchar(32) NOT NULL COMMENT ' status',
  `retry_count` bigint NOT NULL DEFAULT '0' COMMENT 'retry_count',
  `chain_id` bigint NOT NULL COMMENT 'chain id',
  PRIMARY KEY (`id`),
  KEY `status_index` (`status`),
  KEY `idx_event_hash` (`sync_block_id`,`block_log_indexed`,`tx_hash`,`event_hash`)
) ENGINE=InnoDB AUTO_INCREMENT=1000000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
```

2.5 sync_events_history

```
CREATE TABLE `sync_events_history` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `chain_id` bigint NOT NULL COMMENT 'chain id',
  `sync_block_id` bigint NOT NULL COMMENT ' sync block id',
  `block_time` bigint NOT NULL COMMENT ' block time',
  `block_number` bigint NOT NULL COMMENT ' block height',
  `block_hash` varchar(66) NOT NULL COMMENT ' block hash',
  `block_log_indexed` bigint NOT NULL COMMENT ' block log index',
  `tx_index` bigint NOT NULL COMMENT ' tx index',
  `tx_hash` varchar(66) NOT NULL COMMENT ' tx hash',
  `event_name` varchar(32) NOT NULL COMMENT ' event name',
  `event_hash` varchar(66) NOT NULL COMMENT ' event hash',
  `contract_address` varchar(42) NOT NULL COMMENT ' contract address',
  `data` json NOT NULL COMMENT ' data content',
  `status` varchar(32) NOT NULL COMMENT ' status',
  `retry_count` int NOT NULL COMMENT 'retry_count',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1000000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
```

2.6 sync_tasks

```
CREATE TABLE `sync_tasks` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `chain_type` tinyint NOT NULL DEFAULT '0' COMMENT 'chain type',
  `chain_id` bigint NOT NULL COMMENT 'chain id',
  `latest_block` bigint NOT NULL COMMENT ' handle block height',
  `latest_tx` bigint NOT NULL DEFAULT '0' COMMENT 'latest_tx',
  `start_block` bigint NOT NULL COMMENT ' start block',
  `end_block` bigint NOT NULL COMMENT ' end block',
  `handle_num` bigint NOT NULL COMMENT ' handle num',
  `contracts` text NOT NULL COMMENT ' contracts address, multiple use, split',
  `status` tinyint NOT NULL DEFAULT '0' COMMENT ' status',
  PRIMARY KEY (`id`),
  KEY `status_index` (`status`)
) ENGINE=InnoDB AUTO_INCREMENT=1000000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
```

2.7 deposit_history

```
CREATE TABLE `deposit_history` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `btc_block_number` bigint NOT NULL COMMENT 'btc_block_number',
  `btc_tx_index` bigint NOT NULL COMMENT 'btc_tx_index',
  `btc_tx_hash` varchar(128) NOT NULL COMMENT ' btc_tx_hash',
  `btc_tx_type` tinyint NOT NULL DEFAULT '0' COMMENT 'chain type',
  `btc_froms` json DEFAULT NULL COMMENT ' btc_froms',
  `btc_from` varchar(128) NOT NULL COMMENT ' btc_from',
  `btc_tos` json DEFAULT NULL COMMENT ' btc_tos',
  `btc_to` varchar(128) NOT NULL COMMENT ' btc_to',
  `btc_from_aa_address` varchar(128) NOT NULL COMMENT ' btc_from_aa_address',
  `btc_from_evm_address` varchar(128) NOT NULL COMMENT ' btc_from_evm_address',
  `btc_value` bigint NOT NULL COMMENT 'btc_tx_index',
  `b2_tx_from` varchar(128) NOT NULL COMMENT ' b2_tx_from',
  `b2_tx_hash` varchar(128) NOT NULL COMMENT ' b2_tx_hash',
  `b2_tx_nonce` bigint NOT NULL COMMENT 'b2_tx_nonce',
  `b2_tx_status` bigint NOT NULL COMMENT 'b2_tx_status',
  `b2_tx_retry` bigint NOT NULL COMMENT 'b2_tx_retry',
  `b2_eoa_tx_from` varchar(128) NOT NULL COMMENT ' b2_eoa_tx_from',
  `b2_eoa_tx_nonce` bigint NOT NULL COMMENT 'b2_eoa_tx_nonce',
  `b2_eoa_tx_hash` varchar(128) NOT NULL COMMENT ' b2_eoa_tx_hash',
  `b2_eoa_tx_status` bigint NOT NULL COMMENT 'b2_eoa_tx_status',
  `btc_block_time` datetime NOT NULL COMMENT 'btc_block_time',
  `callback_status` bigint NOT NULL COMMENT ' callback_status',
  `listener_status` bigint NOT NULL COMMENT ' listener_status',
  `b2_tx_check` bigint NOT NULL COMMENT ' b2_tx_check',
  `status` tinyint NOT NULL DEFAULT '0' COMMENT 'status',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1000000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
```

### Config

#### Yaml config

Before starting the service, you need to handle the service configuration files. Please refer to configuration
files [builder.yaml](../../applications/config/listener.yaml)、[proposer.yaml](../../applications/config/proposer.yaml)
、[validator.yaml](../../applications/config/validator.yaml) and [builder.yaml](../../applications/config/builder.yaml)
for specific details.

#### Env config

listener.env

```
APP_LOG_LEVEL=6

APP_PARTICLE_URL=https://rpc.particle.network/evm-chain
APP_PARTICLE_CHAINID=1123
APP_PARTICLE_PROJECTUUID=0000000000000000000000000000000000000000
APP_PARTICLE_PROJECTKEY=0000000000000000000000000000000000000000
APP_PARTICLE_AAPUBKEYAPI=https://bridge-aa-dev.bsquared.network

APP_DATABASE_USERNAME=root
APP_DATABASE_PASSWORD=123456
APP_DATABASE_HOST=127.0.0.1
APP_DATABASE_PORT=3306
APP_DATABASE_DBNAME=b2_message
APP_DATABASE_LOGLEVEL=4

APP_BITCOIN_NAME=bitcoin
APP_BITCOIN_STATUS=true
APP_BITCOIN_CHAINTYPE=2
APP_BITCOIN_CHAINID=0
APP_BITCOIN_MAINNET=false
APP_BITCOIN_RPCURL=127.0.0.1:8085
APP_BITCOIN_SAFEBLOCKNUMBER=3
APP_BITCOIN_LISTENADDRESS=muGFcyjuyURJJsXaLXHCm43jLBmGPPU7ME
APP_BITCOIN_BLOCKINTERVAL=6000
APP_BITCOIN_TOCHAINID=1123
APP_BITCOIN_TOCONTRACTADDRESS=0x0000000000000000000000000000000000000000
APP_BITCOIN_BTCUSER=test
APP_BITCOIN_BTCPASS=test
APP_BITCOIN_DISABLETLS=false

APP_BSQUARED_NAME=bsquared
APP_BSQUARED_STATUS=true
APP_BSQUARED_CHAINTYPE=1
APP_BSQUARED_MAINNET=false
APP_BSQUARED_CHAINID=1123
APP_BSQUARED_RPCURL=127.0.0.1:8084
APP_BSQUARED_SAFEBLOCKNUMBER=1
APP_BSQUARED_LISTENADDRESS=0x0000000000000000000000000000000000000000
APP_BSQUARED_BLOCKINTERVAL=2000
APP_BSQUARED_BUILDERS=0x0000000000000000000000000000000000000000000000000000000000000000

APP_ARBITRUM_NAME=arbitrum
APP_ARBITRUM_STATUS=true
APP_ARBITRUM_CHAINTYPE=1
APP_ARBITRUM_CHAINID=421614
APP_ARBITRUM_MAINNET=false
APP_ARBITRUM_RPCURL=127.0.0.1:8083
APP_ARBITRUM_SAFEBLOCKNUMBER=1
APP_ARBITRUM_LISTENADDRESS=0x0000000000000000000000000000000000000000
APP_ARBITRUM_BLOCKINTERVAL=100
APP_ARBITRUM_BUILDERS=0x0000000000000000000000000000000000000000000000000000000000000000
```

proposer.env

```
APP_LOG_LEVEL=6

APP_DATABASE_USERNAME=root
APP_DATABASE_PASSWORD=123456
APP_DATABASE_HOST=127.0.0.1
APP_DATABASE_PORT=3306
APP_DATABASE_DBNAME=b2_message
APP_DATABASE_LOGLEVEL=4

APP_BSQUARED_NAME=bsquared
APP_BSQUARED_STATUS=true
APP_BSQUARED_CHAINTYPE=1
APP_BSQUARED_MAINNET=false
APP_BSQUARED_CHAINID=1123
APP_BSQUARED_RPCURL=127.0.0.1:8081
APP_BSQUARED_SAFEBLOCKNUMBER=1
APP_BSQUARED_LISTENADDRESS=0x0000000000000000000000000000000000000000
APP_BSQUARED_BLOCKINTERVAL=2000
APP_BSQUARED_NODEKEY=0000000000000000000000000000000000000000000000000000000000000000
APP_BSQUARED_NODEPORT=20000
APP_BSQUARED_SIGNATUREWEIGHT=1
APP_BSQUARED_VALIDATORS=0x0000000000000000000000000000000000000000

APP_ARBITRUM_NAME=arbitrum
APP_ARBITRUM_STATUS=true
APP_ARBITRUM_CHAINTYPE=1
APP_ARBITRUM_CHAINID=421614
APP_ARBITRUM_MAINNET=false
APP_ARBITRUM_RPCURL=127.0.0.1:8082
APP_ARBITRUM_SAFEBLOCKNUMBER=1
APP_ARBITRUM_LISTENADDRESS=0x0000000000000000000000000000000000000000
APP_ARBITRUM_BLOCKINTERVAL=100
APP_ARBITRUM_NODEKEY=0000000000000000000000000000000000000000000000000000000000000000
APP_ARBITRUM_NODEPORT=20001
APP_ARBITRUM_SIGNATUREWEIGHT=1
APP_ARBITRUM_VALIDATORS=0x0000000000000000000000000000000000000000

APP_BITCOIN_NAME=bitcoin
APP_BITCOIN_STATUS=true
APP_BITCOIN_CHAINTYPE=2
APP_BITCOIN_CHAINID=0
APP_BITCOIN_MAINNET=false
APP_BITCOIN_RPCURL=127.0.0.1:8083
APP_BITCOIN_SAFEBLOCKNUMBER=3
APP_BITCOIN_LISTENADDRESS=muGFcyjuyURJJsXaLXHCm43jLBmGPPU7ME
APP_BITCOIN_BLOCKINTERVAL=6000
APP_BITCOIN_TOCHAINID=1123
APP_BITCOIN_TOCONTRACTADDRESS=0x0000000000000000000000000000000000000000
APP_BITCOIN_BTCUSER=000000000000000000
APP_BITCOIN_BTCPASS=000000000000000000
APP_BITCOIN_DISABLETLS=true
APP_BITCOIN_NODEKEY=0000000000000000000000000000000000000000000000000000000000000000
APP_BITCOIN_NODEPORT=20002
APP_BITCOIN_SIGNATUREWEIGHT=1
APP_BITCOIN_VALIDATORS=0x0000000000000000000000000000000000000000

APP_PARTICLE_URL=https://rpc.particle.network/evm-chain
APP_PARTICLE_CHAINID=1123
APP_PARTICLE_PROJECTUUID=000000000000000000
APP_PARTICLE_PROJECTKEY=000000000000000000
APP_PARTICLE_AAPUBKEYAPI=https://bridge-aa-dev.bsquared.network
```

validator.env

```
APP_LOG_LEVEL=6

APP_BSQUARED_NAME=bsquared
APP_BSQUARED_STATUS=true
APP_BSQUARED_CHAINTYPE=1
APP_BSQUARED_MAINNET=false
APP_BSQUARED_CHAINID=1123
APP_BSQUARED_RPCURL=127.0.0.1:8081
APP_BSQUARED_SAFEBLOCKNUMBER=1
APP_BSQUARED_LISTENADDRESS=0x0000000000000000000000000000000000000000
APP_BSQUARED_BLOCKINTERVAL=2000
APP_BSQUARED_NODEKEY=0000000000000000000000000000000000000000000000000000000000000000
APP_BSQUARED_ENDPOINT=/ip4/127.0.0.1/tcp/20000/p2p/16Uiu2HAkwynt59WSsNRS9sk1aszgeQ1PXUS8ax3a3tsewaVMgvZX
APP_BSQUARED_SIGNATUREWEIGHT=1

APP_ARBITRUM_NAME=arbitrum
APP_ARBITRUM_STATUS=true
APP_ARBITRUM_CHAINTYPE=1
APP_ARBITRUM_CHAINID=421614
APP_ARBITRUM_MAINNET=false
APP_ARBITRUM_RPCURL=127.0.0.1:8082
APP_ARBITRUM_SAFEBLOCKNUMBER=1
APP_ARBITRUM_LISTENADDRESS=0x0000000000000000000000000000000000000000
APP_ARBITRUM_BLOCKINTERVAL=100
APP_ARBITRUM_NODEKEY=0000000000000000000000000000000000000000000000000000000000000000
APP_ARBITRUM_ENDPOINT=/ip4/127.0.0.1/tcp/20001/p2p/16Uiu2HAkwynt59WSsNRS9sk1aszgeQ1PXUS8ax3a3tsewaVMgvZX
APP_ARBITRUM_SIGNATUREWEIGHT=1

APP_BITCOIN_NAME=bitcoin
APP_BITCOIN_STATUS=true
APP_BITCOIN_CHAINTYPE=2
APP_BITCOIN_CHAINID=0
APP_BITCOIN_MAINNET=false
APP_BITCOIN_RPCURL=127.0.0.1:8083
APP_BITCOIN_SAFEBLOCKNUMBER=3
APP_BITCOIN_LISTENADDRESS=muGFcyjuyURJJsXaLXHCm43jLBmGPPU7ME
APP_BITCOIN_BLOCKINTERVAL=6000
APP_BITCOIN_TOCHAINID=1123
APP_BITCOIN_TOCONTRACTADDRESS=0x0000000000000000000000000000000000000000
APP_BITCOIN_BTCUSER=000000000000000000
APP_BITCOIN_BTCPASS=000000000000000000
APP_BITCOIN_DISABLETLS=true
APP_BITCOIN_NODEKEY=0000000000000000000000000000000000000000000000000000000000000000
APP_BITCOIN_ENDPOINT=/ip4/127.0.0.1/tcp/20001/p2p/16Uiu2HAkwynt59WSsNRS9sk1aszgeQ1PXUS8ax3a3tsewaVMgvZX
APP_BITCOIN_SIGNATUREWEIGHT=1

APP_PARTICLE_URL=https://rpc.particle.network/evm-chain
APP_PARTICLE_CHAINID=1123
APP_PARTICLE_PROJECTUUID=000000000000000000
APP_PARTICLE_PROJECTKEY=000000000000000000
APP_PARTICLE_AAPUBKEYAPI=https://bridge-aa-dev.bsquared.network
```

builder.env

```
APP_LOG_LEVEL=6
APP_PARTICLE_URL=https://rpc.particle.network/evm-chain
APP_PARTICLE_CHAINID=1123
APP_PARTICLE_PROJECTUUID=0000000000000000000000000000000000000000
APP_PARTICLE_PROJECTKEY=0000000000000000000000000000000000000000
APP_PARTICLE_AAPUBKEYAPI=https://bridge-aa-dev.bsquared.network

APP_DATABASE_USERNAME=root
APP_DATABASE_PASSWORD=123456
APP_DATABASE_HOST=127.0.0.1
APP_DATABASE_PORT=3306
APP_DATABASE_DBNAME=b2_message
APP_DATABASE_LOGLEVEL=4

APP_BITCOIN_NAME=bitcoin
APP_BITCOIN_STATUS=true
APP_BITCOIN_CHAINTYPE=2
APP_BITCOIN_CHAINID=0
APP_BITCOIN_MAINNET=false
APP_BITCOIN_RPCURL=127.0.0.1:8085
APP_BITCOIN_SAFEBLOCKNUMBER=3
APP_BITCOIN_LISTENADDRESS=muGFcyjuyURJJsXaLXHCm43jLBmGPPU7ME
APP_BITCOIN_BLOCKINTERVAL=6000
APP_BITCOIN_TOCHAINID=1123
APP_BITCOIN_TOCONTRACTADDRESS=0x0000000000000000000000000000000000000000
APP_BITCOIN_BTCUSER=test
APP_BITCOIN_BTCPASS=test
APP_BITCOIN_DISABLETLS=false

APP_BSQUARED_NAME=bsquared
APP_BSQUARED_STATUS=true
APP_BSQUARED_CHAINTYPE=1
APP_BSQUARED_MAINNET=false
APP_BSQUARED_CHAINID=1123
APP_BSQUARED_RPCURL=127.0.0.1:8084
APP_BSQUARED_SAFEBLOCKNUMBER=1
APP_BSQUARED_LISTENADDRESS=0x0000000000000000000000000000000000000000
APP_BSQUARED_BLOCKINTERVAL=2000
APP_BSQUARED_BUILDERS=0x0000000000000000000000000000000000000000000000000000000000000000

APP_ARBITRUM_NAME=arbitrum
APP_ARBITRUM_STATUS=true
APP_ARBITRUM_CHAINTYPE=1
APP_ARBITRUM_CHAINID=421614
APP_ARBITRUM_MAINNET=false
APP_ARBITRUM_RPCURL=127.0.0.1:8083
APP_ARBITRUM_SAFEBLOCKNUMBER=1
APP_ARBITRUM_LISTENADDRESS=0x0000000000000000000000000000000000000000
APP_ARBITRUM_BLOCKINTERVAL=100
APP_ARBITRUM_BUILDERS=0x0000000000000000000000000000000000000000000000000000000000000000
```

### Quick start

Start using configuration files:

```
$ ./listener -f=listener.yaml
$ ./proposer -f=proposer.yaml
$ ./validator -f=validator.yaml
$ ./builder -f=builder.yaml
```

Start by specifying environment variables:

```
$ APP_LOG_LEVEL=6 ./listener -f=listener.yaml
$ APP_LOG_LEVEL=6 ./proposer -f=proposer.yaml
$ APP_LOG_LEVEL=6 ./validator -f=validator.yaml
$ APP_LOG_LEVEL=6 ./builder -f=builder.yaml
```

Please modify according to the specific configuration and start the service.