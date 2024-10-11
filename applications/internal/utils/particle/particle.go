package particle

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"net/http"
)

type BtcAccount struct {
	ChainId             int    `json:"chainId"`
	IsDeployed          bool   `json:"isDeployed"`
	EoaAddress          string `json:"eoaAddress"`
	FactoryAddress      string `json:"factoryAddress"`
	EntryPointAddress   string `json:"entryPointAddress"`
	SmartAccountAddress string `json:"smartAccountAddress"`
	Owner               string `json:"owner"`
	Name                string `json:"name"`
	Version             string `json:"version"`
	Index               int    `json:"index"`
	BtcPublicKey        string `json:"btcPublicKey"`
}

type Response struct {
	Jsonrpc string       `json:"jsonrpc"`
	Id      int          `json:"id"`
	Result  []BtcAccount `json:"result"`
	Error   *Error       `json:"error"`
	ChainId int          `json:"chainId"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Param struct {
	Name         string `json:"name"`
	Version      string `json:"version"`
	BtcPublicKey string `json:"btcPublicKey"`
}

type Params struct {
	Method  string  `json:"method"`
	Params  []Param `json:"params"`
	Id      int64   `json:"id"`
	Jsonrpc string  `json:"jsonrpc"`
}

func GetBtcAccount(url string, chainId int, projectUuid string, projectKey string, btcPublicKeys []string) ([]BtcAccount, error) {
	URL := fmt.Sprintf("%s?chainId=%d&projectUuid=%s&projectKey=%s", url, chainId, projectUuid, projectKey)
	var params []Param
	for _, btcPublicKey := range btcPublicKeys {
		params = append(params, Param{
			Name:         "BTC",
			Version:      "2.0.0",
			BtcPublicKey: btcPublicKey,
		})
	}
	value, err := json.Marshal(Params{
		Method:  "particle_aa_getBTCAccount",
		Id:      0,
		Params:  params,
		Jsonrpc: "2.0",
	})
	fmt.Println(string(value))
	if err != nil {
		return nil, err
	}
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(string(value)).
		Post(URL)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, errors.New(resp.Status())
	}
	var response Response
	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		return nil, err
	}
	if response.Error != nil && response.Error.Code != 0 {
		return nil, errors.New(response.Error.Message)
	}
	return response.Result, nil
}
