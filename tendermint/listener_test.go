package tendermint

import (
	"context"
	"github.com/stretchr/testify/require"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestGetAllTxResults(t *testing.T) {
	homedir, _ := os.UserHomeDir()
	fileName := strings.Join([]string{homedir,"/.persistenceBridge/chain.json"},"")
	chain, _ := InitializeAndStartChain(fileName, "336h", homedir)
	ctx := context.Background()
	result, err := getAllTxResults(ctx, chain,0 )
	if err != nil {
		t.Errorf("Error getting all Tx Results: %v",err)
	}
	require.Equal(t, reflect.TypeOf([]*coretypes.ResultTx{}) , reflect.TypeOf(result))
}
