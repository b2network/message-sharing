package config

import (
	"bsquared.network/message-sharing-applications/internal/enums"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"strconv"
	"strings"
)

type AppConfig struct {
	Log      LogConfig
	Database Database
	Bsquared Blockchain
	Bitcoin  Blockchain
	Arbitrum Blockchain
	Particle Particle
	Bridges  string
}

type LogConfig struct {
	Level uint32
}

type Database struct {
	UserName string
	Password string
	Host     string
	Port     int64
	DbName   string
	LogLevel int64
}

type Blockchain struct {
	Status            bool
	Name              string
	ChainType         enums.ChainType
	ChainId           int64
	RpcUrl            string
	SafeBlockNumber   int64
	ListenAddress     string
	BlockInterval     int64
	Mainnet           bool
	ToChainId         int64
	ToContractAddress string
	BtcUser           string
	BtcPass           string
	DisableTLS        bool
	NodePort          int
	NodeKey           string
	Endpoint          string
	SignatureWeight   int64
	Validators        []string
	Builders          []string
}

type Particle struct {
	AAPubKeyAPI string
	Url         string
	ChainId     int
	ProjectUuid string
	ProjectKey  string
}

func LoadConfig(input string) AppConfig {
	path, filename, suffix := parsePath(input)
	fmt.Printf("path: %s\n filename: %s\n suffix: %s\n", path, filename, suffix)
	v := viper.New()
	v.SetConfigName(filename)
	v.AddConfigPath(path)
	v.SetConfigType(suffix)
	v.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	v.SetEnvKeyReplacer(replacer)

	var config AppConfig

	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	v.SetEnvPrefix("app")

	if err := v.Unmarshal(&config); err != nil {
		panic(err)
	}

	return config
}

func ParseBridges(input string) (map[int64]string, error) {
	bridges := make(map[int64]string)
	pairs := strings.Split(input, ",")
	for _, pair := range pairs {
		parts := strings.Split(pair, ":")
		if len(parts) != 2 {
			return nil, errors.New("parts len err")
		}
		key, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return nil, err
		}
		value := parts[1]
		bridges[key] = value
	}
	return bridges, nil
}

func parsePath(input string) (string, string, string) {
	var path, filename, suffix = ".", "config", "yaml"
	suffix_index := strings.LastIndex(input, ".")
	path_index := strings.LastIndex(input, "/")
	if path_index < suffix_index && suffix_index > -1 {
		suffix = input[suffix_index+1:]
	}
	if path_index > -1 {
		path = input[:path_index]
	}
	if path_index > -1 && suffix_index > -1 && path_index < suffix_index {
		filename = input[path_index+1 : suffix_index]
	} else if path_index > -1 {
		filename = input[path_index+1:]
	}
	return path, filename, suffix
}
