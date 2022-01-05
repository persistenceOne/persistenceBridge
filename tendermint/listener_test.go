//go:build integration

/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package tendermint

import (
	"context"
	"reflect"
	"testing"

	"github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/stretchr/testify/require"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
)

func TestGetAllTxResults(t *testing.T) {
	chain, err := InitializeAndStartChain("336h", constants.DefaultPBridgeHome)
	require.Nil(t, err)

	ctx := context.Background()

	result, err := getAllTxResults(ctx, chain, 0)
	if err != nil {
		t.Errorf("Error getting all Tx Results: %v", err)
	}

	require.Equal(t, reflect.TypeOf([]*coretypes.ResultTx{}), reflect.TypeOf(result))
}
