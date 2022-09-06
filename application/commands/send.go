package commands

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/persistenceOne/persistenceBridge/application/outgoingTx"
	"github.com/persistenceOne/persistenceBridge/tendermint"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
	"github.com/spf13/cobra"
	"log"
)

func SendCommand() *cobra.Command {
	sendCommand := &cobra.Command{
		Use:   "send [to_address] [amount]",
		Short: "send the amount to corresponding tendermint address",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			homePath, err := cmd.Flags().GetString(constants.FlagPBridgeHome)
			if err != nil {
				return err
			}

			timeout, err := cmd.Flags().GetString(constants.FlagTimeOut)
			if err != nil {
				log.Fatalln(err)
			}

			setAndSealConfig(homePath)

			toAddress, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			tmAddress := configuration.GetAppConfig().Tendermint.GetWrapAddress()
			logging.Info("Bridge (Wrap) Tendermint Address:", tmAddress)

			msg := &bankTypes.MsgSend{
				FromAddress: tmAddress,
				ToAddress:   toAddress.String(),
				Amount:      sdk.NewCoins(amount),
			}

			chain, err := tendermint.InitializeAndStartChain(timeout, homePath)
			if err != nil {
				return err
			}

			response, err := outgoingTx.LogMessagesAndBroadcast(chain, []sdk.Msg{msg}, 0)
			if err != nil {
				return err
			}
			logging.Info("MsgSend TxHash", response.TxHash)

			return nil
		},
	}

	sendCommand.Flags().String(constants.FlagTimeOut, constants.DefaultTimeout, "timeout time for connecting to rpc")
	sendCommand.Flags().String(constants.FlagPBridgeHome, constants.DefaultPBridgeHome, "home for pBridge")

	return sendCommand
}
