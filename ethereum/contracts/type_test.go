package contracts

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/persistenceOne/persistenceBridge/application/casp"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/ethereum/abi/liquidStaking"
	test "github.com/persistenceOne/persistenceBridge/utilities/testing"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestContracts(t *testing.T) {
	test.SetTestConfig()
	tmAddress, err := casp.GetTendermintAddress()
	require.Equal(t, nil, err)
	ethAddress, err := casp.GetEthAddress()
	require.Equal(t, nil, err)
	configuration.SetCASPAddresses(tmAddress, ethAddress)

	contract := LiquidStaking
	contract.SetAddress(common.HexToAddress(configuration.GetAppConfig().Ethereum.LiquidStakingAddress))

	require.Equal(t, "LIQUID_STAKING", contract.GetName())
	require.Equal(t, common.HexToAddress(configuration.GetAppConfig().Ethereum.LiquidStakingAddress), contract.GetAddress())
	require.Equal(t, abi.ABI{}, contract.GetABI())
	contract.SetABI(liquidStaking.LiquidStakingMetaData.ABI)
	contractABI, err := abi.JSON(strings.NewReader(liquidStaking.LiquidStakingMetaData.ABI))
	require.Equal(t, nil, err)
	require.Equal(t, contractABI, contract.GetABI())
	i := 0
	for k := range contract.GetSDKMsgAndSender() {
		if i == 1 {
			require.Equal(t, "unStake", k)
		} else {
			require.Equal(t, "stake", k)
		}
		i += 1
	}

	// TODO Need correct tx hash of stake tx of LiquidStaking contract in Ropsten
	//ethereumClient, err := ethclient.Dial(configuration.GetAppConfig().Ethereum.EthereumEndPoint)
	//require.Equal(t, nil, err)

	// Test tx in block interrupted
	//ctx, _ := context.WithCancel(context.Background())
	//tx, _, _ := ethereumClient.TransactionByHash(ctx, common.HexToHash("0x8e08d80c37c884467b9b48a77e658711615a5cfde43f95fccfb3b95ee66cd6ea"))

	//method, _, err := contract.GetMethodAndArguments(tx.Data())
	//require.Equal(t, nil, err)
	//require.Equal(t, "stake", method.Name)

}
