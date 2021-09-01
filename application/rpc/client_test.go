package rpc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestAddValidator(t *testing.T) {
	validatorAddress, err := sdk.ValAddressFromBech32("cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf")
	dirName, err := os.UserHomeDir()

	require.Equal(t, nil, err)
	validatorName := "binance"
	rpcEndpoint := "127.0.0.1:4040"

	database, err := db.OpenDB(filepath.Join(dirName, "/persistence/persistenceBridge/application") + "/db")
	defer database.Close()

	require.Equal(t, nil, err)

	go StartServer(rpcEndpoint)

	require.Equal(t, nil, err)
	validators, err2 := AddValidator(db.Validator{
		Address: validatorAddress,
		Name:   validatorName ,
	}, rpcEndpoint)

	validatorsGet, err2 := db.GetValidators()
	require.Equal(t, validators, validatorsGet)
	require.Equal(t, nil, err2)

}

func TestRemoveValidator(t *testing.T) {
	validatorAddress, err := sdk.ValAddressFromBech32("cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf")
	dirName, err := os.UserHomeDir()

	require.Equal(t, nil, err)
	validatorName := "binance"
	rpcEndpoint := "127.0.0.1:4040"

	database, err := db.OpenDB(filepath.Join(dirName, "/persistence/persistenceBridge/application") + "/db")
	defer database.Close()

	require.Equal(t, nil, err)

	validators, err2 := AddValidator(db.Validator{
		Address: validatorAddress,
		Name:   validatorName ,
	}, rpcEndpoint)

	validators, err2 = RemoveValidator(validatorAddress, rpcEndpoint)
	require.Equal(t, nil, err2)

	validatorsGet, err2 := db.GetValidators()

	require.Equal(t, nil, err2)
	require.Equal(t, validators, validatorsGet)
}

func TestShowValidators(t *testing.T) {
	rpcEndpoint := "127.0.0.1:4040"
	dirName, err := os.UserHomeDir()

	database, err := db.OpenDB(filepath.Join(dirName, "/persistence/persistenceBridge/application") + "/db")
	defer database.Close()

	require.Equal(t, nil, err)

	validators, err2 := ShowValidators("", rpcEndpoint)
	require.Equal(t, nil, err2)

	validatorsGet, err2 := db.GetValidators()

	require.Equal(t, nil, err2)
	require.Equal(t, validators, validatorsGet)

}
