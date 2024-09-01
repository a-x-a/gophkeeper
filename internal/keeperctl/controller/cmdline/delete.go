package cmdline

import (
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/cobra"

	"github.com/a-x-a/gophkeeper/internal/keeperctl/app"
	"github.com/a-x-a/gophkeeper/internal/keeperctl/entity"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [secret id] [flags]",
	Short: "Delete the secret",
	Args:  cobra.MinimumNArgs(1),
	RunE:  doDelete,
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}

func doDelete(cmd *cobra.Command, args []string) error {
	id, err := uuid.FromString(args[0])
	if err != nil {
		return err
	}

	clientApp, err := app.FromContext(cmd.Context())
	if err != nil {
		return err
	}

	if err := clientApp.Usecases.Secrets.Delete(cmd.Context(), clientApp.AccessToken, id); err != nil {
		clientApp.Log.Debug().Err(err).Msg("")

		return entity.Unwrap(err)
	}

	return nil
}
