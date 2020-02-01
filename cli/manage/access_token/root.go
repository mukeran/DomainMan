package access_token

import "github.com/spf13/cobra"

func Register() *cobra.Command {
	accessTokenRootCommand := &cobra.Command{
		Use:   "access_token [list/create/delete]",
		Short: "manage access token",
		Example: "access_token list\n" +
			"access_token create [NAME]\n" +
			"access_token delete [ACCESS_TOKEN_ID]",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		DisableFlagsInUseLine: true,
	}
	// create
	accessTokenRootCommand.AddCommand(RegisterCreate())
	// delete
	accessTokenRootCommand.AddCommand(RegisterDelete())
	// list
	accessTokenRootCommand.AddCommand(RegisterList())
	return accessTokenRootCommand
}
