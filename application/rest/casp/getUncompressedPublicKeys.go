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
	"net/http"

	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/rest/responses/casp"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
)

func GetUncompressedTMPublicKeys() (casp.UncompressedPublicKeysResponse, error) {
	return getUncompressedPublicKeys(configuration.GetAppConfig().Tendermint.CoinType)
}

const ethCoinType = 60

func GetUncompressedEthPublicKeys() (casp.UncompressedPublicKeysResponse, error) {
	return getUncompressedPublicKeys(ethCoinType)
}

func getUncompressedPublicKeys(coinType uint32) (casp.UncompressedPublicKeysResponse, error) {
	var response casp.UncompressedPublicKeysResponse

	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			// nolint we might like to skip it by purpose
			// nolint: gosec
			InsecureSkipVerify: configuration.GetAppConfig().CASP.TLSInsecureSkipVerify,
		},
	}}

	request, err := http.NewRequest("GET", fmt.Sprintf("%s/casp/api/v1.0/mng/vaults/%s/coins/%d/accounts/0/chains/all/addresses?encoding=uncompressed", configuration.GetAppConfig().CASP.URL, configuration.GetAppConfig().CASP.VaultID, coinType), http.NoBody)
	if err != nil {
		return response, err
	}

	request.Header.Set("authorization", configuration.GetAppConfig().CASP.APIToken)

	var resp *http.Response

	resp, err = client.Do(request)
	if err != nil {
		return response, err
	}

	defer func(body io.Closer) {
		innerErr := body.Close()
		if innerErr != nil {
			logging.Error(fmt.Errorf("%w: %v", ErrCASPUncompressedKeys, innerErr))
		}
	}(resp.Body)

	// Read the response body
	var body []byte

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(body, &response)

	return response, err
}
