package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

var Version = ""

func GetVersion() *cobra.Command {
	addCommand := &cobra.Command{
		Use:   "version",
		Short: "Print the application binary version information",
		RunE: func(cmd *cobra.Command, _ []string) error {
			fmt.Println(Version)
			return nil
		},
	}
	return addCommand
}
