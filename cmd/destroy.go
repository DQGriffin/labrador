package cmd

import (
	"strings"

	"github.com/DQGriffin/labrador/internal/cli/console"
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
			&cli.BoolFlag{
				Name:    "force",
				Usage:   "Take extra steps to force the deletion of resources",
				EnvVars: []string{"DRY_RUN"},
			},
			&cli.StringFlag{
				Name:    "stage-types",
				Usage:   "Restrict destroy operations for stage types in a comma-separated list",
				EnvVars: []string{"STAGE_TYPES"},
			},
		},
		Before: func(c *cli.Context) error {
			console.SetColorEnabled(!c.Bool("no-color"))
			console.SetDebugOutputEnabled(c.Bool("debug"))

			return nil
		},
		Action: func(c *cli.Context) error {
			console.Info("Destroy")

			if c.String("env-file") != "" {
				helpers.LoadEnvFile(c.String("env-file"))
			}
			utils.ReadCliArgs(c)

			var projectPath = "project.json"
			if c.String("project") != "" {
				projectPath = c.String("project")
			} else {
				console.Info("Project config file path not specified. Assuming project.json")
			}

			var isDryRun = c.Bool("dry-run")

			config, err := helpers.LoadProject(projectPath)

			if err != nil {
				console.Error("Could not load project configuration")
				console.Fatal(err.Error())
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

			force := c.Bool("force")

			commandErr := commands.HandleDestroyCommand(config, isDryRun, force, &stageTypesMap, env)
			return commandErr
		},
	}
}
