package main

import (
	"fmt"

	"adelhub.com/voiceline/internal/config"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"
)

func RegisterCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "register",
		Usage: "Register a new user with email and password",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "email",
				Usage:    "User's email address for registration",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "password",
				Usage:    "User's password (must be at least 8 characters)",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			email := c.String("email")
			password := c.String("password")
			if err := RegisterUser(cfg, email, password); err != nil {
				return fmt.Errorf("failed to register: %v", err)
			}
			return nil
		},
	}
}

func LoginCommand(cfg *config.Config, oauth2Config oauth2.Config) *cli.Command {
	return &cli.Command{
		Name:  "login",
		Usage: "Login a user with email and password or SSO",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "email",
				Usage: "User's email address",
			},
			&cli.StringFlag{
				Name:  "password",
				Usage: "User's password",
			},
			&cli.BoolFlag{
				Name:  "sso",
				Usage: "Login via SSO",
			},
		},
		Action: func(c *cli.Context) error {
			email := c.String("email")
			password := c.String("password")
			sso := c.Bool("sso")

			if sso {
				if err := LoginWithSSO(cfg, oauth2Config); err != nil {
					return fmt.Errorf("failed to login via SSO: %v", err)
				}
			} else if email != "" && password != "" {
				if err := LoginWithEmail(cfg, email, password); err != nil {
					return fmt.Errorf("failed to login with email: %v", err)
				}
			} else {
				return fmt.Errorf("please provide either --email and --password or --sso")
			}

			return nil
		},
	}
}
