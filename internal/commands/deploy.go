package commands

import (
	"fmt"

	"github.com/DQGriffin/labrador/internal/cli/console"
	"github.com/DQGriffin/labrador/internal/helpers"
	"github.com/DQGriffin/labrador/internal/services/aws"
	"github.com/DQGriffin/labrador/pkg/types"
	lambdaTypes "github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

func HandleDeployCommand(config types.LabradorConfig, stageTypesMap *map[string]bool, existingLambdas map[string]lambdaTypes.FunctionConfiguration, existingBuckets map[string]bool, onlyCreate bool, onlyUpdate bool) {
	for _, stage := range config.Project.Stages {

		if helpers.IsStageActionable(&stage, stageTypesMap) {
			if stage.Type == "lambda" {
				deployLambdaStage(&stage, existingLambdas, onlyCreate, onlyUpdate)
			} else if stage.Type == "s3" {
				deployS3Stage(&stage, existingBuckets, onlyCreate, onlyUpdate)
			} else if stage.Type == "api" {
				deployApiGatewayStage(&stage, onlyCreate, onlyUpdate)
			} else {
				console.Warn("unknown stage type: ", stage.Type)
			}
		}
	}
}

func deployLambdaStage(stage *types.Stage, existingLambdas map[string]lambdaTypes.FunctionConfiguration, onlyCreate bool, onlyUpdate bool) {
	console.Headingf("[Stage - %s - %s]", stage.Name, stage.Type)

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

func deployApiGatewayStage(stage *types.Stage, onlyCreate bool, onlyUpdate bool) {
	console.Headingf("[Stage - %s - %s]", stage.Name, stage.Type)

	for _, gatewayConfig := range stage.Gateways {
		for _, gateway := range gatewayConfig.Gateways {

			if onlyUpdate {
				console.Debugf("Skipping creating api gateway %s because --only-update is set", *gateway.Name)
				continue
			}

			err := aws.CreateApiGateway(&gateway)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}

	console.Info()
}

func deployS3Stage(stage *types.Stage, existingBuckets map[string]bool, onlyCreate bool, onlyUpdate bool) error {
	console.Headingf("[Stage - %s - %s]", stage.Name, stage.Type)

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
