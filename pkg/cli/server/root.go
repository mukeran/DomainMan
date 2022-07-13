package server

import (
	"DomainMan/pkg/cli/server/start"
	"github.com/spf13/cobra"
)

func Register() *cobra.Command {
	serverRootCommand := &cobra.Command{
		Use:   "server [COMMAND]",
		Short: "Start or manage server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		DisableFlagsInUseLine: true,
	}
	// Start
	serverRootCommand.AddCommand(start.Register())
	return serverRootCommand
}
