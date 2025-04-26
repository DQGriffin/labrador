package cmd

import (
	"fmt"
	"os"
	"strings"

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
			&cli.BoolFlag{
				Name:  "verbose",
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

			if c.String("aws-account-id") != "" {
				fmt.Println("Using AWS account ID provided in flag")
				os.Setenv("AWS_ACCOUNT_ID", c.String("aws-account_id"))
				return nil
			}

			account, accErr := aws.GetAccountID()
			if accErr != nil {
				// Let's not stop execution here
				fmt.Println("Error", accErr.Error())
			}

			fmt.Println("Account ID", account)
			os.Setenv("AWS_ACCOUNT_ID", account)

			return nil
		},
		Action: func(c *cli.Context) error {
			fmt.Println("Inspecting...")

			var projectPath = "project.json"
			if c.String("project") != "" {
				projectPath = c.String("project")
			} else {
				fmt.Println("Project config file path not specified. Assuming project.json")
			}

			config, err := helpers.LoadProject(projectPath)

			if err != nil {
				fmt.Println("Error: Could not load project configuration")
				fmt.Println(err.Error())
				os.Exit(1)
			}

			verbose := c.Bool("verbose")

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
