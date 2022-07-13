package access_token

import (
	"DomainMan/pkg/database"
	"DomainMan/pkg/models"
	"fmt"
	"github.com/bndr/gotabulate"
	"github.com/spf13/cobra"
	"strconv"
	"time"
)

func RegisterList() *cobra.Command {
	listCommand := &cobra.Command{
		Use:     "list [--name [QUERY_NAME]] [--master/non-master] [--can-issue/cannot-issue] [--show-token] [--offset [OFFSET]] [--limit [LIMIT]]",
		Short:   "List available access tokens",
		Example: "list --show-token",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return database.ConnectUsingEnv(false)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			db := database.DB
			flags := cmd.Flags()
			id, _ := flags.GetUint("id")
			name, _ := flags.GetString("name")
			offset, _ := flags.GetInt("offset")
			limit, _ := flags.GetInt("limit")
			master, _ := flags.GetBool("master")
			nonMaster, _ := flags.GetBool("non-master")
			canIssue, _ := flags.GetBool("can-issue")
			cannotIssue, _ := flags.GetBool("cannot-issue")
			showToken, _ := flags.GetBool("show-token")
			query := db.Table(models.TableAccessToken).Offset(offset).Limit(limit)
			if id != 0 {
				query = query.Where("id = ?", id)
			}
			if name != "" {
				query = query.Where("name like ?", name)
			}
			if master {
				query = query.Where("is_master = true")
			}
			if nonMaster {
				query = query.Where("is_master = false")
			}
			if canIssue {
				query = query.Where("can_issue = true")
			}
			if cannotIssue {
				query = query.Where("can_issue = false")
			}
			var headers []string
			if !showToken {
				query = query.Select([]string{"id", "created_at", "updated_at", "name", "is_master", "can_issue", "issuer_id"})
				headers = []string{"ID", "Created At", "Updated At", "Name", "Is Master", "Can Issue", "Issuer ID"}
			} else {
				query = query.Select([]string{"id", "created_at", "updated_at", "name", "token", "is_master", "can_issue", "issuer_id"})
				headers = []string{"ID", "Created At", "Updated At", "Name", "Token", "Is Master", "Can Issue", "Issuer ID"}
			}
			var accessTokens []models.AccessToken
			if v := query.Find(&accessTokens); v.Error != nil {
				return v.Error
			}
			if len(accessTokens) == 0 {
				fmt.Println("No entries")
			} else {
				fmt.Printf("Find %v entries\n", len(accessTokens))
				var out [][]interface{}
				for _, accessToken := range accessTokens {
					if !showToken {
						out = append(out, []interface{}{strconv.FormatUint(uint64(accessToken.ID), 10),
							accessToken.CreatedAt.Format(time.RFC3339), accessToken.UpdatedAt.Format(time.RFC3339),
							accessToken.Name, accessToken.IsMaster, accessToken.CanIssue,
							strconv.FormatUint(uint64(accessToken.IssuerID), 10),
						})
					} else {
						out = append(out, []interface{}{strconv.FormatUint(uint64(accessToken.ID), 10),
							accessToken.CreatedAt.Format(time.RFC3339), accessToken.UpdatedAt.Format(time.RFC3339),
							accessToken.Name, accessToken.Token, accessToken.IsMaster, accessToken.CanIssue,
							strconv.FormatUint(uint64(accessToken.IssuerID), 10),
						})
					}
				}
				t := gotabulate.Create(out)
				t.SetHeaders(headers)
				t.SetAlign("right")
				fmt.Println(t.Render("grid"))
			}
			return nil
		},
		DisableFlagsInUseLine: true,
	}
	flags := listCommand.Flags()
	flags.Uint("id", 0, "Specify id")
	flags.String("name", "", "Query name")
	flags.IntP("offset", "O", 0, "Page control - offset")
	flags.IntP("limit", "L", 10, "Page control - limit")
	flags.Bool("master", false, "Show master access tokens")
	flags.Bool("non-master", false, "Show no master access tokens")
	flags.Bool("can-issue", false, "Show \"can-issue\" access tokens")
	flags.Bool("cannot-issue", false, "Show has no \"can-issue\" access tokens")
	flags.Bool("show-token", false, "Show hidden tokens")
	return listCommand
}
