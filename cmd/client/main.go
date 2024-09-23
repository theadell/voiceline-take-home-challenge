package main

import (
	"log"
	"os"

	"adelhub.com/voiceline/internal/config"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
)

func main() {

	cfg := config.LoadConfig()

	config := oauth2.Config{
		ClientID:     cfg.OAuth2DeviceClientId,
		ClientSecret: cfg.OAuth2DeviceCientIdExt,
		Scopes:       []string{"profile", "email", "openid"},
		Endpoint:     endpoints.Google,
	}

	app := &cli.App{
		Name:  "auth-client",
		Usage: "Helper CLI for user login and registration",
		Commands: []*cli.Command{
			RegisterCommand(cfg),
			LoginCommand(cfg, config),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}
