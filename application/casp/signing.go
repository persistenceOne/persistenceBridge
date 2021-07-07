package casp

import (
	"fmt"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	caspQueries "github.com/persistenceOne/persistenceBridge/application/rest/casp"
	caspResponses "github.com/persistenceOne/persistenceBridge/application/rest/responses/casp"
	"log"
	"time"
)

func GetCASPSigningOperationID(dataToSign []string, publicKeys []string) (string, error) {
	for {
		signDataResponse, busy, err := caspQueries.SignData(dataToSign, publicKeys)
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
			if err.Error() == constants.OPERATION_ID_NOT_FOUND {
				return caspResponses.SignOperationResponse{}, fmt.Errorf("operation id not found")
			}
			log.Printf("Error while getting sign operation %v\n", err)
			time.Sleep(configuration.GetAppConfig().CASP.SignatureWaitTime)
			continue
		}
		if signOperationResponse.Status == "PENDING" {
			time.Sleep(configuration.GetAppConfig().CASP.SignatureWaitTime)
			continue
		}
		return signOperationResponse, nil
	}
}
