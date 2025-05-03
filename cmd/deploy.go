package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/DQGriffin/labrador/internal/cli/console"
	"github.com/DQGriffin/labrador/internal/commands"
	"github.com/DQGriffin/labrador/internal/constants"
	"github.com/DQGriffin/labrador/internal/helpers"
	"github.com/DQGriffin/labrador/internal/services/aws"
	"github.com/DQGriffin/labrador/pkg/utils"
	lambdaTypes "github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/urfave/cli/v2"
)

func DeployCommand(flags []cli.Flag) *cli.Command {
	return &cli.Command{
		Name:  "deploy",
		Usage: "Deploy Lambda functions defined in your config",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "env",
				Usage:   "Deployment environment",
				EnvVars: []string{"LABRADOR_ENV"},
			},
			&cli.StringFlag{
				Name:    "project",
				Usage:   "Path to project file",
				EnvVars: []string{"PROJECT_PATH"},
			},
			&cli.BoolFlag{
				Name:  "dry-run",
				Usage: "Preview operations before taking action",
			},
			&cli.StringFlag{
				Name:    "env-file",
				Usage:   "Path to env file",
				EnvVars: []string{"ENV_FILE"},
			},
			&cli.StringFlag{
				Name:    "stages",
				Usage:   "Comma-separated list of stage types to deploy (e.g. lambda,s3)",
				EnvVars: []string{"DEPLOY_STAGES"},
			},
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
			&cli.IntFlag{
				Name:  "wait-time",
				Usage: "Wait time when waiting for changes to propagate",
			},
			&cli.BoolFlag{
				Name:  "no-wait",
				Usage: "Disable waiting for changes to propagate",
			},
		},
		Before: func(c *cli.Context) error {
			console.SetColorEnabled(!c.Bool("no-color"))
			console.SetDebugOutputEnabled(c.Bool("debug"))
			console.SetVerboseOutputEnabled(c.Bool("verbose"))

			if c.Int("wait-time") < 0 {
				console.Fatal("wait time cannot be negative")
			}

			if c.Bool("only-create") && c.Bool("only-update") {
				return fmt.Errorf("you can't use --only-create and --only-update at the same time")
			}

			if c.String("env-file") != "" {
				helpers.LoadEnvFile(c.String("env-file"))
			}
			utils.ReadCliArgs(c)

			if c.String("aws-account-id") != "" {
				console.Debug("Using AWS account ID provided in flag")
				os.Setenv("AWS_ACCOUNT_ID", c.String("aws-account_id"))
				return nil
			}

			account, accErr := aws.GetAccountID()
			if accErr != nil {
				// Let's not stop execution here
				console.Error(accErr.Error())
			}

			console.Debug("Account ID", account)
			os.Setenv("AWS_ACCOUNT_ID", account)

			return nil
		},
		Action: func(c *cli.Context) error {
			var projectPath = "project.json"
			if c.String("project") != "" {
				projectPath = c.String("project")
			} else {
				console.Info("Project config file path not specified. Assuming project.json")
			}

			if c.String("env-file") != "" {
				helpers.LoadEnvFile(c.String("env-file"))
			}

			config, err := helpers.LoadProject(projectPath)

			if err != nil {
				console.Error("Could not load project configuration")
				console.Fatal(err.Error())
				os.Exit(1)
			}

			utils.ReadCliArgs(c)

			stageTypesMap := make(map[string]bool)
			if c.String("stage-types") != "" {
				stageTypes := strings.Split(c.String("stage-types"), ",")

				for _, stageType := range stageTypes {
					stageTypesMap[stageType] = true
				}
			}

			var existingLambdas = make(map[string]lambdaTypes.FunctionConfiguration)
			if len(stageTypesMap) == 0 || stageTypesMap["lambda"] {
				console.Verbose("Listing lambdas in AWS")
				accountLambdas, err := aws.ListLambdas()
				if err != nil {
					console.Fatal("Could not list lambdas in AWS account. Check permissions ", err.Error())
				}
				existingLambdas = accountLambdas
				console.Verbosef("Found %d lambdas in AWS", len(accountLambdas))
			}

			var existingBuckets = make(map[string]bool)
			if len(stageTypesMap) == 0 || stageTypesMap["s3"] {
				console.Verbose("Listing S3 buckets in AWS")
				ctx, cfg, err := aws.GetConfig("us-east-1")

				if err != nil {
					return err
				}

				client := aws.GetClient(cfg)
				accountBuckets, bucketErr := aws.ListBuckets(ctx, client)
				if bucketErr != nil {
					console.Fatal("Could not list buckets in AWS account. Check permissions ", bucketErr.Error())
				}
				existingBuckets = accountBuckets
				console.Verbosef("Found %d S3 buckets in AWS", len(accountBuckets))
			}

			var existingApiGateways = make(map[string]string)

			if len(stageTypesMap) == 0 || stageTypesMap["api"] {
				console.Verbose("Listing API Gateways in AWS")
				accountApiGateways, gatewayErr := aws.ListApiGateways(os.Getenv("AWS_REGION"))
				if gatewayErr != nil {
					console.Fatal(gatewayErr.Error())
				}

				existingApiGateways = accountApiGateways
				console.Debugf("Found %d API gateways in AWS", len(accountApiGateways))
			}

			onlyCreate := c.Bool("only-create")
			onlyUpdate := c.Bool("only-update")

			propagationWaitTime := constants.DEFAULT_WAIT_TIME
			if c.Uint("wait-time") != 0 {
				propagationWaitTime = c.Int("wait-time")
			}
			if c.Bool("no-wait") {
				propagationWaitTime = 0
			}

			commands.HandleDeployCommand(config, &stageTypesMap, existingLambdas, existingBuckets, &existingApiGateways, onlyCreate, onlyUpdate, propagationWaitTime)

			console.Info("Done")
			return nil
		},
	}
}
