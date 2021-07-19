package ethereum

import (
	"context"

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
				err := handleTransaction(client, ctx, transaction, contract, kafkaProducer, protoCodec)
				if err != nil {
					logging.Error("Failed to process ethereum tx:", transaction.Hash().String())
					return err
				}
			}
		}
	}
	return nil
}

func handleTransaction(client *ethclient.Client, ctx *context.Context, transaction *types.Transaction, contract contracts2.ContractI, kafkaProducer *sarama.SyncProducer, protoCodec *codec.ProtoCodec) error {
	receipt, err := client.TransactionReceipt(*ctx, transaction.Hash())
	if err != nil {
		logging.Error("Unable to get receipt of tx:", transaction.Hash().String(), "Error:", err)
		return err
	}

	if receipt.Status == 1 {
		logging.Info("Received Ethereum Tx:", transaction.Hash().String())
		method, arguments, err := contract.GetMethodAndArguments(transaction.Data())
		if err != nil {
			logging.Fatal("Error in getting method and arguments of:", contract.GetName(), "Error:", err)
		}

		if processFunc, ok := contract.GetMethods()[method.RawName]; ok {
			err = processFunc(kafkaProducer, protoCodec, arguments)
			if err != nil {
				logging.Fatal("Error in processing arguments of contract:", contract.GetName(), "method:", method.RawName, "for TX:", transaction.Hash().String(), "Error", err)
			}
		}
	}
	return nil
}
