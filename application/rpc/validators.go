/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package rpc

import (
	"encoding/json"
	"net/http"

	"github.com/persistenceOne/persistenceBridge/application/db"
)

type validatorResponse struct {
	Validators []db.Validator
}

func validators(w http.ResponseWriter, r *http.Request) {
	var errResponse errorResponse

	validators, err := db.GetValidators()
	if err != nil {
		errResponse.Message = err.Error()

		b, err := json.Marshal(errResponse)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		_, err = w.Write(b)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		return
	}

	validatorResponse := validatorResponse{Validators: validators}

	b, err := json.Marshal(validatorResponse)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	_, err = w.Write(b)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	return
}
