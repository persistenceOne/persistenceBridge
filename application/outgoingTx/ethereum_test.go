package outgoingTx

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/persistenceOne/persistenceBridge/application/casp"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/ethereum/abi/tokenWrapper"
	test "github.com/persistenceOne/persistenceBridge/utilities/testing"
	"math/big"
	"reflect"
	"regexp"
	"strings"

	"github.com/stretchr/testify/require"
	"testing"
)

func TestEthereumWrapToken(t *testing.T) {
	configuration.InitConfig()
	appConfig := test.GetCmdWithConfig()
	configuration.SetConfig(&appConfig)
	
	ethaddress, _ := casp.GetEthAddress()
	wrapTokenMsg := []WrapTokenMsg{{
		Address: ethaddress,
		Amount:  &big.Int{}},
	}

	ethereumClient, errorInClient := ethclient.Dial("wss://mainnet.infura.io/ws/v3/b21966541db246d398fb31402eec2c14")
	require.Nil(t, errorInClient,"Error getting ETH client!")
	ethWrapToken, errEthWT := EthereumWrapToken(ethereumClient,wrapTokenMsg)
	if errEthWT != nil {
		t.Errorf("Failed getting ETH wrap token for: %v",ethaddress)
	}
	re := regexp.MustCompile(`0x[0-9a-fA-F]{64}`)
	require.NotNil(t, ethWrapToken)
	require.Equal(t, true, re.MatchString(ethWrapToken.String()))
	require.Equal(t, reflect.TypeOf(common.Hash{}),reflect.TypeOf(ethWrapToken))
	require.NotEqual(t, nil,ethWrapToken)
	require.Equal(t,32, len(ethWrapToken))
}

func TestSendTxToEth(t *testing.T){
	configuration.InitConfig()
	appConfig := test.GetCmdWithConfig()
	configuration.SetConfig(&appConfig)
	
	ethclientClient, errorInClient := ethclient.Dial("wss://mainnet.infura.io/ws/v3/b21966541db246d398fb31402eec2c14")
	require.Nil(t, errorInClient,"Error getting ETH client!")
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
	re := regexp.MustCompile(`0x[0-9a-fA-F]{64}`)
	require.Equal(t, true, re.MatchString(txToETHhash.String()))
	require.Equal(t, reflect.TypeOf(common.Hash{}),reflect.TypeOf(txToETHhash))
	require.NotNil(t, txToETHhash)
	require.LessOrEqual(t, 0, len(txToETHhash))
	require.Equal(t, 32, len(txToETHhash))
}

func TestGetEthSignature(t *testing.T) {
	configuration.InitConfig()
	appConfig := test.GetCmdWithConfig()
	configuration.SetConfig(&appConfig)
	
	ethclientClient, errorInClient := ethclient.Dial("wss://mainnet.infura.io/ws/v3/b21966541db246d398fb31402eec2c14")
	require.Nil(t, errorInClient,"Error getting ETH client!")
	address, _ := casp.GetEthAddress()
	ctx := context.Background()
	chainID, errC := ethclientClient.ChainID(ctx)
	require.Nil(t, errC,"Error getting chainID")
	gasPrice, errP := ethclientClient.SuggestGasPrice(ctx)
	require.Nil(t, errP,"Error getting gas price")
	nonce, errN := ethclientClient.PendingNonceAt(ctx, ethBridgeAdmin)
	require.Nil(t, errN,"Error getting nonce")
	tokenWrapperABI, errTWABI := abi.JSON(strings.NewReader(tokenWrapper.TokenWrapperABI))
	require.Nil(t, errTWABI,"Error getting token Wrapper ABI!")
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
	require.Nil(t, errorGettingSignature,"Error getting signature response")
	require.Equal(t, reflect.TypeOf([]byte{}),reflect.TypeOf(ethSignature))
	require.Equal(t, reflect.TypeOf(0),reflect.TypeOf(signatureResponse))
	require.Equal(t, 64, len(ethSignature))
	require.NotEqual(t, -1,signatureResponse)
	require.NotNil(t,ethSignature )
}


func TestSetEthBridgeAdmin(t *testing.T) {
	configuration.InitConfig()
	appConfig := test.GetCmdWithConfig()
	configuration.SetConfig(&appConfig)
	
	ethBridgeAdminErro := setEthBridgeAdmin()
	re := regexp.MustCompile(`^0x[0-9a-fA-F]{40}$`)
	require.Equal(t, true,re.MatchString(ethBridgeAdmin.String()))
	require.Nil(t, nil,ethBridgeAdminErro,"Eth Admin setting failed")
	require.NotEqual(t, "0x0000000000000000000000000000000000000000",ethBridgeAdmin,"ETH Bridge Admin alreadu set")
	require.Equal(t, nil,ethBridgeAdminErro)
}
