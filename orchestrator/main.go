/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/debug"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	serverCmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/cosmos/cosmos-sdk/version"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/persistenceOne/persistenceBridge/application"
	"github.com/persistenceOne/persistenceBridge/application/commands"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/spf13/cobra"
	tendermintClient "github.com/tendermint/tendermint/libs/cli"

	"os"
)

const flagInvalidCheckPeriod = "invalid-check-period"

var invalidCheckPeriod uint

func main() {
	encodingConfig := application.MakeEncodingConfig()
	initClientCtx := client.Context{}.
		WithJSONMarshaler(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TransactionConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(authTypes.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastBlock).
		WithHomeDir(constants.DefaultPBridgeHome)

	cobra.EnableCommandSorting = false

	rootCommand := &cobra.Command{
		Use:   "persistenceBridge",
		Short: "Persistence Bridge Orchestrator Daemon (server)",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}

			return server.InterceptConfigsPreRunHandler(cmd)
		},
	}

	rootCommand.AddCommand(tendermintClient.NewCompletionCmd(rootCommand, true))
	rootCommand.AddCommand(debug.Cmd())
	rootCommand.AddCommand(version.NewVersionCommand())
	rootCommand.PersistentFlags().UintVar(
		&invalidCheckPeriod,
		flagInvalidCheckPeriod,
		0,
		"Assert registered invariants every N blocks",
	)
	rootCommand.AddCommand(commands.StartCommand(initClientCtx))
	rootCommand.AddCommand(commands.AddCommand(initClientCtx))
	rootCommand.AddCommand(commands.RemoveCommand(initClientCtx))
	rootCommand.AddCommand(commands.ShowCommand(initClientCtx))
	rootCommand.AddCommand(commands.InitCommand())

	if err := serverCmd.Execute(rootCommand, constants.DefaultPBridgeHome); err != nil {
		switch e := err.(type) {
		case server.ErrorCode:
			os.Exit(e.Code)

		default:
			os.Exit(1)
		}
	}
}
