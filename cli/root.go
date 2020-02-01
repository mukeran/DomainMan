package cli

import (
	"DomainMan/cli/manage"
	"DomainMan/cli/migrate"
	"DomainMan/cli/server"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

const (
	Version = "1.0"
)

var (
	rootCommand *cobra.Command
)

func register() {
	// Root
	rootCommand = &cobra.Command{
		Use: "DomainMan [OPTIONS] COMMAND [ARG...]",
		Short: "DomainMan is a domain management platform which implements the domain DNS management, WHOIS query and" +
			" Domain/Subdomain watch.",
		SilenceUsage:     true,
		SilenceErrors:    true,
		TraverseChildren: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		Version:               Version,
		DisableFlagsInUseLine: true,
	}
	// Migrate
	rootCommand.AddCommand(migrate.Register())
	// Server
	rootCommand.AddCommand(server.Register())
	// Manage
	rootCommand.AddCommand(manage.Register())
}

func Run() {
	register()
	if err := rootCommand.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
