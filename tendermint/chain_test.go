/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package tendermint

import (
	"github.com/persistenceOne/persistenceBridge/application/constants"
	test "github.com/persistenceOne/persistenceBridge/utilities/testing"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
)

func TestInitializeAndStartChain(t *testing.T) {
	test.SetTestConfig()
	initAndStartChain, err := InitializeAndStartChain("336h", constants.DefaultPBridgeHome)
	require.Equal(t, nil, err)
	re := regexp.MustCompile(`^cosmos$`)
	require.Equal(t, true, re.MatchString(initAndStartChain.AccountPrefix))
	require.NotNil(t, initAndStartChain)
}
