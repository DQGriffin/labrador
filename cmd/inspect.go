package cmd

import (
	"fmt"
	"os"

	"github.com/DQGriffin/labrador/internal/commands"
	"github.com/DQGriffin/labrador/internal/helpers"
	"github.com/DQGriffin/labrador/pkg/utils"
	"github.com/urfave/cli/v2"
)

func InspectCommand(flags []cli.Flag) *cli.Command {
	return &cli.Command{
		Name:  "inspect",
		Usage: "Inspect project configurations",
		Flags: flags,
		Action: func(c *cli.Context) error {
			fmt.Println("Inspecting...")

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

			config, err := helpers.LoadProject(projectPath)

			if err != nil {
				fmt.Println("Error: Could not load project configuration")
				fmt.Println(err.Error())
				os.Exit(1)
			}

			verbose := c.Bool("verbose")

			commands.HandleInspectCommand(&config, "plain", verbose)

			return nil
		},
	}
}
