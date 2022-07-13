package cli

import (
	"DomainMan/pkg/cli/manage"
	"DomainMan/pkg/cli/migrate"
	"DomainMan/pkg/cli/server"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

const (
	Version = "1.0"
)

var (
	rootCommand            *cobra.Command
	dbDialect, dbParameter string
)

func register() {
	// Root
	cobra.OnInitialize(initCommand)
	rootCommand = &cobra.Command{
		Use: "dm-cli [OPTIONS] COMMAND [ARG...]",
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
	rootCommand.PersistentFlags().StringVar(&dbDialect, "db-dialect", "sqlite", "Specify database dialect (mysql, sqlite)")
	rootCommand.PersistentFlags().StringVar(&dbParameter, "db-parameter", "db.sqlite3", "Specify database parameter (e.g.: db.sqlite3, mysql://root:root@127.0.0.1:3306/DomainMan)")
	// Migrate
	rootCommand.AddCommand(migrate.Register())
	// Server
	rootCommand.AddCommand(server.Register())
	// Manage
	rootCommand.AddCommand(manage.Register())
}

func initCommand() {
	if os.Getenv("DOMAINMAN_DATABASE_DIALECT") == "" {
		os.Setenv("DOMAINMAN_DATABASE_DIALECT", dbDialect)
	}
	if os.Getenv("DOMAINMAN_DATABASE_PARAMETER") == "" {
		os.Setenv("DOMAINMAN_DATABASE_PARAMETER", dbParameter)
	}
}

func Run() {
	register()
	if err := rootCommand.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
