package commands

import (
	"errors"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/persistenceOne/persistenceBridge/application"
	constants2 "github.com/persistenceOne/persistenceBridge/application/constants"
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

			db, err := application.OpenDB(homePath + "/db")
			validators, err := application.GetValidators()
			if err != nil {
				err2 := application.SetValidators([]string{})
				if err2 != nil {
					return err2
				}
			}
			defer db.Close()

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
			err = application.SetValidators(newValidators)
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
