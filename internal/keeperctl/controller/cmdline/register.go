package cmdline

import (
	"github.com/spf13/cobra"

	"github.com/a-x-a/gophkeeper/internal/keeperctl/app"
	"github.com/a-x-a/gophkeeper/internal/keeperctl/entity"
)

var registerCmd = &cobra.Command{
	Use:   "register [flags]",
	Short: "Register a new user",
	RunE:  doRegister,
}

func init() {
	rootCmd.AddCommand(registerCmd)
}

func doRegister(cmd *cobra.Command, args []string) error {
	clientApp, err := app.FromContext(cmd.Context())
	if err != nil {
		return err
	}

	accessToken, err := clientApp.Usecases.Users.Register(
		cmd.Context(),
		cfg.Username,
		clientApp.Key,
	)
	if err != nil {
		clientApp.Log.Debug().Err(err).Msg("")

		return entity.Unwrap(err)
	}

	clientApp.Log.Debug().Str("access-token", accessToken).Msg("New user successfully created")

	return nil
}
