//go:build integration

/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package casp

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/utilities/test"
)

func TestGetSignOperation(t *testing.T) {
	configuration.SetConfig(test.GetCmdWithConfig())

	ctx := context.Background()

	operationID := "dba4017b-6e88-4693-8c4c-372d283534ad"
	responseRecieved, err := GetSignOperation(ctx, operationID)
	require.Nil(t, err)
	require.Equal(t, responseRecieved.StatusText, "Completed")
	require.Equal(t, responseRecieved.IsApproved, true)

	_, err = GetSignOperation(ctx, "dba4017b-6e88-4693-8c4c-372d283534ae")
	require.Equal(t, "Operation not found", err.Error())
}
