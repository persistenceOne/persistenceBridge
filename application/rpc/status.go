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

type statusResponse struct {
	EthBlockHeight   int64
	TendermintHeight int64
	NextEpoch        int64
	TotalWrapped     string
}

type errorResponse struct {
	Message string
}

func status(w http.ResponseWriter, _ *http.Request) {
	var errResponse errorResponse

	cosmosStatus, err := db.GetCosmosStatus()
	if err != nil {
		errResponse.Message = err.Error()

		b, err := json.Marshal(errResponse)
		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		_, err = w.Write(b)
		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
		}

		return
	}

	ethStatus, err := db.GetEthereumStatus()
	if err != nil {
		errResponse.Message = err.Error()

		b, err := json.Marshal(errResponse)
		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		_, err = w.Write(b)
		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		return
	}

	unboundEpoch, err := db.GetUnboundEpochTime()
	if err != nil {
		errResponse.Message = err.Error()

		b, err := json.Marshal(errResponse)
		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		_, err = w.Write(b)
		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		return
	}

	totalWrapped, err := db.GetTotalTokensWrapped()
	if err != nil {
		errResponse.Message = err.Error()

		b, err := json.Marshal(errResponse)
		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		_, err = w.Write(b)
		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		return
	}

	status := statusResponse{
		EthBlockHeight:   ethStatus.LastCheckHeight,
		TendermintHeight: cosmosStatus.LastCheckHeight,
		NextEpoch:        unboundEpoch.Epoch,
		TotalWrapped:     totalWrapped.String(),
	}

	b, err := json.Marshal(status)
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	_, err = w.Write(b)
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	return
}
