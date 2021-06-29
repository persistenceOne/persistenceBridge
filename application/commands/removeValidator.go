package commands

import (
	"errors"
	"github.com/cosmos/cosmos-sdk/client"
	constants2 "github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/spf13/cobra"
	"log"
)

func RemoveCommand(initClientCtx client.Context) *cobra.Command {
	removeCommand := &cobra.Command{
		Use:   "remove [ValoperAddress]",
		Short: "Remove validator address to signing group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			homePath, err := cmd.Flags().GetString(constants2.FlagPBridgeHome)
			if err != nil {
				log.Fatalln(err)
			}

			database, err := db.OpenDB(homePath + "/db")
			validators, err := db.GetValidators()
			if err != nil {
				err2 := db.SetValidators([]string{})
				if err2 != nil {
					return err2
				}
			}
			defer database.Close()

			// TODO validate if validator is correct and already present
			//
			var newValidators []string
			for _, validator := range validators {
				if validator != args[0] {
					newValidators = append(newValidators, validator)
				}
			}

			// check that atleast one validator is present.
			if len(newValidators) == 0 {
				return errors.New("cannot remove all validators, need atleast one")
			}
			err = db.SetValidators(newValidators)
			if err != nil {
				return err
			}
			log.Printf("Updated set of validators: %v", newValidators)
			return nil
		},
	}

	removeCommand.Flags().String(constants2.FlagPBridgeHome, constants2.DefaultPBridgeHome, "home for pBridge")
	return removeCommand
}
