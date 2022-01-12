/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package casp

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/rest/responses/casp"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
)

type SignDataRequest struct {
	DataToSign              []string `json:"dataToSign"`
	Description             string   `json:"description"`
	ProviderData            string   `json:"providerData"`
	Details                 string   `json:"details"`
	PublicKeys              []string `json:"publicKeys"`
	AllowConcurrentKeyUsage bool     `json:"allowConcurrentKeyUsage"`
}

func PostSignData(ctx context.Context, dataToSign, publicKeys []string, description string) (casp.PostSignDataResponse, error) {
	var response casp.PostSignDataResponse

	// Encode the data
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

	request, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/casp/api/v1.0/mng/vaults/%s/sign", configuration.GetAppConfig().CASP.URL, configuration.GetAppConfig().CASP.VaultID), responseBody)
	if err != nil {
		return response, err
	}

	request.Header.Set("authorization", configuration.GetAppConfig().CASP.APIToken)
	request.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(request)
	if err != nil {
		return response, err
	}

	defer func(body io.Closer) {
		err := body.Close()
		if err != nil {
			logging.Error(err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(body, &response)
	if err != nil || response.OperationID == "" {
		logging.Error("posting casp sign data, err:", err, "Body:", string(body))

		var errResponse casp.ResponseError

		err = json.Unmarshal(body, &errResponse)
		if err != nil {
			logging.Error("CASP SIGNING ERROR (Unknown error response type):", err)

			return response, err
		}

		return response, errResponse
	}

	return response, err
}
