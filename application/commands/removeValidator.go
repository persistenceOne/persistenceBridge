package commands

import (
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
			validatorAddress, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			homePath, err := cmd.Flags().GetString(constants2.FlagPBridgeHome)
			if err != nil {
				log.Fatalln(err)
			}

			database, err := db.OpenDB(homePath + "/db")
			defer database.Close()

			err = db.DeleteValidator(validatorAddress)
			if err != nil {
				return err
			}

			validators, err := db.GetValidators()
			if err != nil {
				return err
			}

			log.Printf("Updated set of validators: %v\n", validators)
			if len(validators) == 0 {
				log.Println("IMPORTANT: No validator present!!!")
			}
			return nil
		},
	}

	removeCommand.Flags().String(constants2.FlagPBridgeHome, constants2.DefaultPBridgeHome, "home for pBridge")
	return removeCommand
}
