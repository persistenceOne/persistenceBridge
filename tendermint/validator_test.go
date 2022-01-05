//go:build integration

package tendermint

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/relayer/relayer"
	"github.com/stretchr/testify/require"

	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/utilities/test"
)

func TestHandleSlashedOrAboutToBeSlashed(t *testing.T) {
	configuration.SetConfig(test.GetCmdWithConfig())

	const (
		ethStart = 4772131
		tmStart  = 1
	)

	database, err := db.InitializeDB(constants.TestHomeDir, tmStart, ethStart)
	require.Nil(t, err)

	defer func() {
		err = database.Close()
		require.Nil(t, err)
	}()

	var chain *relayer.Chain
	chain, err = InitializeAndStartChain("336h", constants.DefaultPBridgeHome)
	require.Nil(t, err)

	var address types.ValAddress
	address, err = types.ValAddressFromBech32("cosmosvaloper1efz2js35e4kncmzjmnnu9tul45k8r9etwmkpcp")
	require.Nil(t, err)

	err = db.SetValidator(database, db.Validator{
		Address: address,
		Name:    "test1",
	})
	require.Equal(t, nil, err)

	var cosmosStatus db.Status
	cosmosStatus, err = db.GetCosmosStatus(database)
	require.Equal(t, nil, err)

	processHeight := cosmosStatus.LastCheckHeight + 1
	CheckValidators(chain, database, processHeight)
}
