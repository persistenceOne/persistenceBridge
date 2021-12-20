/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package outgoingtx

import (
	"context"
	"math/big"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/persistenceOne/persistenceBridge/application/casp"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/ethereum/abi/tokenWrapper"
	test "github.com/persistenceOne/persistenceBridge/utilities/testing"

	"github.com/stretchr/testify/require"
)

func TestEthereumWrapToken(t *testing.T) {
	configuration.InitConfig()
	configuration.SetConfig(test.GetCmdWithConfig())

	ethAddress, _ := casp.GetEthAddress()
	wrapTokenMsg := []WrapTokenMsg{{
		Address: ethAddress,
		Amount:  &big.Int{}},
	}

	ethereumClient, err := ethclient.Dial(configuration.GetAppConfig().Ethereum.EthereumEndPoint)
	require.Nil(t, err)

	var txHash common.Hash
	txHash, err = EthereumWrapToken(ethereumClient, wrapTokenMsg)
	require.NotNil(t, txHash)
	require.Nil(t, err)
	require.Equal(t, reflect.TypeOf(common.Hash{}), reflect.TypeOf(txHash))
	require.NotEqual(t, nil, txHash)
	require.Equal(t, 32, len(txHash))

	re := regexp.MustCompile(`0x[0-9a-fA-F]{64}`)
	require.Equal(t, true, re.MatchString(txHash.String()))
}

func TestSendTxToEth(t *testing.T) {
	configuration.InitConfig()
	configuration.SetConfig(test.GetCmdWithConfig())

	ethClient, errorInClient := ethclient.Dial(configuration.GetAppConfig().Ethereum.EthereumEndPoint)
	require.Nil(t, errorInClient, "Error getting ETH client!")

	ethAddress, _ := casp.GetEthAddress()
	tokenWrapperABI, err := abi.JSON(strings.NewReader(tokenWrapper.TokenWrapperABI))
	require.Nil(t, err)

	addresses := make([]common.Address, 1)
	amounts := make([]*big.Int, 1)
	addresses[0] = ethAddress
	amounts[0] = big.NewInt(1)

	txdata, err := tokenWrapperABI.Pack("generateUTokensInBatch", addresses, amounts)
	require.Nil(t, err)

	txToETHhash, err := sendTxToEth(ethClient, &ethAddress, nil, txdata)
	require.Nil(t, err)

	re := regexp.MustCompile(`0x[0-9a-fA-F]{64}`)
	require.Equal(t, true, re.MatchString(txToETHhash.String()))
	require.Equal(t, reflect.TypeOf(common.Hash{}), reflect.TypeOf(txToETHhash))
	require.NotNil(t, txToETHhash)
	require.LessOrEqual(t, 0, len(txToETHhash))
	require.Equal(t, 32, len(txToETHhash))
}

func TestGetEthSignature(t *testing.T) {
	configuration.InitConfig()
	configuration.SetConfig(test.GetCmdWithConfig())

	ethClient, errorInClient := ethclient.Dial(configuration.GetAppConfig().Ethereum.EthereumEndPoint)
	require.Nil(t, errorInClient, "Error getting ETH client!")

	address, _ := casp.GetEthAddress()
	ctx := context.Background()
	chainID, err := ethClient.ChainID(ctx)
	require.Nil(t, err)

	gasPrice, err := ethClient.SuggestGasPrice(ctx)
	require.Nil(t, err)

	nonce, err := ethClient.PendingNonceAt(ctx, ethBridgeAdmin)
	require.Nil(t, err)

	tokenWrapperABI, err := abi.JSON(strings.NewReader(tokenWrapper.TokenWrapperABI))
	require.Nil(t, err)

	addresses := make([]common.Address, 1)
	amounts := make([]*big.Int, 1)
	addresses[0] = address
	amounts[0] = big.NewInt(1)
	txdata, err := tokenWrapperABI.Pack("generateUTokensInBatch", addresses, amounts)
	require.Nil(t, err)

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		Value:    nil,
		Gas:      configuration.GetAppConfig().Ethereum.GasLimit,
		GasPrice: gasPrice.Add(gasPrice, big.NewInt(4000000000)),
		Data:     txdata,
		To:       &address,
	})

	signer := types.NewEIP155Signer(chainID)
	ethSignature, signatureResponse, errorGettingSignature := getEthSignature(tx, signer)
	require.Nil(t, errorGettingSignature, "Error getting signature response")
	require.Equal(t, reflect.TypeOf([]byte{}), reflect.TypeOf(ethSignature))
	require.Equal(t, reflect.TypeOf(0), reflect.TypeOf(signatureResponse))
	require.Equal(t, 64, len(ethSignature))
	require.NotEqual(t, -1, signatureResponse)
	require.NotNil(t, ethSignature)
}

func TestSetEthBridgeAdmin(t *testing.T) {
	configuration.InitConfig()
	configuration.SetConfig(test.GetCmdWithConfig())

	err := setEthBridgeAdmin()
	require.Nil(t, err)

	re := regexp.MustCompile(`^0x[0-9a-fA-F]{40}$`)
	require.Equal(t, true, re.MatchString(ethBridgeAdmin.String()))
	require.NotEqual(t, EthEmptyAddress, ethBridgeAdmin, "ETH Bridge Admin alreadu set")
}
