package casp

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/rest/responses/casp"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"testing"
)


func Test_getUncompressedPublicKeys(t *testing.T) {
	pStakeConfig := configuration.InitConfig()
	_, err := toml.DecodeFile(filepath.Join("/Users/gokuls/.persistenceBridge/", "config.toml"), &pStakeConfig)
	if err != nil {
		log.Fatalf("Error decoding pStakeConfig file: %v\n", err.Error())
	}

	var response casp.UncompressedPublicKeysResponse
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}}

	request, err := http.NewRequest("GET", fmt.Sprintf("%s/casp/api/v1.0/mng/vaults/%s/coins/%d/accounts/0/chains/all/addresses?encoding=uncompressed", configuration.GetAppConfig().CASP.URL, configuration.GetAppConfig().CASP.VaultID, 118), nil)
	require.Equal(t, nil, err)
	request.Header.Set("authorization", configuration.GetAppConfig().CASP.APIToken)
	resp, err := client.Do(request)
	require.Equal(t, nil, err)

	defer resp.Body.Close()
	//Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	require.Equal(t, nil, err)

	err = json.Unmarshal(body, &response)

	funcResponse, err := getUncompressedPublicKeys(118)

	require.Equal(t, funcResponse, response)

	request, err = http.NewRequest("GET", fmt.Sprintf("%s/casp/api/v1.0/mng/vaults/%s/coins/%d/accounts/0/chains/all/addresses?encoding=uncompressed", configuration.GetAppConfig().CASP.URL, configuration.GetAppConfig().CASP.VaultID, 60), nil)
	require.Equal(t, nil, err)
	request.Header.Set("authorization", configuration.GetAppConfig().CASP.APIToken)
	resp, err = client.Do(request)
	require.Equal(t, nil, err)

	defer resp.Body.Close()
	//Read the response body
	body, err = ioutil.ReadAll(resp.Body)
	require.Equal(t, nil, err)

	err = json.Unmarshal(body, &response)

	funcResponse, err = getUncompressedPublicKeys(60)

	require.Equal(t, funcResponse, response)
}
