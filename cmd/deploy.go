package cmd

import (
	"fmt"
	"os"

	"github.com/DQGriffin/labrador/internal/commands"
	"github.com/DQGriffin/labrador/internal/helpers"
	"github.com/DQGriffin/labrador/internal/services/aws"
	"github.com/DQGriffin/labrador/pkg/utils"
	"github.com/urfave/cli/v2"
)

func DeployCommand(flags []cli.Flag) *cli.Command {
	return &cli.Command{
		Name:  "deploy",
		Usage: "Deploy Lambda functions defined in your config",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "only-create",
				Usage: "Only create new resources, skip updating existing ones",
			},
			&cli.BoolFlag{
				Name:  "only-update",
				Usage: "Only update resources, skip creating new ones",
			},
		},
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

			existingLambdas, err := aws.ListLambdas()

			if err != nil {
				fmt.Println("Error: Could not list lambdas in AWS account. Check permissions ", err.Error())
				os.Exit(1)
			}

			ctx, cfg, err := aws.GetConfig("us-east-1")

			if err != nil {
				return err
			}

			client := aws.GetClient(cfg)

			existingBuckets, bucketErr := aws.ListBuckets(ctx, client)

			if bucketErr != nil {
				fmt.Println("Error: Could not list buckets in AWS account. Check permissions ", bucketErr.Error())
				os.Exit(1)
			}

			onlyCreate := c.Bool("only-create")
			onlyUpdate := c.Bool("only-update")

			commands.HandleDeployCommand(config, existingLambdas, existingBuckets, onlyCreate, onlyUpdate)

			fmt.Println("Done")
			return nil
		},
	}
}
