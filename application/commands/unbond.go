package commands

import (
	"github.com/Shopify/sarama"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/persistenceOne/persistenceBridge/application"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
	"github.com/spf13/cobra"
	"os"
)

func UnbondCommand() *cobra.Command {
	unbondCommand := &cobra.Command{
		Use:   "unbond [amount]",
		Short: "unbond amount from delegation",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			homePath, err := cmd.Flags().GetString(constants.FlagPBridgeHome)
			if err != nil {
				return err
			}

			setAndSealConfig(homePath)

			amount, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}

			tmAddress := configuration.GetAppConfig().Tendermint.GetWrapAddress()
			logging.Info("Bridge (Wrap) Tendermint Address:", tmAddress)

			kafkaConfig := utils.SaramaConfig()
			producer := utils.NewProducer(configuration.GetAppConfig().Kafka.Brokers, kafkaConfig)
			defer func(producer sarama.SyncProducer) {
				err := producer.Close()
				if err != nil {
					logging.Error("failed to close producer in topic: MsgUnbond via CLI, error:", err)
				}
			}(producer)

			encodingConfig := application.MakeEncodingConfig()
			clientContext := client.Context{}.
				WithCodec(encodingConfig.Marshaler).
				WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
				WithTxConfig(encodingConfig.TransactionConfig).
				WithLegacyAmino(encodingConfig.Amino).
				WithInput(os.Stdin).
				WithAccountRetriever(authTypes.AccountRetriever{}).
				WithHomeDir(homePath).
				WithViper("").
				WithBroadcastMode(configuration.GetAppConfig().Tendermint.BroadcastMode)

			protoCodec := codec.NewProtoCodec(clientContext.InterfaceRegistry)

			msg := &stakingTypes.MsgUndelegate{
				DelegatorAddress: tmAddress,
				ValidatorAddress: "",
				Amount:           amount,
			}

			msgBytes, err := protoCodec.MarshalInterface(msg)
			if err != nil {
				logging.Error("unable to marshal message-bytes")
				return err
			}

			err = utils.ProducerDeliverMessage(msgBytes, utils.EthUnbond, producer)
			if err != nil {
				logging.Error("failed to produce message from: EthUnbond to: MsgUnbond")
				return err
			}
			return nil
		},
	}
	unbondCommand.Flags().String(constants.FlagPBridgeHome, constants.DefaultPBridgeHome, "home for pBridge")
	return unbondCommand
}
