package cmdline

import (
	"github.com/spf13/cobra"

	"github.com/a-x-a/gophkeeper/internal/keeperctl/app"
	"github.com/a-x-a/gophkeeper/internal/keeperctl/entity"
)

func login(cmd *cobra.Command, _ []string) error {
	clientApp, err := app.FromContext(cmd.Context())
	if err != nil {
		return err
	}

	token, err := clientApp.Usecases.Auth.Login(
		cmd.Context(),
		cfg.Username,
		clientApp.Key,
	)
	if err != nil {
		clientApp.Log.Debug().Err(err).Msg("")

		return entity.Unwrap(err)
	}

	clientApp.AccessToken = token
	clientApp.Log.Debug().
		Str("access-token", token).
		Msg("Login successful")

	return nil
}
