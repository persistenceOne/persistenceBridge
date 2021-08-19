package casp

import (
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	test "github.com/persistenceOne/persistenceBridge/utilities/testing"
	"github.com/stretchr/testify/require"
	"testing"
)


func TestGetUncompressedPublicKeys(t *testing.T) {

	configuration.InitConfig()
	appconfig := test.GetCmdWithConfig()
	configuration.SetConfig(&appconfig)
	funcResponse, err := getUncompressedPublicKeys(118)
	require.Equal(t, nil, err)
	require.Equal(t, funcResponse.AccountName, "tendermint")
	require.Equal(t, 1, len(funcResponse.PublicKeys))

}
