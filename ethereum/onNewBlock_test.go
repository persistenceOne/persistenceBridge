/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package ethereum

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/persistenceOne/persistenceBridge/application/casp"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	constants2 "github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	test "github.com/persistenceOne/persistenceBridge/utilities/testing"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestOnNewBlock(t *testing.T) {
	test.SetTestConfig()
	tmAddress, err := casp.GetTendermintAddress()
	require.Equal(t, nil, err)
	ethAddress, err := casp.GetEthAddress()
	require.Equal(t, nil, err)
	configuration.SetCASPAddresses(tmAddress, ethAddress)

	ethereumClient, err := ethclient.Dial(configuration.GetAppConfig().Ethereum.EthereumEndPoint)
	require.Equal(t, nil, err)
	ctx := context.Background()
	kafkaProducer := utils.NewProducer(configuration.GetAppConfig().Kafka.GetBrokersList(), utils.SaramaConfig())
	latestEthHeight, err := ethereumClient.BlockNumber(ctx)

	database, err := db.OpenDB(constants2.TestDbDir)
	require.Nil(t, err)
	defer database.Close()

	TxhashFail := common.HexToHash("0xb96560e8ef15a0d86f8156daddf6f2421d962f5a37dd8e2ba212b05eddaf0b59")

	Address := common.BytesToAddress([]byte("0xce3f57a8de9aa69da3289871b5fee5e77ffcf480"))
	amt := new(big.Int)
	amt.SetInt64(1000)
	wrapTokenMsg := db.WrapTokenMsg{
		Address:       Address,
		StakingAmount: amt,
	}
	txd := []db.WrapTokenMsg{wrapTokenMsg}

	ethTransaction := db.OutgoingEthereumTransaction{
		TxHash:   TxhashFail,
		Messages: txd,
	}

	err = db.SetOutgoingEthereumTx(ethTransaction)
	require.Equal(t, nil, err)
	err = onNewBlock(ctx, latestEthHeight, ethereumClient, &kafkaProducer, nil)
	require.Equal(t, nil, err)

	TxhashSuccess := common.HexToHash("0x8e08d80c37c884467b9b48a77e658711615a5cfde43f95fccfb3b95ee66cd6ea")
	Address = common.BytesToAddress([]byte("0x477573f212a7bdd5f7c12889bd1ad0aa44fb82aa"))
	amt = new(big.Int)
	amt.SetInt64(1000)
	wrapTokenMsg = db.WrapTokenMsg{
		Address:       Address,
		StakingAmount: amt,
	}
	txd = []db.WrapTokenMsg{wrapTokenMsg}

	ethTransaction = db.OutgoingEthereumTransaction{
		TxHash:   TxhashSuccess,
		Messages: txd,
	}

	err = db.SetOutgoingEthereumTx(ethTransaction)
	require.Equal(t, nil, err)
	err = onNewBlock(ctx, latestEthHeight, ethereumClient, &kafkaProducer, nil)
	require.Equal(t, nil, err)

}
