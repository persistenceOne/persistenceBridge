/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package ethereum

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"

	"github.com/persistenceOne/persistenceBridge/application/casp"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/application/outgoingtx"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"github.com/persistenceOne/persistenceBridge/utilities/test"
)

func TestOnNewBlock(t *testing.T) {
	pStakeConfig := configuration.InitConfig()
	configuration.SetConfig(test.GetCmdWithConfig())

	tmAddress, err := casp.GetTendermintAddress()
	require.Nil(t, err)

	configuration.SetPStakeAddress(tmAddress)

	ethereumClient, err := ethclient.Dial(configuration.GetAppConfig().Ethereum.EthereumEndPoint)
	require.Nil(t, err)

	kafkaProducer := utils.NewProducer(pStakeConfig.Kafka.Brokers, utils.SaramaConfig())

	ctx := context.Background()
	latestEthHeight, err := ethereumClient.BlockNumber(ctx)
	require.Nil(t, err)

	database, err := db.OpenDB(constants.TestDBDir)
	require.Nil(t, err)

	defer database.Close()

	txhashFail := common.HexToHash("0xb96560e8ef15a0d86f8156daddf6f2421d962f5a37dd8e2ba212b05eddaf0b59")

	address := common.BytesToAddress([]byte("0xce3f57a8de9aa69da3289871b5fee5e77ffcf480"))
	amt := new(big.Int)
	amt.SetInt64(1000)

	wrapTokenMsg := outgoingtx.WrapTokenMsg{
		Address: address,
		Amount:  amt,
	}

	txd := []outgoingtx.WrapTokenMsg{wrapTokenMsg}

	ethTransaction := db.OutgoingEthereumTransaction{
		TxHash:   txhashFail,
		Messages: txd,
	}

	err = db.SetOutgoingEthereumTx(ethTransaction)
	require.Nil(t, err)

	err = onNewBlock(ctx, latestEthHeight, ethereumClient, kafkaProducer)
	require.Nil(t, err)

	TxhashSuccess := common.HexToHash("0x8e08d80c37c884467b9b48a77e658711615a5cfde43f95fccfb3b95ee66cd6ea")
	address = common.BytesToAddress([]byte("0x477573f212a7bdd5f7c12889bd1ad0aa44fb82aa"))
	amt = new(big.Int)
	amt.SetInt64(1000)

	wrapTokenMsg = outgoingtx.WrapTokenMsg{
		Address: address,
		Amount:  amt,
	}

	txd = []outgoingtx.WrapTokenMsg{wrapTokenMsg}

	ethTransaction = db.OutgoingEthereumTransaction{
		TxHash:   TxhashSuccess,
		Messages: txd,
	}

	err = db.SetOutgoingEthereumTx(ethTransaction)
	require.Nil(t, err)

	err = onNewBlock(ctx, latestEthHeight, ethereumClient, kafkaProducer)
	require.Nil(t, err)
}
