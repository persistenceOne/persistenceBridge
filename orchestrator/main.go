/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"github.com/spf13/cobra"

	"github.com/persistenceOne/persistenceBridge/application/commands"
)

func main() {
	cobra.EnableCommandSorting = false

	rootCommand := &cobra.Command{
		Use:   "persistenceBridge",
		Short: "Persistence Bridge Orchestrator Daemon (server)",
	}

	rootCommand.AddCommand(commands.InitCommand())
	rootCommand.AddCommand(commands.AddCommand())
	rootCommand.AddCommand(commands.ShowCommand())
	rootCommand.AddCommand(commands.StartCommand())
	rootCommand.AddCommand(commands.RemoveCommand())

	if err := rootCommand.Execute(); err != nil {
		panic(err)
	}
}
