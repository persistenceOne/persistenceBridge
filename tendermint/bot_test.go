package tendermint

import (
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/persistenceOne/persistenceBridge/application/db"
	test "github.com/persistenceOne/persistenceBridge/utilities/testing"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHandleSlashedOrAboutToBeSlashed(t *testing.T) {
	configuration.InitConfig()
	configuration.SetConfig(test.GetCmdWithConfig())
	var ethStart int64 = 4772131
	var tmStart int64 = 1
	database, err := db.InitializeDB(constants.TestHomeDir, tmStart, ethStart)
	chain, _ := InitializeAndStartChain("336h", constants.DefaultPBridgeHome)
	address, err := types.ValAddressFromBech32("cosmosvaloper1efz2js35e4kncmzjmnnu9tul45k8r9etwmkpcp")
	require.Equal(t, nil, err)
	err = db.SetValidator(db.Validator{
		Address: address,
		Name:    "test1",
	})
	require.Equal(t, nil, err)
	validatorList, err := db.GetValidators()
	require.Equal(t, nil, err)
	cosmosStatus, err := db.GetCosmosStatus()
	require.Equal(t, nil, err)
	processHeight := cosmosStatus.LastCheckHeight + 1
	err = handleSlashedOrAboutToBeSlashed(chain, validatorList, processHeight)
	require.Equal(t, nil, err)
	database.Close()
}
