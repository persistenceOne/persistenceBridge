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
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	constants2 "github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/application/rpc"
	tendermint2 "github.com/persistenceOne/persistenceBridge/tendermint"
	"github.com/spf13/cobra"
)

func AddCommand() *cobra.Command {
	addCommand := &cobra.Command{
		Use:   "add [validatorOperatorAddress] [name]",
		Short: "Add validator address to signing group",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			homePath, err := cmd.Flags().GetString(constants2.FlagPBridgeHome)
			if err != nil {
				log.Fatalln(err)
			}

			pStakeConfig := configuration.InitConfig()
			_, err = toml.DecodeFile(filepath.Join(homePath, "config.toml"), &pStakeConfig)
			if err != nil {
				log.Fatalf("Error decoding pStakeConfig file: %v\n", err.Error())
			}
			_, err = tendermint2.SetBech32PrefixesAndPStakeWrapAddress()
			if err != nil {
				log.Fatalln(err)
			}
			configuration.ValidateAndSeal()
			validatorAddress, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			validatorName := args[1]

			rpcEndpoint, err := cmd.Flags().GetString(constants2.FlagRPCEndpoint)
			if err != nil {
				log.Fatalln(err)
			}
			var validators []db.Validator
			database, err := db.OpenDB(homePath + "/db")
			if err != nil {
				log.Printf("Db is already open: %v", err)
				log.Printf("sending rpc")
				var err2 error
				validators, err2 = rpc.AddValidator(db.Validator{
					Address: validatorAddress,
					Name:    validatorName,
				}, rpcEndpoint)
				if err2 != nil {
					return err2
				}
			} else {
				defer database.Close()

				err2 := db.SetValidator(db.Validator{
					Address: validatorAddress,
					Name:    validatorName,
				})
				if err2 != nil {
					return err2
				}

				validators, err2 = db.GetValidators()
				if err2 != nil {
					return err2
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

	addCommand.Flags().String(constants2.FlagRPCEndpoint, constants2.DefaultRPCEndpoint, "rpc endpoint for bridge relayer")
	addCommand.Flags().String(constants2.FlagPBridgeHome, constants2.DefaultPBridgeHome, "home for pBridge")
	return addCommand
}
