/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package casp

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/rest/responses/casp"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
)

func GetUncompressedTMPublicKeys() (casp.UncompressedPublicKeysResponse, error) {
	return getUncompressedPublicKeys(configuration.GetAppConfig().Tendermint.CoinType)
}

func GetUncompressedEthPublicKeys() (casp.UncompressedPublicKeysResponse, error) {
	return getUncompressedPublicKeys(60)
}

func getUncompressedPublicKeys(coinType uint32) (casp.UncompressedPublicKeysResponse, error) {
	var response casp.UncompressedPublicKeysResponse

	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: configuration.GetAppConfig().CASP.TLSInsecureSkipVerify,
		},
	}}

	request, err := http.NewRequest("GET", fmt.Sprintf("%s/casp/api/v1.0/mng/vaults/%s/coins/%d/accounts/0/chains/all/addresses?encoding=uncompressed", configuration.GetAppConfig().CASP.URL, configuration.GetAppConfig().CASP.VaultID, coinType), nil)
	if err != nil {
		return response, err
	}

	request.Header.Set("authorization", configuration.GetAppConfig().CASP.APIToken)

	var resp *http.Response

	resp, err = client.Do(request)
	if err != nil {
		return response, err
	}

	defer func(Body io.ReadCloser) {
		innerErr := Body.Close()
		if innerErr != nil {
			logging.Error(fmt.Errorf("casp error while getting UncompressedPublicKeys: %v", innerErr))
		}
	}(resp.Body)

	// Read the response body
	var body []byte

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(body, &response)

	return response, err
}
