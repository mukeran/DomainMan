package migrate

import (
	"DomainMan/pkg/cli/migrate/templates"
	"DomainMan/pkg/database"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"strings"
)

func Register() *cobra.Command {
	migrateRootCommand := &cobra.Command{
		Use:     "migrate [TEMPLATE]",
		Short:   "Migrate the database to the newest status",
		Example: "migrate initial - Using the initial template",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return cmd.Help()
			}
			err := database.ConnectUsingEnv(false)
			if err != nil {
				return err
			}
			db := database.DB
			tpl := strings.ToLower(args[0])
			switch tpl {
			case "initial":
				err = templates.Initial{}.Execute(db)
				if err != nil {
					log.Printf("Failed to execute template initial")
					return err
				}
			default:
				return fmt.Errorf("no such template %v", tpl)
			}
			log.Printf("Successfully executed template %v", tpl)
			return nil
		},
		DisableFlagsInUseLine: true,
	}
	return migrateRootCommand
}
