package cmd

import (
	"github.com/DQGriffin/labrador/internal/cli/console"
	"github.com/DQGriffin/labrador/internal/commands/add"
	"github.com/DQGriffin/labrador/internal/helpers"
	"github.com/DQGriffin/labrador/pkg/utils"
	"github.com/urfave/cli/v2"
)

func AddCommand(flags []cli.Flag) *cli.Command {
	return &cli.Command{
		Name:  "add",
		Usage: "Add a resource",
		Before: func(c *cli.Context) error {
			console.SetColorEnabled(!c.Bool("no-color"))
			console.SetDebugOutputEnabled(c.Bool("debug"))
			return nil
		},
		Subcommands: []*cli.Command{
			{
				Name:  "stage",
				Usage: "Add a new stage to the project",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "type",
						Aliases: []string{"t"},
						Usage:   "Type of the stage (ex. lambda, s3)",
					},
					&cli.StringFlag{
						Name:    "name",
						Aliases: []string{"n"},
						Usage:   "Name of the stage",
					},
					&cli.StringFlag{
						Name:    "project",
						Aliases: []string{"p"},
						Usage:   "Path to the project file",
					},
					&cli.StringFlag{
						Name:    "output",
						Aliases: []string{"o"},
						Usage:   "Path to output templates",
					},
				},
				Before: func(c *cli.Context) error {
					if c.String("type") == "" {
						console.Fatal("You must provide a stage type via --type")
					}

					if c.String("name") == "" {
						console.Fatal("You must provide a stage name via --name")
					}

					if c.String("project") == "" {
						console.Fatal("You must provide a project via --project")
					}

					if c.String("output") == "" {
						console.Fatal("You must provide an output path via --output")
					}

					if c.String("env-file") != "" {
						helpers.LoadEnvFile(c.String("env-file"))
					}
					utils.ReadCliArgs(c)

					return nil
				},
				Action: func(c *cli.Context) error {
					console.Debug("Add stage")
					projectPath := c.String("project")
					stageType := c.String("type")
					stageName := c.String("name")
					outputPath := c.String("output")

					err := add.HandleAddStage(projectPath, stageType, stageName, outputPath)

					if err != nil {
						console.Fatal(err.Error())
					}

					return nil
				},
			},
		},
	}
}
