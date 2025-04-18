package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/DQGriffin/labrador/cmd"
)

var globalFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    "env",
		Usage:   "Deployment environment",
		EnvVars: []string{"LABRADOR_ENV"},
	},
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
		Name:    "aws-region",
		Usage:   "AWS region",
		EnvVars: []string{"AWS_REGION"},
	},
	&cli.StringFlag{
		Name:    "project",
		Usage:   "Path to project file",
		EnvVars: []string{"PROJECT_PATH"},
	},
	&cli.BoolFlag{
		Name:  "only-create",
		Usage: "Only create new functions, skip updating existing ones",
	},
	&cli.BoolFlag{
		Name:  "only-update",
		Usage: "Only update existing functions, skip creating new ones",
	},
	&cli.StringFlag{
		Name:    "env-file",
		Usage:   "Path to env file",
		EnvVars: []string{"ENV_FILE"},
	},
}

func main() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }
	app := &cli.App{
		Name:  "labrador",
		Usage: "Deploy and manage AWS resources",
		Flags: globalFlags,
		Commands: []*cli.Command{
			cmd.DeployCommand(globalFlags),
			cmd.InitCommand(globalFlags),
			cmd.PlanCommand(globalFlags),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
