package cmd

import (
	"fmt"
	"os"
	"strings"

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
			&cli.StringFlag{
				Name:  "stage-types",
				Usage: "Restrict deployment to specific stage types",
			},
		},
		Before: func(c *cli.Context) error {
			if c.Bool("only-create") && c.Bool("only-update") {
				return fmt.Errorf("you can't use --only-create and --only-update at the same time")
			}

			if c.String("env-file") != "" {
				helpers.LoadEnvFile(c.String("env-file"))
			}
			utils.ReadCliArgs(c)

			if c.String("aws-account-id") != "" {
				fmt.Println("Using AWS account ID provided in flag")
				os.Setenv("AWS_ACCOUNT_ID", c.String("aws-account_id"))
				return nil
			}

			account, accErr := aws.GetAccountID()
			if accErr != nil {
				// Let's not stop execution here
				fmt.Println("Error", accErr.Error())
			}

			fmt.Println("Account ID", account)
			os.Setenv("AWS_ACCOUNT_ID", account)

			return nil
		},
		Action: func(c *cli.Context) error {
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

			stageTypesMap := make(map[string]bool)
			if c.String("stage-types") != "" {
				stageTypes := strings.Split(c.String("stage-types"), ",")

				for _, stageType := range stageTypes {
					stageTypesMap[stageType] = true
				}
			}

			commands.HandleDeployCommand(config, &stageTypesMap, existingLambdas, existingBuckets, onlyCreate, onlyUpdate)

			fmt.Println("Done")
			return nil
		},
	}
}
