/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package tendermint

import (
	"reflect"
	"regexp"
	"testing"

	"github.com/cosmos/relayer/relayer"
	"github.com/stretchr/testify/require"

	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	test "github.com/persistenceOne/persistenceBridge/utilities/testing"
)

func TestInitializeAndStartChain(t *testing.T) {
	configuration.InitConfig()
	configuration.SetConfig(test.GetCmdWithConfig())
	initAndStartChain, err := InitializeAndStartChain("336h", constants.DefaultPBridgeHome)
	require.Nil(t, err)
	re := regexp.MustCompile(`^cosmos$`)
	require.Equal(t, true, re.MatchString(initAndStartChain.AccountPrefix))
	require.Equal(t, reflect.TypeOf(&relayer.Chain{}), reflect.TypeOf(initAndStartChain))
	require.NotNil(t, initAndStartChain)
}
