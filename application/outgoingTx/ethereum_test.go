package outgoingTx

import (
	"context"
	"github.com/BurntSushi/toml"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/persistenceOne/persistenceBridge/application/casp"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/ethereum/abi/tokenWrapper"
	"math/big"
	"strings"

	"github.com/stretchr/testify/require"
	"log"
	"path/filepath"
	"testing"
)

func TestEthereumWrapToken(t *testing.T) {
	pStakeConfig := configuration.InitConfig()
	_, err := toml.DecodeFile(filepath.Join("/Users/ankitkumar/.persistenceBridge/", "config.toml"), &pStakeConfig)
	if err != nil {
		log.Fatalf("Error decoding pStakeConfig file: %v\n", err.Error())
	}
	ethaddress, _ := casp.GetEthAddress()
	wrapTokenMsg := []WrapTokenMsg{{
		Address: ethaddress,
		Amount:  &big.Int{},},
	}

	ethereumClient, error_in_client := ethclient.Dial(pStakeConfig.Ethereum.EthereumEndPoint)
	if error_in_client != nil {
		t.Errorf("Error getting ETH client!")
	}
	ethWrapToken, error := EthereumWrapToken(ethereumClient,wrapTokenMsg)
	if error != nil {
		t.Errorf("Failed getting ETH wrap token for: %v",ethaddress)
	}
	require.NotNil(t, ethWrapToken)
	require.NotEqual(t, nil,ethWrapToken)
	require.Equal(t,32, len(ethWrapToken))
}

func Test_sendTxToEth(t *testing.T){
	pStakeConfig := configuration.InitConfig()
	_, err := toml.DecodeFile(filepath.Join("/Users/ankitkumar/.persistenceBridge/", "config.toml"), &pStakeConfig)
	if err != nil {
		log.Fatalf("Error decoding pStakeConfig file: %v\n", err.Error())
	}
	ethclientClient, errorInClient := ethclient.Dial(pStakeConfig.Ethereum.EthereumEndPoint)
	if errorInClient != nil {
		t.Errorf("Error getting ETH client!")
	}
	ethaddress, _ := casp.GetEthAddress()
	tokenWrapperABI, err := abi.JSON(strings.NewReader(tokenWrapper.TokenWrapperABI))
	addresses := make([]common.Address, 1)
	amounts := make([]*big.Int, 1)
	addresses[0] = ethaddress
	amounts[0] = big.NewInt(1)
	txdata, error_txdata := tokenWrapperABI.Pack("generateUTokensInBatch", addresses, amounts)
	if error_txdata != nil {
		t.Errorf("Error generating TX data with error: \n %v",error_txdata)
	}
	txToETHhash, err := sendTxToEth(ethclientClient,&ethaddress,nil, txdata)
	if err != nil {
		t.Errorf("Error sending TX to ETH with error: \n %v",err)
	}
	require.NotNil(t, txToETHhash)
	require.LessOrEqual(t, 0, len(txToETHhash))
	require.Equal(t, 32, len(txToETHhash))
}

func Test_getEthSignature(t *testing.T) {
	pStakeConfig := configuration.InitConfig()
	_, err := toml.DecodeFile(filepath.Join("/Users/ankitkumar/.persistenceBridge/", "config.toml"), &pStakeConfig)
	if err != nil {
		log.Fatalf("Error decoding pStakeConfig file: %v\n", err.Error())
	}
	ethclientClient, errorInClient := ethclient.Dial(pStakeConfig.Ethereum.EthereumEndPoint)
	if errorInClient != nil {
		t.Errorf("Error getting ETH client!")
	}
	address, _ := casp.GetEthAddress()
	ctx := context.Background()
	chainID, err := ethclientClient.ChainID(ctx)
	gasPrice, err := ethclientClient.SuggestGasPrice(ctx)
	nonce, err := ethclientClient.PendingNonceAt(ctx, ethBridgeAdmin)
	tokenWrapperABI, err := abi.JSON(strings.NewReader(tokenWrapper.TokenWrapperABI))
	addresses := make([]common.Address, 1)
	amounts := make([]*big.Int, 1)
	addresses[0] = address
	amounts[0] = big.NewInt(1)
	txdata, error_txdata := tokenWrapperABI.Pack("generateUTokensInBatch", addresses, amounts)
	if error_txdata != nil {
		t.Errorf("Error generating TX data with error: \n %v",error_txdata)
	}
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		Value:    nil,
		Gas:      configuration.GetAppConfig().Ethereum.GasLimit,
		GasPrice: gasPrice.Add(gasPrice, big.NewInt(4000000000)),
		Data:     txdata,
		To:       &address,
	})
	signer := types.NewEIP155Signer(chainID)
	ethSignature, signatureResponse, errorGettingSignature := getEthSignature(tx,signer)
	if errorGettingSignature != nil {
		t.Errorf("")
	}
	require.NotEqual(t, -1,signatureResponse)
	require.NotNil(t,ethSignature )
}


func Test_setEthBridgeAdmin(t *testing.T) {
	pStakeConfig := configuration.InitConfig()
	_, err := toml.DecodeFile(filepath.Join("/Users/ankitkumar/.persistenceBridge/", "config.toml"), &pStakeConfig)
	if err != nil {
		log.Fatalf("Error decoding pStakeConfig file: %v\n", err.Error())
	}
	ethBridgeAdminErro := setEthBridgeAdmin()
	require.Nil(t, nil,ethBridgeAdminErro,"Eth Admin setting failed")
}
