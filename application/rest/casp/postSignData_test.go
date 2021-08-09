package casp

import (
	"github.com/BurntSushi/toml"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/stretchr/testify/require"
	"log"
	"path/filepath"
	"testing"
)

func TestSignData(t *testing.T) {
	pStakeConfig := configuration.InitConfig()
	_, err := toml.DecodeFile(filepath.Join("/Users/gokuls/.persistenceBridge/", "config.toml"), &pStakeConfig)
	if err != nil {
		log.Fatalf("Error decoding pStakeConfig file: %v\n", err.Error())
	}
	dataToSign := []string{"55C53F5D490297900CEFA825D0C8E8E9532EE8A118ABE7D8570762CD38BE9818"}
	description := "eth"
	publicKeys := []string{"3056301006072A8648CE3D020106052B8104000A0342000413109ECEADCBF6122EF44184B207F8C6820E509497792DDFB166BC090A0FB4447CFFCE16BAAF9EC7F57D14C02641B3A6A698614D973ED744E725A85E62535DA4"}

	responseRecieved, boolrecieved, err := SignData(dataToSign, publicKeys, description)

	require.Equal(t, nil, err)
	require.Equal(t, false, boolrecieved)
	require.Equal(t, "f5f87fd1-08e2-4d0a-bc9d-cbbb59f1c47b", responseRecieved.OperationID)


}
