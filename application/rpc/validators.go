/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package rpc

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
)

type validatorResponse struct {
	Validators []db.Validator
}

func validators(w http.ResponseWriter, _ *http.Request) {
	var errResponse errorResponse

	validators, err := db.GetValidators()
	if err != nil {
		errResponse.Message = err.Error()

		var b []byte

		b, err = json.Marshal(errResponse)
		if err != nil {
			_, httpErr := w.Write([]byte(err.Error()))
			if httpErr != nil {
				logging.Error(fmt.Sprintf("%v: %v, marshall error %v", ErrHTTPWriter, httpErr, err))
			}

			return
		}

		_, err = w.Write(b)
		if err != nil {
			_, httpErr := w.Write([]byte(err.Error()))
			logging.Error(fmt.Sprintf("%v: %v, previous error %v", ErrHTTPWriter, httpErr.Error(), err))

			return
		}

		return
	}

	response := validatorResponse{Validators: validators}

	var b []byte

	b, err = json.Marshal(response)
	if err != nil {
		_, httpErr := w.Write([]byte(err.Error()))
		if httpErr != nil {
			logging.Error(fmt.Sprintf("%v: %v, marshall error %v", ErrHTTPWriter, httpErr, err))
		}

		return
	}

	_, err = w.Write(b)
	if err != nil {
		_, httpErr := w.Write([]byte(err.Error()))
		if httpErr != nil {
			logging.Error(fmt.Sprintf("%v: %v, previous error %v", ErrHTTPWriter, httpErr, err))
		}
	}
}
