package commands

import (
	"github.com/BurntSushi/toml"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	constants2 "github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/application/rpc"
	tendermint2 "github.com/persistenceOne/persistenceBridge/tendermint"
	"github.com/spf13/cobra"
	"log"
	"path/filepath"
)

func ShowCommand() *cobra.Command {
	showCommand := &cobra.Command{
		Use:   "show",
		Short: "show validators address to signing group",
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

			rpcEndpoint, err := cmd.Flags().GetString(constants2.FlagRPCEndpoint)
			if err != nil {
				log.Fatalln(err)
			}
			var validators []db.Validator
			database, err := db.OpenDB(homePath + "/db")
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

				validators, err = db.GetValidators()
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
	showCommand.Flags().String(constants2.FlagRPCEndpoint, constants2.DefaultRPCEndpoint, "rpc endpoint for bridge relayer")
	showCommand.Flags().String(constants2.FlagPBridgeHome, constants2.DefaultPBridgeHome, "home for pBridge")
	return showCommand
}
