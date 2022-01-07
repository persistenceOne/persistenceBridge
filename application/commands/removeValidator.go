/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package commands

import (
	"context"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/application/rpc"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"github.com/persistenceOne/persistenceBridge/tendermint"
)

func RemoveCommand() *cobra.Command {
	removeCommand := &cobra.Command{
		Use:   "remove [validatorOperatorAddress]",
		Short: "Remove validator address to signing group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			homePath, err := cmd.Flags().GetString(constants.FlagPBridgeHome)
			if err != nil {
				log.Fatalln(err)
			}

			pStakeConfig := configuration.GetAppConfig()
			_, err = toml.DecodeFile(filepath.Join(homePath, "config.toml"), &pStakeConfig)
			if err != nil {
				log.Fatalf("Error decoding pStakeConfig file: %v\n", err.Error())
			}

			// fixme: init with deps and timeout
			ctx := context.Background()

			_, err = tendermint.SetBech32PrefixesAndPStakeWrapAddress(ctx)
			if err != nil {
				log.Fatalln(err)
			}

			configuration.ValidateAndSeal()

			validatorAddress, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			rpcEndpoint, err := cmd.Flags().GetString(constants.FlagRPCEndpoint)
			if err != nil {
				log.Fatalln(err)
			}

			kafkaPorts, err := cmd.Flags().GetString(constants.FlagKafkaPorts)
			if err != nil {
				log.Fatalln(err)
			}

			config := utils.SaramaConfig()
			producer := utils.NewProducer(strings.Split(kafkaPorts, ","), config)

			defer func() {
				innerErr := producer.Close()
				if innerErr != nil {
					log.Printf("failed to close producer in topic: %v\n", utils.MsgUnbond)
				}
			}()

			var validators []db.Validator

			database, err := db.OpenDB(filepath.Join(homePath, "db"))
			if err != nil {
				log.Printf("Db is already open: %v", err)
				log.Printf("sending rpc")

				validators, err = rpc.RemoveValidator(validatorAddress, rpcEndpoint)
				if err != nil {
					return err
				}
			} else {
				defer database.Close()

				err = db.DeleteValidator(database, validatorAddress)
				if err != nil {
					return err
				}

				validators, err = db.GetValidators(database)
				if err != nil {
					return err
				}
			}

			if len(validators) == 0 {
				log.Println("IMPORTANT: No validator present to redelegate!!!")

				return ErrNoValidators
			}

			log.Printf("Total validators %d:\n", len(validators))

			for i, validator := range validators {
				log.Printf("%d. %s - %s\n", i+1, validator.Name, validator.Address.String())
			}

			time.Sleep(1 * time.Minute)

			err = utils.ProducerDeliverMessage(validatorAddress, utils.Redelegate, producer)
			if err != nil {
				log.Printf("failed to produce message to topic %v\n", utils.Redelegate)
				return err
			}

			return nil
		},
	}

	removeCommand.Flags().String(constants.FlagRPCEndpoint, constants.DefaultRPCEndpoint, "rpc endpoint for bridge relayer")
	removeCommand.Flags().String(constants.FlagPBridgeHome, constants.DefaultPBridgeHome(), "home for pBridge")
	removeCommand.Flags().String(constants.FlagKafkaPorts, constants.DefaultKafkaPorts, "broker ports kafka is running on")

	return removeCommand
}
