/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package contracts

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	constants2 "github.com/persistenceOne/persistenceBridge/application/constants"
	test "github.com/persistenceOne/persistenceBridge/utilities/testing"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestContracts(t *testing.T) {
	contract := LiquidStaking
	contractName := contract.GetName()
	contractAddress := contract.GetAddress()
	cABI := contract.GetABI()
	cMethods := contract.GetSDKMsgAndSender()
	configuration.InitConfig()
	configuration.SetConfig(test.GetCmdWithConfig())

	require.Equal(t, "LIQUID_STAKING", contractName)
	require.Equal(t, common.HexToAddress(constants2.LiquidStakingAddress), contractAddress)
	require.Equal(t, abi.ABI{}, cABI)
	contract.SetABI(constants2.LiquidStakingABI)
	contractABI, err := abi.JSON(strings.NewReader(constants2.LiquidStakingABI))
	require.Equal(t, nil, err)
	require.Equal(t, contractABI, contract.GetABI())
	i := 0
	for k := range cMethods {
		if i == 1 {
			require.Equal(t, "unStake", k)
		} else {
			require.Equal(t, "stake", k)
		}
		i += 1

	}
	ethereumClient, err := ethclient.Dial(configuration.GetAppConfig().Ethereum.EthereumEndPoint)
	require.Equal(t, nil, err)

	// Test tx in block interupted
	ctx, _ := context.WithCancel(context.Background())
	tx, _, _ := ethereumClient.TransactionByHash(ctx, common.HexToHash("0x8e08d80c37c884467b9b48a77e658711615a5cfde43f95fccfb3b95ee66cd6ea"))

	method, _, err := contract.GetMethodAndArguments(tx.Data())
	require.Equal(t, nil, err)
	require.Equal(t, "stake", method.Name)

}
