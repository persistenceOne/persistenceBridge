/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package rpc

import (
	"encoding/json"
	"net/http"

	sdk "github.com/cosmos/cosmos-sdk/types"
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

		var b []byte

		b, err = json.Marshal(errResponse)
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

	var ethStatus db.Status

	ethStatus, err = db.GetEthereumStatus()
	if err != nil {
		errResponse.Message = err.Error()

		var b []byte

		b, err = json.Marshal(errResponse)
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

	var unboundEpoch db.UnboundEpochTime

	unboundEpoch, err = db.GetUnboundEpochTime()
	if err != nil {
		errResponse.Message = err.Error()

		var b []byte

		b, err = json.Marshal(errResponse)
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

	var totalWrapped sdk.Int

	totalWrapped, err = db.GetTotalTokensWrapped()
	if err != nil {
		errResponse.Message = err.Error()

		var b []byte

		b, err = json.Marshal(errResponse)
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

	var b []byte

	b, err = json.Marshal(status)
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
