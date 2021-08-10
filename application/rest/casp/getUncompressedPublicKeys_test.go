package casp

import (
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/utilities/config"
	"github.com/stretchr/testify/require"
	"testing"
)


func Test_getUncompressedPublicKeys(t *testing.T) {

	configuration.InitConfig()
	appconfig := config.SetConfig()
	configuration.SetConfig(&appconfig)
	funcResponse, err := getUncompressedPublicKeys(118)
	require.Equal(t, nil, err)
	require.Equal(t, funcResponse.AccountName, "tendermint")
	require.Equal(t, 1, len(funcResponse.PublicKeys))


}
