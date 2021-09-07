package ethereum

import (
	"context"
	"fmt"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"

	"github.com/Shopify/sarama"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	contracts2 "github.com/persistenceOne/persistenceBridge/ethereum/contracts"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
)

func handleBlock(client *ethclient.Client, ctx *context.Context, block *types.Block, kafkaProducer *sarama.SyncProducer, protoCodec *codec.ProtoCodec) error {
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
				err := collectEthTx(client, ctx, protoCodec, transaction, contract)
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

func collectEthTx(client *ethclient.Client, ctx *context.Context, protoCodec *codec.ProtoCodec, transaction *types.Transaction, contract contracts2.ContractI) error {
	receipt, err := client.TransactionReceipt(*ctx, transaction.Hash())
	if err != nil {
		logging.Error("Unable to get receipt of tx:", transaction.Hash().String(), "Error:", err)
		return err
	}

	if receipt.Status == 1 {
		logging.Info("Received Ethereum Tx:", transaction.Hash().String())
		method, arguments, err := contract.GetMethodAndArguments(transaction.Data())
		if err != nil {
			return fmt.Errorf("unable to get method and arguments of: %s Error: %s", contract.GetName(), err.Error())
		}

		if processFunc, ok := contract.GetSDKMsgAndSender()[method.RawName]; ok {
			msg, sender, err := processFunc(arguments)
			if err != nil {
				return fmt.Errorf("failed to process arguments of contract: %s method: %s for TX: %s Error: %s", contract.GetName(), method.RawName, transaction.Hash().String(), err.Error())
			}
			exists := db.CheckEthereumIncomingTxProduced(transaction.Hash())
			if !exists {
				msgBytes, err := protoCodec.MarshalInterface(msg)
				if err != nil {
					return err
				}
				err = db.AddToPendingEthereumIncomingTx(db.EthereumIncomingTx{
					TxHash:          transaction.Hash(),
					ProducedToKafka: false,
					MsgBytes:        msgBytes,
					Sender:          sender,
					MsgType:         msg.Type(),
				})
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func produceToKafka(kafkaProducer *sarama.SyncProducer) {
	ethInTxs, err := db.GetProduceToKafkaEthereumIncomingTxs()
	if err != nil {
		logging.Fatal(err)
	}
	for _, ethTxToTM := range ethInTxs {
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
		logging.Info("Adding to [ETH Listener] kafka producer:", producer, "of txHash:", ethTxToTM.TxHash.String(), "msgType:", ethTxToTM.MsgType, "sender:", ethTxToTM.Sender.String())
		err = utils.ProducerDeliverMessage(ethTxToTM.MsgBytes, producer, *kafkaProducer)
		if err != nil {
			logging.Fatal("Failed to add msg to kafka queue [ETH Listener], producer:", producer, "txHash:", ethTxToTM.TxHash.String(), "sender:", ethTxToTM.Sender.String(), "error:", err)
		}
		err := db.SetEthereumIncomingTxProduced(ethTxToTM)
		if err != nil {
			logging.Fatal(err)
		}
	}
}
