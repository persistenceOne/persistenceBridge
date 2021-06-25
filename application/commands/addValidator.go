package commands

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
)

func AddCommand(initClientCtx client.Context) *cobra.Command {
	addCommand := &cobra.Command{
		Use:   "add [ValoperAddress]",
		Short: "Add validator address to signing group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			return nil
		},
	}
	return addCommand
}
