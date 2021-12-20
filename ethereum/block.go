/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package ethereum

import (
	"context"
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/persistenceOne/persistenceBridge/application/db"
	contracts2 "github.com/persistenceOne/persistenceBridge/ethereum/contracts"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
)

func handleBlock(ctx context.Context, client *ethclient.Client, block *types.Block, kafkaProducer sarama.SyncProducer, protoCodec *codec.ProtoCodec) error {
	for _, transaction := range block.Transactions() {
		if transaction.To() != nil {
			var contract contracts2.ContractI

			switch transaction.To().String() {
			case contracts2.LiquidStaking.GetAddress().String():
				contract = &contracts2.LiquidStaking
			case contracts2.TokenWrapper.GetAddress().String():
				contract = &contracts2.TokenWrapper
			default:
			}

			if contract != nil {
				err := collectEthTx(ctx, client, protoCodec, transaction, contract)
				if err != nil {
					logging.Error("Failed to process ethereum tx:", transaction.Hash().String())

					return err
				}
			}
		}
	}

	produceToKafka(kafkaProducer)

	return nil
}

func collectEthTx(ctx context.Context, client *ethclient.Client, protoCodec *codec.ProtoCodec, transaction *types.Transaction, contract contracts2.ContractI) error {
	receipt, err := client.TransactionReceipt(ctx, transaction.Hash())
	if err != nil {
		logging.Error("Unable to get receipt of tx:", transaction.Hash().String(), "Error:", err)

		return err
	}

	if receipt.Status == 1 {
		logging.Info("Received Ethereum Tx:", transaction.Hash().String())

		var (
			method    *abi.Method
			arguments []interface{}
		)

		method, arguments, err = contract.GetMethodAndArguments(transaction.Data())
		if err != nil {
			return fmt.Errorf("%w of: %s Error: %s", ErrUnableGetMethodAndArgument, contract.GetName(), err.Error())
		}

		if processFunc, ok := contract.GetSDKMsgAndSender()[method.RawName]; ok {
			var (
				msg    sdk.Msg
				sender common.Address
			)

			msg, sender, err = processFunc(arguments)
			if err != nil {
				return fmt.Errorf("%w: %s method: %s for TX: %s Error: %s",
					ErrFailedProcessArguments, contract.GetName(), method.RawName, transaction.Hash().String(), err.Error())
			}

			// Do not check for EthereumTxToKafka exists.
			if !db.CheckIncomingEthereumTxExists(transaction.Hash()) {
				var msgBytes []byte

				msgBytes, err = protoCodec.MarshalInterface(msg)
				if err != nil {
					return err
				}

				err = db.AddIncomingEthereumTx(&db.IncomingEthereumTx{
					TxHash:   transaction.Hash(),
					MsgBytes: msgBytes,
					Sender:   sender,
					MsgType:  msg.Type(),
				})
				if err != nil {
					return err
				}

				err = db.AddEthereumTxToKafka(&db.EthereumTxToKafka{
					TxHash: transaction.Hash(),
				})

				if err != nil {
					return fmt.Errorf("%w: %v", ErrCannotAddTxToKafka, err)
				}
			}
		}
	}

	return nil
}

func produceToKafka(kafkaProducer sarama.SyncProducer) {
	ethTxToKafkaList, err := db.GetAllEthereumTxToKafka()
	if err != nil {
		logging.Fatal(err)
	}

	for _, tx := range ethTxToKafkaList {
		ethTxToTM, err := db.GetIncomingEthereumTx(tx.TxHash)
		if err != nil {
			logging.Fatal(err)
		}

		producer := ""

		switch ethTxToTM.MsgType {
		case bankTypes.TypeMsgSend:
			producer = utils.MsgSend
		case stakingTypes.TypeMsgDelegate:
			producer = utils.MsgDelegate
		case stakingTypes.TypeMsgUndelegate:
			producer = utils.EthUnbond
		default:
			logging.Fatal("unknown msg type [ETH Listener]")
		}

		logging.Info("[ETH Listener] Adding to kafka producer:", producer, "of txHash:", ethTxToTM.TxHash.String(), "msgType:", ethTxToTM.MsgType, "sender:", ethTxToTM.Sender.String())

		err = utils.ProducerDeliverMessage(ethTxToTM.MsgBytes, producer, kafkaProducer)
		if err != nil {
			logging.Fatal("Failed to add msg to kafka queue [ETH Listener], producer:", producer, "txHash:", ethTxToTM.TxHash.String(), "sender:", ethTxToTM.Sender.String(), "error:", err)
		}

		err = db.DeleteEthereumTxToKafka(ethTxToTM.TxHash)
		if err != nil {
			logging.Fatal(err)
		}
	}
}
