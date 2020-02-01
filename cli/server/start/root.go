package start

import (
	"DomainMan/pkg/mq"
	"DomainMan/pkg/server"
	"github.com/spf13/cobra"
)

func Register() *cobra.Command {
	startRootCommand := &cobra.Command{
		Use:     "start",
		Short:   "Start server",
		Example: "server start --port 80",
		RunE: func(cmd *cobra.Command, args []string) error {
			mq.Init()
			server.Init()
			return server.Run()
		},
		DisableFlagsInUseLine: true,
	}
	return startRootCommand
}
