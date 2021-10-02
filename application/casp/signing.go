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

// SendDataToSign returns operation id from CASP
func SendDataToSign(dataToSign []string, publicKeys []string, isEthTx bool) (string, error) {
	description := "tm"
	if isEthTx {
		description = "eth"
	}
	signDataResponse, err := caspQueries.PostSignData(dataToSign, publicKeys, description)
	if err != nil {
		return "", err
	}
	return signDataResponse.OperationID, nil
}

func GetCASPSignature(operationID string) (caspResponses.SignOperationResponse, error) {
	if operationID == "" {
		return caspResponses.SignOperationResponse{}, fmt.Errorf("empty operationID")
	}
	attempts := 0
	for {
		time.Sleep(configuration.GetAppConfig().CASP.WaitTime)

		signOperationResponse, err := caspQueries.GetSignOperation(operationID)
		if err != nil {
			logging.Error("CASP sign operation:", operationID, " Error:", err)
			return signOperationResponse, err
		}
		attempts++
		if signOperationResponse.Status == constants.PENDING {
			logging.Info("CASP signing operation pending for", operationID)
			err = signOperationResponse.GetPendingParticipantsApprovals()
			if err != nil {
				logging.Error("attempt:", attempts, err)
			}
			if attempts >= configuration.GetAppConfig().CASP.MaxAttempts {
				return signOperationResponse, fmt.Errorf("unable to get approvals for operation: %s", operationID)
			}
			continue
		}
		return signOperationResponse, nil
	}
}
