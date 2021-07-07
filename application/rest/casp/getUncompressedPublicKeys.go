package casp

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/persistenceOne/persistenceBridge/application/rest/responses/casp"
	"io/ioutil"
	"net/http"
)

func GetUncompressedTMPublicKeys() (casp.UncompressedPublicKeysResponse, error) {
	return getUncompressedPublicKeys(118)
}

func GetUncompressedEthPublicKeys() (casp.UncompressedPublicKeysResponse, error) {
	return getUncompressedPublicKeys(60)
}

func getUncompressedPublicKeys(coinType uint32) (casp.UncompressedPublicKeysResponse, error) {
	var response casp.UncompressedPublicKeysResponse
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}}

	request, err := http.NewRequest("GET", fmt.Sprintf("%s/casp/api/v1.0/mng/vaults/%s/coins/%d/accounts/0/chains/all/addresses?encoding=uncompressed", "https://65.2.149.241:443", "509fd89a-762a-40ec-bd4b-0745b06e2d3d", coinType), nil)
	if err != nil {
		return response, err
	}
	request.Header.Set("authorization", "Bearer cHVuZWV0TmV3QXBpa2V5MTI6OWM1NDBhMzAtNTQ5NC00ZDdhLTljODktODA3MDZiNWNhYzQ1")
	resp, err := client.Do(request)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()
	//Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}
	err = json.Unmarshal(body, &response)
	return response, err
}
