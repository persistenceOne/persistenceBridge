package commands

import (
	"errors"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	constants2 "github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/application/rpc"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"github.com/spf13/cobra"
	"log"
	"strings"
	"time"
)

func RemoveCommand(initClientCtx client.Context) *cobra.Command {
	removeCommand := &cobra.Command{
		Use:   "remove [ValoperAddress]",
		Short: "Remove validator address to signing group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			validatorAddress, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			homePath, err := cmd.Flags().GetString(constants2.FlagPBridgeHome)
			if err != nil {
				log.Fatalln(err)
			}
			rpcEndpoint, err := cmd.Flags().GetString(constants2.FlagRPCEndpoint)
			if err != nil {
				log.Fatalln(err)
			}
			var validators []sdk.ValAddress
			database, err := db.OpenDB(homePath + "/db")
			if err != nil {
				log.Printf("Db is already open: %v", err)
				log.Printf("sending rpc")
				var err2 error
				validators, err2 = rpc.RemoveValidator(validatorAddress, rpcEndpoint)
				if err2 != nil {
					return err2
				}
			} else {
				defer database.Close()
				err = db.DeleteValidator(validatorAddress)
				if err != nil {
					return err
				}

				var err2 error
				validators, err2 = db.GetValidators()
				if err2 != nil {
					return err2
				}

			}
			log.Printf("Updated set of validators: %v\n", validators)
			if len(validators) == 0 {
				log.Println("IMPORTANT: No validator present to redelegate!!!")
				return errors.New("Needs to have atleast one validator to redelegate to.")
			}

			time.Sleep(1 * time.Minute)
			kafkaPorts, err := cmd.Flags().GetString(constants2.FlagKafkaPorts)
			if err != nil {
				log.Fatalln(err)
			}
			config := utils.SaramaConfig()
			producer := utils.NewProducer(strings.Split(kafkaPorts, ","), config)
			defer func() {
				err := producer.Close()
				if err != nil {
					log.Printf("failed to close producer in topic: %v\n", utils.MsgUnbond)
				}
			}()

			err = utils.ProducerDeliverMessage(validatorAddress, utils.Redelegate, producer)
			if err != nil {
				log.Printf("failed to produce message to topic %v\n", utils.Redelegate)
				return err
			}

			return nil
		},
	}
	removeCommand.Flags().String(constants2.FlagRPCEndpoint, constants2.DefaultRPCEndpoint, "rpc endpoint for bridge relayer")
	removeCommand.Flags().String(constants2.FlagPBridgeHome, constants2.DefaultPBridgeHome, "home for pBridge")
	removeCommand.Flags().String(constants2.FlagKafkaPorts, constants2.DefaultKafkaPorts, "broker ports kafka is running on")
	return removeCommand
}
