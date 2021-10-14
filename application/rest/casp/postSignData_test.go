package casp

import (
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	test "github.com/persistenceOne/persistenceBridge/utilities/testing"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSignData(t *testing.T) {
	test.SetTestConfig()
	dataToSign := []string{"55C53F5D490297900CEFA825D0C8E8E9532EE8A118ABE7D8570762CD38BE9818"}
	description := "eth"
	publicKeys := []string{configuration.GetAppConfig().CASP.TendermintPublicKey}
	_, err := PostSignData(dataToSign, publicKeys, description)
	require.Equal(t, nil, err)
}
