package ethereum

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/persistenceOne/persistenceBridge/application/casp"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	constants2 "github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/application/outgoingTx"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	test "github.com/persistenceOne/persistenceBridge/utilities/testing"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestOnNewBlock(t *testing.T) {
	pStakeConfig := configuration.InitConfig()
	configuration.SetConfig(test.GetCmdWithConfig())
	tmAddress, err := casp.GetTendermintAddress()
	require.Equal(t, nil, err)

	configuration.SetPStakeAddress(tmAddress)

	ethereumClient, err := ethclient.Dial(configuration.GetAppConfig().Ethereum.EthereumEndPoint)
	fmt.Println(configuration.GetAppConfig().Ethereum.EthereumEndPoint)
	require.Equal(t, nil, err)
	ctx := context.Background()
	kafkaProducer := utils.NewProducer(pStakeConfig.Kafka.Brokers, utils.SaramaConfig())
	latestEthHeight, err := ethereumClient.BlockNumber(ctx)
	txReceipt, err := ethereumClient.TransactionReceipt(ctx, common.HexToHash("0x034efa147e1ae645c5a6749fd6d19ec7f0c602a0ed22b122a13d76741b71764f"))

	fmt.Println(txReceipt.Status)

	database, err := db.OpenDB(constants2.TestDbDir)
	require.Nil(t, err)
	defer database.Close()

	TxhashFail := common.HexToHash("0x1fda0765d7803e4e9056a1bd849b4e17e92703710c278d6b5ce4d24e1eeca072")

	Address := common.BytesToAddress([]byte("0xA9739b5BdAfe956DEAa8b2e695c7d4f1DF7Bc1D6"))
	amt := new(big.Int)
	amt.SetInt64(1000)
	wrapTokenMsg := outgoingTx.WrapTokenMsg{
		Address: Address,
		Amount:  amt,
	}
	txd := []outgoingTx.WrapTokenMsg{wrapTokenMsg}

	ethTransaction := db.OutgoingEthereumTransaction{
		TxHash:   TxhashFail,
		Messages: txd,
	}

	err = db.SetOutgoingEthereumTx(ethTransaction)
	require.Equal(t, nil, err)
	err = onNewBlock(ctx, latestEthHeight, ethereumClient, &kafkaProducer)
	require.Equal(t, nil, err)



}
