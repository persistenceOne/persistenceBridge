/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package casp

import (
	"context"
	"fmt"
	"time"

	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	caspQueries "github.com/persistenceOne/persistenceBridge/application/rest/casp"
	caspResponses "github.com/persistenceOne/persistenceBridge/application/rest/responses/casp"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
)

const (
	tendermintDescription = "tm"
	ethereumDescription   = "eth"
)

// SendDataToSign returns operation id from CASP
func SendDataToSign(ctx context.Context, dataToSign []string, publicKeys []string, isEthTx bool) (string, error) {
	description := tendermintDescription

	if isEthTx {
		description = ethereumDescription
	}

	signDataResponse, err := caspQueries.PostSignData(ctx, dataToSign, publicKeys, description)
	if err != nil {
		return "", err
	}

	return signDataResponse.OperationID, nil
}

func GetCASPSignature(ctx context.Context, operationID string) (caspResponses.SignOperationResponse, error) {
	if operationID == "" {
		return caspResponses.SignOperationResponse{}, ErrEmptyOperationID
	}

	attempts := uint(0)

	for {
		time.Sleep(configuration.GetAppConfig().CASP.WaitTime)

		signOperationResponse, err := caspQueries.GetSignOperation(ctx, operationID)
		if err != nil {
			logging.Error("CASP sign operation:", operationID, " Error:", err)

			return signOperationResponse, err
		}

		attempts++

		if signOperationResponse.Status == constants.Pending {
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
