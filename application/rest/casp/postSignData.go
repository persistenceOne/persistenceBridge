package casp

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/persistenceOne/persistenceBridge/application/rest/responses/casp"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
	"io/ioutil"
	"net/http"
)

type SignDataRequest struct {
	DataToSign              []string `json:"dataToSign"`
	Description             string   `json:"description"`
	ProviderData            string   `json:"providerData"`
	Details                 string   `json:"details"`
	PublicKeys              []string `json:"publicKeys"`
	AllowConcurrentKeyUsage bool     `json:"allowConcurrentKeyUsage"`
}

func SignData(dataToSign []string, publicKeys []string, description string) (casp.PostSignDataResponse, bool, error) {
	var response casp.PostSignDataResponse
	//Encode the data
	postBody, _ := json.Marshal(SignDataRequest{
		DataToSign:              dataToSign,
		Description:             description,
		ProviderData:            "",
		Details:                 "",
		PublicKeys:              publicKeys,
		AllowConcurrentKeyUsage: configuration.GetAppConfig().CASP.AllowConcurrentKeyUsage,
	})
	responseBody := bytes.NewBuffer(postBody)
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
		},
	}}
	request, err := http.NewRequest("POST", fmt.Sprintf("%s/casp/api/v1.0/mng/vaults/%s/sign", configuration.GetAppConfig().CASP.URL, configuration.GetAppConfig().CASP.VaultID), responseBody)
	if err != nil {
		return response, false, err
	}
	request.Header.Set("authorization", configuration.GetAppConfig().CASP.GetAPIToken())
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
	if err != nil || response.OperationID == "" {
		logging.Error("posting casp sign data, err:", err, "Body:", string(body))
		var errResponse casp.ErrorResponse
		err = json.Unmarshal(body, &errResponse)
		if err != nil {
			logging.Error("CASP SIGNING ERROR (Unknown response type):", err)
			return response, false, err
		}
		if errResponse.Title == constants.VAULT_BUSY {
			return response, true, nil
		} else {
			return response, false, errResponse
		}
	}
	return response, false, nil
}
