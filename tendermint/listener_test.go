package tendermint

import (
	"context"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/stretchr/testify/require"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	"reflect"
	"testing"
)

func TestGetAllTxResults(t *testing.T) {
	chain, _ := InitializeAndStartChain("336h", constants.DefaultPBridgeHome)
	ctx := context.Background()
	result, err := getAllTxResults(ctx, chain, 0)
	if err != nil {
		t.Errorf("Error getting all Tx Results: %v", err)
	}
	require.Equal(t, reflect.TypeOf([]*coretypes.ResultTx{}), reflect.TypeOf(result))
}
