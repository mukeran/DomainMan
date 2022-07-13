package start

import (
	"DomainMan/pkg/api"
	"DomainMan/pkg/database"
	"DomainMan/pkg/mq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func Register() *cobra.Command {
	startRootCommand := &cobra.Command{
		Use:     "start",
		Short:   "Start server",
		Example: "server start --listen 0.0.0.0:8899",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := database.ConnectUsingEnv(false)
			if err != nil {
				logrus.Errorf("Failed to connect database")
				return err
			}
			mq.Init()
			api.Init()
			listen, err := cmd.Flags().GetString("listen")
			if err != nil {
				return err
			}
			return api.Run(listen)
		},
		DisableFlagsInUseLine: true,
	}
	startRootCommand.Flags().StringP("listen", "l", "0.0.0.0:8899", "Server's listening address")
	return startRootCommand
}
