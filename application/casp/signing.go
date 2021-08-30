package casp

import (
	"fmt"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	caspQueries "github.com/persistenceOne/persistenceBridge/application/rest/casp"
	caspResponses "github.com/persistenceOne/persistenceBridge/application/rest/responses/casp"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
	"time"
)

// GetCASPSigningOperationID description should be small
func GetCASPSigningOperationID(dataToSign []string, publicKeys []string, description string) (string, error) {
	for {
		signDataResponse, busy, err := caspQueries.SignData(dataToSign, publicKeys, description)
		if busy {
			time.Sleep(configuration.GetAppConfig().CASP.SignatureWaitTime)
		}
		if err != nil {
			return "", err
		}
		return signDataResponse.OperationID, nil
	}
}

func GetCASPSignature(operationID string) (caspResponses.SignOperationResponse, error) {
	if operationID == "" {
		return caspResponses.SignOperationResponse{}, fmt.Errorf("empty operationID")
	}
	for {
		signOperationResponse, err := caspQueries.GetSignOperation(operationID)
		if err != nil {
			logging.Error("CASP sign operation:", operationID, " Error:", err)
			return caspResponses.SignOperationResponse{}, err
		}
		if signOperationResponse.Status == constants.PENDING {
			logging.Info("CASP signing operation pending for", operationID)
			time.Sleep(configuration.GetAppConfig().CASP.SignatureWaitTime)
			continue
		}
		return signOperationResponse, nil
	}
}
