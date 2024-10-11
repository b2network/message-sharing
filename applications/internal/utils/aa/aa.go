package aa

import (
	"bsquared.network/message-sharing-applications/internal/utils/particle"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var AddressNotFoundErrCode = "1001"

type Response struct {
	Code    string
	Message string
	Data    struct {
		Pubkey string
	}
}

func GetPubKey(api, btcAddress string) (*Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", api+"/v1/btc/pubkey/"+btcAddress, nil)
	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	btcResp := Response{}

	err = json.Unmarshal(body, &btcResp)
	if err != nil {
		return nil, err
	}

	return &btcResp, nil
}

func BitcoinAddressToEthAddress(aaPubKeyAPI, bitcoinAddress, particleUrl string, particleChainId int, particleProjectUuid, particleProjectKey string) (string, error) {
	pubkeyResp, err := GetPubKey(aaPubKeyAPI, bitcoinAddress)
	if err != nil {
		return "", err
	}
	if pubkeyResp.Code != "0" {
		if pubkeyResp.Code == AddressNotFoundErrCode {
			return "", fmt.Errorf("AAGetBTCAccount not found")
		}
		return "", fmt.Errorf("get pubkey code err:%v", pubkeyResp)
	}

	accounts, err := particle.GetBtcAccount(particleUrl, particleChainId, particleProjectUuid, particleProjectKey, []string{pubkeyResp.Data.Pubkey})
	if err != nil {
		return "", err
	}
	if len(accounts) != 1 {
		return "", fmt.Errorf("AAGetBTCAccount result not match")
	}
	return accounts[0].SmartAccountAddress, nil
}
