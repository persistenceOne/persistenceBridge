/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package commands

import (
	"log"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/application/rpc"
)

func ShowCommand() *cobra.Command {
	showCommand := &cobra.Command{
		Use:   "show",
		Short: "show wrap address, eth bridge admin and validators",
		RunE: func(cmd *cobra.Command, args []string) error {
			homePath, err := cmd.Flags().GetString(constants.FlagPBridgeHome)
			if err != nil {
				log.Fatalln(err)
			}

			setAndSealConfig(homePath)

			rpcEndpoint, err := cmd.Flags().GetString(constants.FlagRPCEndpoint)
			if err != nil {
				log.Fatalln(err)
			}

			var validators []db.Validator

			database, err := db.OpenDB(filepath.Join(homePath, "db"))
			if err != nil {
				log.Printf("Db is already open: %v", err)
				log.Printf("sending rpc to %v", rpcEndpoint)

				var err2 error
				validators, err2 = rpc.ShowValidators("", rpcEndpoint)
				if err2 != nil {
					return err2
				}
			} else {
				defer database.Close()

				validators, err = db.GetValidators(database)
				if err != nil {
					return err
				}
			}

			if len(validators) == 0 {
				log.Fatalln("No validators in db, panic.")
			} else {
				log.Printf("Total validators %d:\n", len(validators))

				for i, validator := range validators {
					log.Printf("%d. %s - %s\n", i+1, validator.Name, validator.Address.String())
				}
			}

			return nil
		},
	}

	showCommand.Flags().String(constants.FlagRPCEndpoint, constants.DefaultRPCEndpoint, "rpc endpoint for bridge relayer")
	showCommand.Flags().String(constants.FlagPBridgeHome, constants.DefaultPBridgeHome(), "home for pBridge")

	return showCommand
}
