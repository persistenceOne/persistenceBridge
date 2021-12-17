/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package casp

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/persistenceOne/persistenceBridge/application/configuration"
	test "github.com/persistenceOne/persistenceBridge/utilities/testing"
)

func TestGetSignOperation(t *testing.T) {
	configuration.InitConfig()
	configuration.SetConfig(test.GetCmdWithConfig())

	operationID := "69544933-2767-4e09-af4a-c2dacb9a20af"
	responseRecieved, err := GetSignOperation(operationID)
	require.Nil(t, err)
	require.Equal(t, responseRecieved.StatusText, "Completed")
	require.Equal(t, responseRecieved.IsApproved, true)
	require.Equal(t, responseRecieved.AccountID, "bd4c618e-8046-4fef-bdaa-9716ade77553")

	_, err = GetSignOperation("")
	require.Equal(t, "Operation not found", err.Error())
}
