package contracts

import (
	"context"
	"fmt"
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
	cName := contract.GetName()
	cAddress := contract.GetAddress()
	cABI := contract.GetABI()
	cMethods := contract.GetSDKMsgAndSender()
	configuration.InitConfig()
	appconfig := test.GetCmdWithConfig()
	configuration.SetConfig(&appconfig)


	require.Equal(t, "LIQUID_STAKING", cName)
	require.Equal(t, common.HexToAddress(constants2.LiquidStakingAddress), cAddress)
	require.Equal(t, abi.ABI{}, cABI)
	contract.SetABI(constants2.LiquidStakingABI)
	contractABI, err := abi.JSON(strings.NewReader(constants2.LiquidStakingABI))
	require.Equal(t, nil, err)
	require.Equal(t, contractABI, contract.GetABI())
	i:=0
	for k := range cMethods {
		if i == 1 {
			require.Equal(t, "unStake", k)
		}else{
			require.Equal(t, "stake", k)
		}
		i+=1

	}
	ethereumClient, err := ethclient.Dial("wss://mainnet.infura.io/ws/v3/b21966541db246d398fb31402eec2c14")
	require.Equal(t, nil, err)

	// Test tx in block interupted
	ctx, _ := context.WithCancel(context.Background())
	tx, _, _ := ethereumClient.TransactionByHash(ctx, common.HexToHash("0x8e08d80c37c884467b9b48a77e658711615a5cfde43f95fccfb3b95ee66cd6ea"))

	method, arguments, err := contract.GetMethodAndArguments(tx.Data())
	require.Equal(t, nil, err)
	require.Equal(t, "stake", method.Name)
	fmt.Println(method,arguments)

}

