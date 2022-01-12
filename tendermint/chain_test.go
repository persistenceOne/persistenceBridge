//go:build integration

/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package tendermint

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/persistenceOne/persistenceBridge/utilities/test"
)

func TestInitializeAndStartChain(t *testing.T) {
	configuration.SetConfig(test.GetCmdWithConfig())

	initAndStartChain, err := InitializeAndStartChain("336h", constants.DefaultPBridgeHome())
	require.Nil(t, err)

	require.NotNil(t, initAndStartChain)

	re := regexp.MustCompile(`^cosmos$`)
	require.Equal(t, true, re.MatchString(initAndStartChain.AccountPrefix))
}
