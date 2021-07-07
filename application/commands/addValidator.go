package commands

import (
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	constants2 "github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/spf13/cobra"
	"log"
)

func AddCommand(initClientCtx client.Context) *cobra.Command {
	addCommand := &cobra.Command{
		Use:   "add [ValoperAddress]",
		Short: "Add validator address to signing group",
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

			err = db.SetValidator(db.Validator{
				Address: validatorAddress,
				Active:  true,
			})
			if err != nil {
				return err
			}

			validators, err := db.GetValidators()
			if err != nil {
				return err
			}
			log.Printf("Updated set of validators: %v\n", validators)
			return nil
		},
	}
	addCommand.Flags().String(constants2.FlagPBridgeHome, constants2.DefaultPBridgeHome, "home for pBridge")
	return addCommand
}
