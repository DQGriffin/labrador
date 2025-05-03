package cmd

import (
	"github.com/DQGriffin/labrador/internal/cli/console"
	"github.com/DQGriffin/labrador/internal/helpers"
	"github.com/DQGriffin/labrador/internal/services/aws"
	"github.com/DQGriffin/labrador/pkg/utils"
	"github.com/urfave/cli/v2"
)

func PlanCommand(flags []cli.Flag) *cli.Command {
	return &cli.Command{
		Name:  "plan",
		Usage: "Preview actions labrador will take",
		Flags: []cli.Flag{
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
		},
		Before: func(c *cli.Context) error {
			console.SetColorEnabled(!c.Bool("no-color"))
			console.SetDebugOutputEnabled(c.Bool("debug"))
			console.SetVerboseOutputEnabled(c.Bool("verbose"))

			return nil
		},
		Action: func(c *cli.Context) error {
			console.Debug("planning...")

			var projectPath = "project.json"
			if c.String("project") != "" {
				projectPath = c.String("project")
			} else {
				console.Info("Project config file path not specified. Assuming project.json")
			}

			config, err := helpers.LoadProject(projectPath)

			if err != nil {
				console.Error("Error: Could not load project configuration")
				console.Fatal(err.Error())
			}

			utils.ReadCliArgs(c)

			existingLambdas, err := aws.ListLambdas()

			if err != nil {
				console.Fatal("An error occured while listing lambdas in the AWS account. ", err.Error())
			}

			var createCount = 0
			var updateCount = 0

			for _, functionGroup := range config.FunctionData {
				for _, function := range functionGroup.Functions {
					if _, exists := existingLambdas[function.Name]; exists {
						console.Info("Will be updated:", function.Name)
						updateCount += 1
					} else {
						console.Info("Will be created", function.Name)
						createCount += 1
					}
				}
			}

			console.Infof("Plan complete: %d to create, %d to update, %d to destroy\n", createCount, updateCount, 0)

			return nil
		},
	}
}
