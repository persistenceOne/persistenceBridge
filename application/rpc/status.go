package rpc

import (
	"encoding/json"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"net/http"
)

type statusResponse struct {
	EthBlockHeight   int64
	TendermintHeight int64
	NextEpoch        int64
}

type errorResponse struct {
	Message string
}

func status(w http.ResponseWriter, r *http.Request) {
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

	status := statusResponse{
		EthBlockHeight:   ethStatus.LastCheckHeight,
		TendermintHeight: cosmosStatus.LastCheckHeight,
		NextEpoch:        unboundEpoch.Epoch,
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
