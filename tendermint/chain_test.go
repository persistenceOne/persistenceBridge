package tendermint

import (
	"github.com/cosmos/relayer/relayer"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	test "github.com/persistenceOne/persistenceBridge/utilities/testing"
	"github.com/stretchr/testify/require"
	"reflect"
	"regexp"
	"testing"
)

func TestInitializeAndStartChain(t *testing.T) {
	configuration.InitConfig()
	configuration.SetConfig(test.GetCmdWithConfig())
	initAndStartChain, err := InitializeAndStartChain("336h", constants.DefaultPBridgeHome)
	require.Equal(t, nil, err)
	re := regexp.MustCompile(`^cosmos$`)
	require.Equal(t, true, re.MatchString(initAndStartChain.AccountPrefix))
	require.Equal(t, reflect.TypeOf(&relayer.Chain{}), reflect.TypeOf(initAndStartChain))
	require.NotNil(t, initAndStartChain)
}
