/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package casp

import (
	test "github.com/persistenceOne/persistenceBridge/utilities/testing"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetSignOperation(t *testing.T) {
	test.SetTestConfig()

	operationID := "dba4017b-6e88-4693-8c4c-372d283534ad"
	responseRecieved, err := GetSignOperation(operationID)

	require.Equal(t, nil, err)
	require.Equal(t, responseRecieved.StatusText, "Completed")
	require.Equal(t, responseRecieved.IsApproved, true)
	_, err = GetSignOperation("dba4017b-6e88-4693-8c4c-372d283534ae")
	require.Equal(t, "Resource not found", err.Error())

}
