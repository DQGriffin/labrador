package cmd

import (
	"os"
	"strings"

	"github.com/DQGriffin/labrador/internal/cli/console"
	"github.com/DQGriffin/labrador/internal/commands"
	"github.com/DQGriffin/labrador/internal/helpers"
	"github.com/DQGriffin/labrador/internal/services/aws"
	"github.com/DQGriffin/labrador/pkg/utils"
	"github.com/urfave/cli/v2"
)

func InspectCommand(flags []cli.Flag) *cli.Command {
	return &cli.Command{
		Name:  "inspect",
		Usage: "Inspect project configurations",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "env",
				Usage:   "Deployment environment",
				EnvVars: []string{"LABRADOR_ENV"},
			},
			&cli.StringFlag{
				Name:    "project",
				Usage:   "Path to project file",
				EnvVars: []string{"PROJECT_PATH"},
			},
			&cli.StringFlag{
				Name:    "env-file",
				Usage:   "Path to env file",
				EnvVars: []string{"ENV_FILE"},
			},
			&cli.BoolFlag{
				Name:  "full",
				Usage: "Output all information for resources",
			},
			&cli.StringFlag{
				Name:  "output",
				Usage: "Set output mode (plain, tree)",
			},
			&cli.StringFlag{
				Name:  "stage-types",
				Usage: "Comma-separated list of stage types to include (e.g., 'lambda,api')",
			},
		},
		Before: func(c *cli.Context) error {
			if c.String("env-file") != "" {
				helpers.LoadEnvFile(c.String("env-file"))
			}
			utils.ReadCliArgs(c)

			console.SetColorEnabled(!c.Bool("no-color"))
			console.SetDebugOutputEnabled(c.Bool("debug"))
			console.SetVerboseOutputEnabled(c.Bool("verbose"))

			if c.String("aws-account-id") != "" {
				console.Debug("Using AWS account ID provided in flag")
				os.Setenv("AWS_ACCOUNT_ID", c.String("aws-account_id"))
				return nil
			}

			account, accErr := aws.GetAccountID()
			if accErr != nil {
				// Let's not stop execution here
				console.Error(accErr.Error())
			}

			console.Debug("Account ID ", account)
			os.Setenv("AWS_ACCOUNT_ID", account)

			return nil
		},
		Action: func(c *cli.Context) error {
			console.Debug("Inspecting...")

			var projectPath = "project.json"
			if c.String("project") != "" {
				projectPath = c.String("project")
			} else {
				console.Info("Project config file path not specified. Assuming project.json")
			}

			config, err := helpers.LoadProject(projectPath)

			if err != nil {
				console.Error("Could not load project configuration")
				console.Fatal(err.Error())
			}

			verbose := c.Bool("full")

			stageTypesMap := make(map[string]bool)
			if c.String("stage-types") != "" {
				stageTypes := strings.Split(c.String("stage-types"), ",")

				for _, stageType := range stageTypes {
					stageTypesMap[stageType] = true
				}
			}

			var outputMode = "tree"
			if c.String("output") != "" {
				outputMode = c.String("output")
			}

			commands.HandleInspectCommand(&config, outputMode, &stageTypesMap, verbose)

			return nil
		},
	}
}
