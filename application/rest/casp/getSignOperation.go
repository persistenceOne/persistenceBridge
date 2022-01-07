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

func GetSignOperation(operationID string) (casp.SignOperationResponse, error) {
	var response casp.SignOperationResponse

	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			// nolint we might like to skip it by purpose
			// nolint: gosec
			// #nosec
			InsecureSkipVerify: configuration.GetAppConfig().CASP.TLSInsecureSkipVerify,
		},
	}}

	request, err := http.NewRequest("GET", fmt.Sprintf("%s/casp/api/v1.0/mng/operations/sign/%s", configuration.GetAppConfig().CASP.URL, operationID), http.NoBody)
	if err != nil {
		return response, err
	}

	request.Header.Set("authorization", configuration.GetAppConfig().CASP.APIToken)

	resp, err := client.Do(request)
	if err != nil {
		return response, err
	}

	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		logging.Error("getting casp sign operation error: ", err, "Body:", string(body))

		var errResponse casp.ErrorResponse

		err = json.Unmarshal(body, &errResponse)
		if err != nil {
			logging.Error("CASP SignOperation ERROR (Unknown response struct type):", err)

			return response, err
		}

		return response, fmt.Errorf("%w: %s", ErrCASPSign, errResponse.Title)
	}

	return response, err
}
