package casp

import (
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	test "github.com/persistenceOne/persistenceBridge/utilities/testing"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetUncompressedPublicKeys(t *testing.T) {
	configuration.InitConfig()
	configuration.SetConfig(test.GetCmdWithConfig())
	funcResponse, err := GetUncompressedTMPublicKeys()
	require.Equal(t, nil, err)
	require.Equal(t, funcResponse.AccountName, "tendermint")
	require.Equal(t, 1, len(funcResponse.PublicKeys))
	funcResponse, err = GetUncompressedEthPublicKeys()
	require.Equal(t, nil, err)
	require.Equal(t, funcResponse.AccountName, "ethereum")
	require.Equal(t, 1, len(funcResponse.PublicKeys))
}
