package cli

import (
	"github.com/spf13/cobra"
)

func (cli *CLI) buildVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Get version of NewCommander CLI",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			showSuccess(cli.version)
		},
	}

	return cmd
}
