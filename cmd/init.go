package cmd

import (
	"encoding/json"
	"os"

	"github.com/DQGriffin/labrador/internal/cli/console"
	"github.com/DQGriffin/labrador/pkg/types"
	"github.com/urfave/cli/v2"
)

func InitCommand(flags []cli.Flag) *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "Generate a sample Labrador project config",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "env",
				Usage:   "Environment to initialize project with (ex. dev, staging, prod)",
				EnvVars: []string{"ENV"},
			},
			&cli.StringFlag{
				Name:    "name",
				Aliases: []string{"n"},
				Usage:   "Name of the project",
				EnvVars: []string{"NAME"},
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "Path to output project config to",
				EnvVars: []string{"OUTPUT"},
			},
		},
		Action: func(c *cli.Context) error {

			projectName := "my_project"
			if c.String("name") != "" {
				projectName = c.String("name")
			}

			projectEnv := "dev"
			if c.String("env") != "" {
				projectEnv = c.String("env")
			}

			outputPath := "labrador_project.json"
			if c.String("output") != "" {
				outputPath = c.String("output")
			}

			project := types.Project{
				Name:        projectName,
				Environment: projectEnv,
				Variables: map[string]string{
					"version": "1.0",
				},
				Stages: []types.Stage{},
			}

			// Marshal the project struct to JSON
			data, err := json.MarshalIndent(project, "", "\t")
			if err != nil {
				panic(err)
			}

			// Write the JSON to a file
			err = os.WriteFile(outputPath, data, 0644)
			if err != nil {
				panic(err)
			}

			console.Infof("Wrote project config to %s", outputPath)
			if c.String("name") == "" && c.String("output") == "" && c.String("env") == "" {
				console.Info("Tip: use --name, --env, or --output next time to customize")
			}

			return nil
		},
	}
}
