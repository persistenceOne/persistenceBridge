/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package casp

import (
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	test "github.com/persistenceOne/persistenceBridge/utilities/testing"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetCASPSignature(t *testing.T) {
	test.SetTestConfig()

	dataToSign := []string{"55C53F5D490297900CEFA825D0C8E8E9532EE8A118ABE7D8570762CD38BE9818"}
	operationID, err := SendDataToSign(dataToSign, []string{configuration.GetAppConfig().CASP.EthereumPublicKey}, true)
	require.Nil(t, err, "Error getting OperationId")
	caspSignature, errCS := GetCASPSignature(operationID)
	require.Nil(t, errCS, "Error getting casp Signature")
	require.Equal(t, caspSignature.IsApproved, true)
	require.Equal(t, caspSignature.Description, "eth")
	require.NotEqual(t, caspSignature.Description, "")
	require.NotEqual(t, caspSignature.AccountID, nil)
	require.NotEqual(t, caspSignature.Signatures, "")
	require.NotNil(t, caspSignature.Signatures)
}

func TestSendDataToSign(t *testing.T) {
	dataToSign := []string{"55C53F5D490297900CEFA825D0C8E8E9532EE8A118ABE7D8570762CD38BE9818"}
	test.SetTestConfig()

	publickey := []string{configuration.GetAppConfig().CASP.TendermintPublicKey}

	caspSignatureOperationID, err := SendDataToSign(dataToSign, publickey, false)
	require.Nil(t, err, "Error getting casp signing OperationID")
	require.NotNil(t, caspSignatureOperationID)
	require.NotEqual(t, "", caspSignatureOperationID, "Empty OperationID")

}
