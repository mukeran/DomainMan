package manage

import (
	"DomainMan/cli/manage/access_token"
	"github.com/spf13/cobra"
)

func Register() *cobra.Command {
	manageRootCommand := &cobra.Command{
		Use:   "manage [access_token/config/domain/registrar/suffix/whois]",
		Short: "Manage the resources in database",
		Example: "manage access_token create [NAME]\n" +
			"manage config set [KEY] [VALUE]\n" +
			"manage domain add [NAME]\n" +
			"manage whois delete [WHOIS_ID]",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		DisableFlagsInUseLine: true,
	}
	// Access Token
	manageRootCommand.AddCommand(access_token.Register())
	return manageRootCommand
}
