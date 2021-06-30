package casp

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/rest/responses/casp"
	"io/ioutil"
	"net/http"
)

type signDataRequest struct {
	DataToSign   []string `json:"dataToSign"`
	Description  string   `json:"description"`
	ProviderData string   `json:"providerData"`
	Details      string   `json:"details"`
	PublicKeys   []string `json:"publicKeys"`
}

func SignData(dataToSign []string, publicKeys []string) (casp.PostSignDataResponse, error) {
	var response casp.PostSignDataResponse
	//Encode the data
	postBody, _ := json.Marshal(signDataRequest{
		DataToSign:   dataToSign,
		Description:  "",
		ProviderData: "",
		Details:      "",
		PublicKeys:   publicKeys,
	})
	responseBody := bytes.NewBuffer(postBody)
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}}
	request, err := http.NewRequest("POST", fmt.Sprintf("%s/casp/api/v1.0/mng/vaults/%s/sign", configuration.GetAppConfig().CASP.URL, configuration.GetAppConfig().CASP.VaultID), responseBody)
	if err != nil {
		return response, err
	}
	request.Header.Set("authorization", configuration.GetAppConfig().CASP.APIToken)
	request.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(request)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, fmt.Errorf("unmarshalling error for casp signing %s", err)
	}
	return response, nil
}
