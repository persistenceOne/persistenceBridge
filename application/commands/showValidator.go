package commands

import (
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	constants2 "github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/application/rpc"
	"github.com/spf13/cobra"
	"log"
)

func ShowCommand(initClientCtx client.Context) *cobra.Command {
	showCommand := &cobra.Command{
		Use:   "show",
		Short: "show validator address to signing group",
		RunE: func(cmd *cobra.Command, args []string) error {
			homePath, err := cmd.Flags().GetString(constants2.FlagPBridgeHome)
			if err != nil {
				log.Fatalln(err)
			}
			rpcEndpoint, err := cmd.Flags().GetString(constants2.FlagRPCEndpoint)
			if err != nil {
				log.Fatalln(err)
			}
			var validators []sdk.ValAddress
			database, err := db.OpenDB(homePath + "/db")
			if err != nil {
				log.Printf("Db is already open: %v", err)
				log.Printf("sending rpc")
				var err2 error
				validators, err2 = rpc.ShowValidators("", rpcEndpoint)
				if err2 != nil {
					return err2
				}
			} else {
				defer database.Close()

				validators, err = db.GetValidators()
				if err != nil {
					return err
				}

			}
			log.Printf("List set of validators: %v\n", validators)
			if len(validators) == 0 {
				log.Println("No validators in db, panic.")
			}
			return nil
		},
	}
	showCommand.Flags().String(constants2.FlagRPCEndpoint, constants2.DefaultRPCEndpoint, "rpc endpoint for bridge relayer")
	showCommand.Flags().String(constants2.FlagPBridgeHome, constants2.DefaultPBridgeHome, "home for pBridge")
	return showCommand
}
