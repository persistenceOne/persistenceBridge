package commands

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
)

func RemoveCommand(initClientCtx client.Context) *cobra.Command {
	removeCommand := &cobra.Command{
		Use:   "remove [ValoperAddress]",
		Short: "Remove validator address to signing group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	return removeCommand
}
