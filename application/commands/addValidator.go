/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package commands

import (
	"log"
	"path/filepath"

	"github.com/BurntSushi/toml"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dgraph-io/badger/v3"
	"github.com/spf13/cobra"

	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/application/rpc"
	"github.com/persistenceOne/persistenceBridge/tendermint"
)

func AddCommand() *cobra.Command {
	const argsCount = 2

	addCommand := &cobra.Command{
		Use:   "add [validatorOperatorAddress] [name]",
		Short: "Add validator address to signing group",
		Args:  cobra.ExactArgs(argsCount),
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

			_, err = tendermint.SetBech32PrefixesAndPStakeWrapAddress()
			if err != nil {
				log.Fatalln(err)
			}

			configuration.ValidateAndSeal()

			var validatorAddress sdk.ValAddress
			validatorAddress, err = sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			validatorName := args[1]

			var rpcEndpoint string
			rpcEndpoint, err = cmd.Flags().GetString(constants.FlagRPCEndpoint)
			if err != nil {
				log.Fatalln(err)
			}

			var validators []db.Validator
			var database *badger.DB

			database, err = db.OpenDB(filepath.Join(homePath, "db"))
			if err != nil {
				log.Printf("Db is already open: %v", err)
				log.Printf("sending rpc")

				validators, err = rpc.AddValidator(db.Validator{
					Address: validatorAddress,
					Name:    validatorName,
				}, rpcEndpoint)
				if err != nil {
					return err
				}
			} else {
				defer func() {
					err = database.Close()
					log.Printf("DB got an error while closing: %v", err)
				}()

				err = db.SetValidator(database, db.Validator{
					Address: validatorAddress,
					Name:    validatorName,
				})
				if err != nil {
					return err
				}

				validators, err = db.GetValidators(database)
				if err != nil {
					return err
				}
			}

			if len(validators) == 0 {
				log.Println("No validators in db, panic.")
			} else {
				log.Printf("Total validators %d:\n", len(validators))

				for i, validator := range validators {
					log.Printf("%d. %s - %s\n", i+1, validator.Name, validator.Address.String())
				}
			}

			return nil
		},
	}

	addCommand.Flags().String(constants.FlagRPCEndpoint, constants.DefaultRPCEndpoint, "rpc endpoint for bridge relayer")
	addCommand.Flags().String(constants.FlagPBridgeHome, constants.DefaultPBridgeHome, "home for pBridge")

	return addCommand
}
