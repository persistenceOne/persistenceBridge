package casp

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/persistenceOne/persistenceBridge/application/rest/responses/casp"
	"io/ioutil"
	"log"
	"net/http"
)

type SignDataRequest struct {
	DataToSign   []string `json:"dataToSign"`
	Description  string   `json:"description"`
	ProviderData string   `json:"providerData"`
	Details      string   `json:"details"`
	PublicKeys   []string `json:"publicKeys"`
}

func SignData(dataToSign []string, publicKeys []string, description string) (casp.PostSignDataResponse, bool, error) {
	var response casp.PostSignDataResponse
	//Encode the data
	postBody, _ := json.Marshal(SignDataRequest{
		DataToSign:   dataToSign,
		Description:  description,
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
		return response, false, err
	}
	request.Header.Set("authorization", configuration.GetAppConfig().CASP.APIToken)
	request.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(request)
	if err != nil {
		return response, false, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, false, err
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error posting casp sign data %s\n", string(body))
		var errResponse casp.ErrorResponse
		err = json.Unmarshal(body, &errResponse)
		if err != nil {
			log.Fatalf("CASP SIGNING ERROR (Unknown response type): %v\n", err)
		}
		if errResponse.Title == constants.VAULT_BUSY {
			return response, true, nil
		}
	}
	return response, false, err
}
