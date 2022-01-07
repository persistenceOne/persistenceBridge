/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package casp

import (
	test "github.com/persistenceOne/persistenceBridge/utilities/testing"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetUncompressedPublicKeys(t *testing.T) {
	test.SetTestConfig()
	funcResponse, err := GetUncompressedTMPublicKeys()
	require.Equal(t, nil, err)
	require.Equal(t, 1, len(funcResponse.Items))
	funcResponse, err = GetUncompressedEthPublicKeys()
	require.Equal(t, nil, err)
	require.Equal(t, 1, len(funcResponse.Items))
}
