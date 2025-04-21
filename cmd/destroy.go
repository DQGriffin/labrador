package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/DQGriffin/labrador/internal/commands"
	"github.com/DQGriffin/labrador/internal/helpers"
	"github.com/DQGriffin/labrador/pkg/utils"
	"github.com/urfave/cli/v2"
)

func DestroyCommand(flags []cli.Flag) *cli.Command {
	return &cli.Command{
		Name:  "destroy",
		Usage: "Destroy resources",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "dry-run",
				Usage:   "Preview actions labarador will take without making changes to infrastructure",
				EnvVars: []string{"DRY_RUN"},
			},
			&cli.StringFlag{
				Name:    "stage-types",
				Usage:   "Restrict destroy operations for stage types in a comma-separated list",
				EnvVars: []string{"STAGE_TYPES"},
			},
		},
		Action: func(c *cli.Context) error {
			fmt.Println("Destroy")

			if c.String("env-file") != "" {
				helpers.LoadEnvFile(c.String("env-file"))
			}
			utils.ReadCliArgs(c)

			var projectPath = "project.json"
			if c.String("project") != "" {
				projectPath = c.String("project")
			} else {
				fmt.Println("Project config file path not specified. Assuming project.json")
			}

			var isDryRun = c.Bool("dry-run")

			config, err := helpers.LoadProject(projectPath)

			if err != nil {
				fmt.Println("Error: Could not load project configuration")
				fmt.Println(err.Error())
				os.Exit(1)
			}

			var env = config.Project.Environment
			if c.String("env") != "" {
				env = c.String("env")
			}

			stageTypesMap := make(map[string]bool)
			if c.String("stage-types") != "" {
				stageTypes := strings.Split(c.String("stage-types"), ",")

				for _, stageType := range stageTypes {
					stageTypesMap[stageType] = true
				}
			}

			fmt.Println(stageTypesMap)

			commandErr := commands.HandleDestroyCommand(config, isDryRun, &stageTypesMap, env)
			return commandErr
		},
	}
}
