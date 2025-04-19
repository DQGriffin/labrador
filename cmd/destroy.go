package cmd

import (
	"fmt"
	"os"

	"github.com/DQGriffin/labrador/internal/commands"
	"github.com/DQGriffin/labrador/internal/helpers"
	"github.com/DQGriffin/labrador/pkg/utils"
	"github.com/urfave/cli/v2"
)

func DestroyCommand(flags []cli.Flag) *cli.Command {
	return &cli.Command{
		Name:  "destroy",
		Usage: "Destroy resources",
		Flags: flags,
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

			commandErr := commands.HandleDestroyCommand(config, isDryRun, env)
			return commandErr
		},
	}
}
