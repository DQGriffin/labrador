package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"

	"github.com/DQGriffin/labrador/cmd"
)

var globalFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    "aws-access-key-id",
		Usage:   "AWS access key ID",
		EnvVars: []string{"AWS_ACCESS_KEY_ID"},
	},
	&cli.StringFlag{
		Name:    "aws-secret-access-key",
		Usage:   "AWS secret access key",
		EnvVars: []string{"AWS_SECRET_ACCESS_KEY"},
	},
	&cli.StringFlag{
		Name:    "aws-account-id",
		Usage:   "AWS account id",
		EnvVars: []string{"AWS_ACCOUNT_ID"},
	},
	&cli.StringFlag{
		Name:    "aws-region",
		Usage:   "AWS region",
		EnvVars: []string{"AWS_REGION"},
	},
	&cli.BoolFlag{
		Name:  "verbose",
		Usage: "Output extra information",
	},
	&cli.BoolFlag{
		Name:  "no-color",
		Usage: "Disable color in output",
	},
	&cli.BoolFlag{
		Name:  "debug",
		Usage: "Output debug information",
	},
}

var Version = "0.1.0"

func main() {
	err := godotenv.Load(".labrador.env")
	if err != nil {
		// Not an error. .labrador.env is optional
	}
	app := &cli.App{
		Name:    "labrador",
		Usage:   "Deploy and manage AWS resources",
		Version: Version,
		Flags:   globalFlags,
		Commands: []*cli.Command{
			cmd.DeployCommand(globalFlags),
			cmd.InitCommand(globalFlags),
			cmd.PlanCommand(globalFlags),
			cmd.DestroyCommand(globalFlags),
			cmd.InspectCommand(globalFlags),
			cmd.AddCommand(globalFlags),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
