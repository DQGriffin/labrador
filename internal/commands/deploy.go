package commands

import (
	"fmt"
	"time"

	"github.com/DQGriffin/labrador/internal/cli/console"
	"github.com/DQGriffin/labrador/internal/helpers"
	"github.com/DQGriffin/labrador/internal/services/aws"
	"github.com/DQGriffin/labrador/pkg/types"
	lambdaTypes "github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

func HandleDeployCommand(config types.LabradorConfig, stageTypesMap *map[string]bool, existingLambdas map[string]lambdaTypes.FunctionConfiguration, existingBuckets map[string]bool, existingApiGateways *map[string]string, onlyCreate bool, onlyUpdate bool, propagationWaitTime int) {
	for _, stage := range config.Project.Stages {

		if helpers.IsStageActionable(&stage, stageTypesMap) {
			if stage.Hooks != nil {
				helpers.RunHooks("preDeploy", stage.Hooks.WorkingDir, &stage.Hooks.PreDeploy, stage.Hooks.SuppressStdout, stage.Hooks.SuppressStderr, stage.Hooks.StopOnError)
			}
			if stage.Type == "lambda" {
				deployLambdaStage(&stage, existingLambdas, onlyCreate, onlyUpdate)
			} else if stage.Type == "s3" {
				deployS3Stage(&stage, existingBuckets, onlyCreate, onlyUpdate)
			} else if stage.Type == "api" {
				deployApiGatewayStage(&stage, existingApiGateways, onlyCreate, onlyUpdate)
			} else if stage.Type == "iam-role" {
				deployIamRoleStage(&stage, onlyCreate, onlyUpdate)

				console.Info("Waiting to let changes propogate")
				time.Sleep(time.Duration(propagationWaitTime) * time.Second)
			} else {
				console.Warn("unknown stage type: ", stage.Type)
			}
			if stage.Hooks != nil {
				helpers.RunHooks("postDeploy", stage.Hooks.WorkingDir, &stage.Hooks.PostDeploy, stage.Hooks.SuppressStdout, stage.Hooks.SuppressStderr, stage.Hooks.StopOnError)
			}
		}
	}
}

func deployLambdaStage(stage *types.Stage, existingLambdas map[string]lambdaTypes.FunctionConfiguration, onlyCreate bool, onlyUpdate bool) {
	console.Heading(stage.ToHeader())

	for _, fnConfig := range stage.Functions {
		for _, fn := range fnConfig.Functions {
			if _, exists := existingLambdas[fn.Name]; exists {
				if onlyCreate {
					console.Debugf("Skipping updating lambda %s because --only-create is set", fn.Name)
					continue
				}

				aws.UpdateLambda(fn)
			} else {
				if onlyUpdate {
					console.Debugf("Skipping creating lambda %s because --only-update is set", fn.Name)
					continue
				}

				aws.CreateLambda(fn)
			}
		}
	}

	console.Info()
}

func deployApiGatewayStage(stage *types.Stage, existingApiGateways *map[string]string, onlyCreate bool, onlyUpdate bool) {
	console.Heading(stage.ToHeader())

	for _, gatewayConfig := range stage.Gateways {
		for _, gateway := range gatewayConfig.Gateways {

			apiId := (*existingApiGateways)[*gateway.Name]
			if apiId == "" {
				if onlyUpdate {
					console.Debugf("Skipping creating api gateway %s because --only-update is set", *gateway.Name)
					continue
				}

				err := aws.CreateApiGateway(&gateway)
				if err != nil {
					console.Error(err.Error())
				}
			} else {
				if onlyCreate {
					console.Debugf("Skipping updating api gateway %s because --only-create is set", *gateway.Name)
					continue
				}

				err := aws.UpdateApiGateway(&gateway, apiId)
				if err != nil {
					console.Error(err.Error())
				}
			}
		}
	}

	console.Info()
}

func deployS3Stage(stage *types.Stage, existingBuckets map[string]bool, onlyCreate bool, onlyUpdate bool) error {
	console.Heading(stage.ToHeader())

	for _, bucketConfig := range stage.Buckets {
		for _, bucket := range bucketConfig.Buckets {
			ctx, cfg, err := aws.GetConfig(*bucket.Region)

			if err != nil {
				return err
			}

			client := aws.GetClient(cfg)

			if _, exists := existingBuckets[*bucket.Name]; exists {
				if onlyCreate {
					console.Debugf("Skipping updating bucket %s because --only-create is set", *bucket.Name)
					continue
				}

				updateErr := aws.UpdateBucket(ctx, *client, bucket)
				if updateErr != nil {
					fmt.Println(updateErr.Error())
				}

			} else {
				if onlyUpdate {
					console.Debugf("Skipping creating bucket %s because --only-update is set", *bucket.Name)
					continue
				}
				createErr := aws.CreateBucket(ctx, cfg, *client, bucket)
				if createErr != nil {
					fmt.Println(createErr.Error())
				}

			}
		}
	}

	console.Info()
	return nil
}

func deployIamRoleStage(stage *types.Stage, onlyCreate, onlyUpdate bool) error {
	console.Heading(stage.ToHeader())

	for _, config := range stage.IamRoles {
		for _, role := range config.Roles {
			if !onlyUpdate {
				err := aws.CreateIamRole(&role)
				if err != nil {
					console.Errorf("failed to create IAM role: %s", err.Error())
				}
			}
		}
	}

	console.Info()
	return nil
}
