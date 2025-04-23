package commands

import (
	"fmt"

	"github.com/DQGriffin/labrador/internal/helpers"
	"github.com/DQGriffin/labrador/internal/services/aws"
	"github.com/DQGriffin/labrador/pkg/types"
	lambdaTypes "github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

func HandleDeployCommand(config types.LabradorConfig, stageTypesMap *map[string]bool, existingLambdas map[string]lambdaTypes.FunctionConfiguration, existingBuckets map[string]bool, onlyCreate bool, onlyUpdate bool) {
	for _, stage := range config.Project.Stages {

		if helpers.IsStageActionable(&stage, stageTypesMap) {
			deployLambdaStage(&stage, existingLambdas, onlyCreate, onlyUpdate)
			deployS3Stage(&stage, existingBuckets, onlyCreate, onlyUpdate)
			deployApiGatewayStage(&stage, onlyCreate, onlyUpdate)
		}
	}
}

func deployLambdaStage(stage *types.Stage, existingLambdas map[string]lambdaTypes.FunctionConfiguration, onlyCreate bool, onlyUpdate bool) {
	fmt.Printf("Deploying Stage: %s\n", stage.Name)
	fmt.Printf("Type: %s\n", stage.Type)

	for _, fnConfig := range stage.Functions {
		for _, fn := range fnConfig.Functions {
			if _, exists := existingLambdas[fn.Name]; exists {
				if !onlyCreate {
					fmt.Println("updating function", fn.Name)
					aws.UpdateLambda(fn)
				}
			} else {
				if !onlyUpdate {
					fmt.Println("creating function", fn.Name)
					aws.CreateLambda(fn)
				}
			}
		}
	}
}

func deployApiGatewayStage(stage *types.Stage, onlyCreate bool, onlyUpdate bool) {
	fmt.Printf("Deploying Stage: %s\n", stage.Name)
	fmt.Printf("Type: %s\n", stage.Type)

	for _, gatewayConfig := range stage.Gateways {
		for _, gateway := range gatewayConfig.Gateways {
			if !onlyUpdate {
				err := aws.CreateApiGateway(&gateway)
				if err != nil {
					fmt.Println(err.Error())
				}
			}
		}
	}
}

func deployS3Stage(stage *types.Stage, existingBuckets map[string]bool, onlyCreate bool, onlyUpdate bool) error {
	for _, bucketConfig := range stage.Buckets {
		for _, bucket := range bucketConfig.Buckets {
			ctx, cfg, err := aws.GetConfig(*bucket.Region)

			if err != nil {
				return err
			}

			client := aws.GetClient(cfg)

			if _, exists := existingBuckets[*bucket.Name]; exists {
				updateErr := aws.UpdateBucket(ctx, *client, bucket)
				if updateErr != nil {
					fmt.Println(updateErr.Error())
				}

			} else {
				createErr := aws.CreateBucket(ctx, cfg, *client, bucket)
				if createErr != nil {
					fmt.Println(createErr.Error())
				}

			}
		}
	}

	return nil
}
