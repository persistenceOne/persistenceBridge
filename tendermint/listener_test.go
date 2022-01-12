//go:build integration

/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package tendermint

import (
	"context"
	"testing"

	"github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/stretchr/testify/require"
)

func TestGetAllTxResults(t *testing.T) {
	chain, _ := InitializeAndStartChain("336h", constants.DefaultPBridgeHome())

	ctx := context.Background()

	result, err := getAllTxResults(ctx, chain, 1)
	if err != nil {
		t.Errorf("Error getting all Tx Results: %v", err)
	}

	require.NotNil(t, result[0])
}
