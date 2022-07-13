package access_token

import (
	"DomainMan/pkg/api/handlers"
	"DomainMan/pkg/database"
	"DomainMan/pkg/models"
	"DomainMan/pkg/random"
	"fmt"
	"github.com/bndr/gotabulate"
	"github.com/spf13/cobra"
	"strconv"
	"time"
)

func RegisterCreate() *cobra.Command {
	createCommand := &cobra.Command{
		Use:     "create [NAME] [--master] [--can-issue]",
		Short:   "Create a new access token",
		Example: "create test --can-issue",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return database.ConnectUsingEnv(false)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return cmd.Help()
			}
			name := args[0]
			flags := cmd.Flags()
			master, _ := flags.GetBool("master")
			canIssue, _ := flags.GetBool("can-issue")
			accessToken := models.AccessToken{
				Name:     name,
				Token:    random.String(handlers.AccessTokenLength, random.DictAlphaNumber),
				IsMaster: master,
				CanIssue: canIssue,
			}
			db := database.DB
			if v := db.Create(&accessToken); v.Error != nil {
				return v.Error
			}
			t := gotabulate.Create([][]interface{}{{strconv.FormatUint(uint64(accessToken.ID), 10),
				accessToken.CreatedAt.Format(time.RFC3339), accessToken.UpdatedAt.Format(time.RFC3339),
				accessToken.Name, accessToken.Token, accessToken.IsMaster, accessToken.CanIssue,
			}})
			t.SetHeaders([]string{"ID", "Created At", "Updated At", "Name", "Token", "Is Master", "Can Issue"})
			t.SetAlign("right")
			fmt.Println(t.Render("grid"))
			return nil
		},
		DisableFlagsInUseLine: true,
	}
	flag := createCommand.Flags()
	flag.Bool("master", false, "set as master")
	flag.Bool("can-issue", false, "set as \"can issue\"")
	return createCommand
}
