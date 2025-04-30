package cmd

import (
	"encoding/json"
	"os"

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
		},
		Action: func(c *cli.Context) error {
			project := types.Project{
				Name:        "my_project",
				Environment: c.String("env"),
				Variables: map[string]string{
					"version": "1.0",
				},
				Stages: []types.Stage{},
			}

			// Marshal the project struct to JSON
			data, err := json.MarshalIndent(project, "", "  ")
			if err != nil {
				panic(err)
			}

			// Write the JSON to a file
			err = os.WriteFile("lab_project.json", data, 0644)
			if err != nil {
				panic(err)
			}

			return nil
		},
	}
}
