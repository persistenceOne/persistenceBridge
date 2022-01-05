/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package rpc

import (
	"encoding/json"
	"net/http"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dgraph-io/badger/v3"

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

func newStatusHandler(database *badger.DB) http.HandlerFunc {
	return func(writer http.ResponseWriter, _ *http.Request) {
		status(database, writer)
	}
}

func status(database *badger.DB, w http.ResponseWriter) {
	var errResponse errorResponse

	cosmosStatus, err := db.GetCosmosStatus(database)
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

	ethStatus, err = db.GetEthereumStatus(database)
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

	unboundEpoch, err = db.GetUnboundEpochTime(database)
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

	totalWrapped, err = db.GetTotalTokensWrapped(database)
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
	}
}
