/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package casp

import (
	"fmt"
	"time"

	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	caspQueries "github.com/persistenceOne/persistenceBridge/application/rest/casp"
	caspResponses "github.com/persistenceOne/persistenceBridge/application/rest/responses/casp"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
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
	attempts := 0
	for {
		time.Sleep(configuration.GetAppConfig().CASP.SignatureWaitTime)

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
			if attempts >= configuration.GetAppConfig().CASP.MaxGetSignatureAttempts {
				return signOperationResponse, fmt.Errorf("unable to get approvals for operation: %s", operationID)
			}
			continue
		}
		return signOperationResponse, nil
	}
}
