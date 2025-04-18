package cmd

import (
	"fmt"
	"os"

	"github.com/DQGriffin/labrador/internal/aws/lambda"
	"github.com/DQGriffin/labrador/internal/helpers"
	"github.com/DQGriffin/labrador/pkg/utils"
	"github.com/urfave/cli/v2"
)

func DeployCommand(flags []cli.Flag) *cli.Command {
	return &cli.Command{
		Name:  "deploy",
		Usage: "Deploy Lambda functions defined in your config",
		Flags: flags,
		Before: func(c *cli.Context) error {
			if c.Bool("only-create") && c.Bool("only-update") {
				return fmt.Errorf("you can't use --only-create and --only-update at the same time")
			}
			return nil
		},
		Action: func(c *cli.Context) error {
			fmt.Println("Deploying...")

			var projectPath = "project.json"
			if c.String("project") != "" {
				projectPath = c.String("project")
			} else {
				fmt.Println("Project config file path not specified. Assuming project.json")
			}

			if c.String("env-file") != "" {
				helpers.LoadEnvFile(c.String("env-file"))
			}

			config, err := helpers.LoadProject(projectPath)

			if err != nil {
				fmt.Println("Error: Could not load project configuration")
				fmt.Println(err.Error())
				os.Exit(1)
			}

			utils.ReadCliArgs(c)

			existingLambdas, err := lambda.ListLambdas()

			if err != nil {
				fmt.Println("Error: Could not list lambdas in AWS account. Check permissions ", err.Error())
				os.Exit(1)
			}

			for _, stage := range config.Project.Stages {
				fmt.Printf("Deploying Stage: %s\n", stage.Name)
				fmt.Printf("Type: %s\n", stage.Type)

				for _, fnConfig := range stage.Functions {
					for _, fn := range fnConfig.Functions {
						if _, exists := existingLambdas[fn.Name]; exists {
							if !c.Bool("only-create") {
								fmt.Println("updating function", fn.Name)
								// lambda.UpdateLambda(fn)
							}
						} else {
							if !c.Bool("only-update") {
								fmt.Println("creating function", fn.Name)
								// lambda.CreateLambda(fn)
							}
						}
					}
				}
			}

			fmt.Println("Done")
			return nil
		},
	}
}
