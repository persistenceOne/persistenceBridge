package ethereum

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/persistenceOne/persistenceBridge/application/casp"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/ethereum/contracts"
	test "github.com/persistenceOne/persistenceBridge/utilities/testing"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_collectEthTx(t *testing.T) {
	configuration.InitConfig()
	appconfig := test.GetCmdWithConfig()
	configuration.SetConfig(&appconfig)
	tmAddress, err := casp.GetTendermintAddress()
	require.Equal(t, nil, err)

	configuration.SetPStakeAddress(tmAddress)

	ethereumClient, err := ethclient.Dial("wss://mainnet.infura.io/ws/v3/b21966541db246d398fb31402eec2c14")
	require.Equal(t, nil, err)
	ctx := context.Background()
	tx, _, _ := ethereumClient.TransactionByHash(ctx, common.HexToHash("8e08d80c37c884467b9b48a77e658711615a5cfde43f95fccfb3b95ee66cd6ea"))
	contract := contracts.LiquidStaking
	coltx, err := collectEthTx(ethereumClient, &ctx, tx, &contract)
	require.Equal(t, nil, err)
	require.Equal(t, "0x8e08d80c37c884467b9b48a77e658711615a5cfde43f95fccfb3b95ee66cd6ea", coltx.txHash )

}




