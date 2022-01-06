/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package commands

import (
	"log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/application/rpc"
	"github.com/spf13/cobra"
)

func AddCommand() *cobra.Command {
	addCommand := &cobra.Command{
		Use:   "add [validatorOperatorAddress] [name]",
		Short: "Add validator address to signing group",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			homePath, err := cmd.Flags().GetString(constants.FlagPBridgeHome)
			if err != nil {
				log.Fatalln(err)
			}

			setAndSealConfig(homePath)

			validatorAddress, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			validatorName := args[1]

			rpcEndpoint, err := cmd.Flags().GetString(constants.FlagRPCEndpoint)
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
				log.Fatalln("No validators found in db.")
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
