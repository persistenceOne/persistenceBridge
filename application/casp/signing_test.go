package casp

import (
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	test "github.com/persistenceOne/persistenceBridge/utilities/testing"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestGetCASPSignature(t *testing.T) {
	configuration.InitConfig()
	configuration.SetConfig(test.GetCmdWithConfig())

	dataToSign := []string{"55C53F5D490297900CEFA825D0C8E8E9532EE8A118ABE7D8570762CD38BE9818"}
	operationID, err := GetCASPSigningOperationID(dataToSign, []string{configuration.GetAppConfig().CASP.EthereumPublicKey}, "eth")
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

func TestGetCASPSigningOperationID(t *testing.T) {
	dataToSign := []string{"55C53F5D490297900CEFA825D0C8E8E9532EE8A118ABE7D8570762CD38BE9818"}
	description := "60"
	publickey := []string{"3056301006072A8648CE3D020106052B8104000A0342000413109ECEADCBF6122EF44184B207F8C6820E509497792DDFB166BC090A0FB4447CFFCE16BAAF9EC7F57D14C02641B3A6A698614D973ED744E725A85E62535DA4"}

	configuration.InitConfig()
	configuration.SetConfig(test.GetCmdWithConfig())

	caspSignatureOperationID, err := GetCASPSigningOperationID(dataToSign, publickey, description)
	require.Nil(t, err, "Error getting casp sigining OperationID")
	require.NotNil(t, caspSignatureOperationID)
	require.Equal(t, reflect.TypeOf(""), reflect.TypeOf(caspSignatureOperationID))
	require.NotEqual(t, "", caspSignatureOperationID, "Empty OperationID")

}
