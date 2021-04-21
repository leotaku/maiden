package cmd

import (
	"fmt"
	"net/url"
	"path"

	"github.com/adrg/xdg"
	"github.com/leotaku/maiden/oauth"
	"github.com/spf13/cobra"
)

var (
	dataHome     = path.Join(xdg.DataHome, "maiden")
	googleScopes = []string{"https://www.googleapis.com/auth/calendar"}
	googleURL, _ = url.Parse("https://apidata.googleusercontent.com/")
)

const googleCredentials = `{
  "installed": {
    "client_id": "298650252192-5t47ioau6lj0a1tb04vddlfarflpjbml.apps.googleusercontent.com",
    "project_id": "handmaiden-1611867768099",
    "auth_uri": "https://accounts.google.com/o/oauth2/auth",
    "token_uri": "https://oauth2.googleapis.com/token",
    "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
    "client_secret": "t3li8TvUuDOTvaH_eAjg1W25",
    "redirect_uris": ["urn:ietf:wg:oauth:2.0:oob", "http://localhost"]
  }
}`

var initCmd = &cobra.Command{
	Use:     "init [flags..]",
	Short:   "Initialize credentials and secrets",
	Version: rootCmd.Version,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		life := oauth.NewLifecycle(dataHome)
		_, _, err := life.Init(googleCredentials, googleScopes...)
		if err != nil {
			return fmt.Errorf("auth: %w", err)
		}

		return nil
	},
	DisableFlagsInUseLine: true,
}
