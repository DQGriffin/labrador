package cmd

import (
	"fmt"
	"os"

	"github.com/DQGriffin/labrador/internal/helpers"
	"github.com/DQGriffin/labrador/internal/services/aws"
	"github.com/DQGriffin/labrador/pkg/utils"
	"github.com/urfave/cli/v2"
)

func PlanCommand(flags []cli.Flag) *cli.Command {
	return &cli.Command{
		Name:  "plan",
		Usage: "Preview actions labrador will take",
		Flags: flags,
		Action: func(c *cli.Context) error {
			fmt.Println("Planning...")

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

			utils.ReadCliArgs(c)

			existingLambdas, err := aws.ListLambdas()

			if err != nil {
				fmt.Println("An error occured while listing lambdas in the AWS account. ", err.Error())
				os.Exit(1)
			}

			var createCount = 0
			var updateCount = 0

			for _, functionGroup := range config.FunctionData {
				for _, function := range functionGroup.Functions {
					if _, exists := existingLambdas[function.Name]; exists {
						fmt.Println("Will be updated:", function.Name)
						updateCount += 1
					} else {
						fmt.Println("Will be created", function.Name)
						createCount += 1
					}
				}
			}

			fmt.Printf("Plan complete: %d to create, %d to update, %d to destroy\n", createCount, updateCount, 0)

			return nil
		},
	}
}
