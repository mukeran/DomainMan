package access_token

import (
	"DomainMan/pkg/database"
	"DomainMan/pkg/models"
	"fmt"
	"github.com/bndr/gotabulate"
	"github.com/spf13/cobra"
)

func RegisterDelete() *cobra.Command {
	deleteCommand := &cobra.Command{
		Use:     "delete [ACCESS_TOKEN_ID...]",
		Short:   "Delete access tokens",
		Example: "delete 1 2 3",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return database.ConnectUsingEnv(false)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}
			var result [][]interface{}
			for _, id := range args {
				if v := database.DB.Delete(models.AccessToken{}, "id = ?", id); v.Error != nil {
					result = append(result, []interface{}{id, "Failed", v.Error})
				} else {
					result = append(result, []interface{}{id, "Successful", "-"})
				}
			}
			headers := []string{"ID", "Result", "Error Message"}
			t := gotabulate.Create(result)
			t.SetHeaders(headers)
			t.SetAlign("right")
			fmt.Println(t.Render("grid"))
			return nil
		},
		DisableFlagsInUseLine: true,
	}
	return deleteCommand
}
