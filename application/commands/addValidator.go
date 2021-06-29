package commands

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/persistenceOne/persistenceBridge/application"
	constants2 "github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/spf13/cobra"
	"log"
)

func AddCommand(initClientCtx client.Context) *cobra.Command {
	addCommand := &cobra.Command{
		Use:   "add [ValoperAddress]",
		Short: "Add validator address to signing group",
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

			// TODO validate if validator is correct and not already present
			//
			validators = append(validators, args[0])
			err = application.SetValidators(validators)
			if err != nil {
				return err
			}

			log.Printf("Updated set of validators: %v", validators)
			return nil
		},
	}
	addCommand.Flags().String(constants2.FlagPBridgeHome, constants2.DefaultPBridgeHome, "home for pBridge")
	return addCommand
}
