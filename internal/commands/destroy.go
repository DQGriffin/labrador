package commands

import (
	"context"

	"github.com/DQGriffin/labrador/internal/cli/console"
	"github.com/DQGriffin/labrador/internal/services/aws"
	internalTypes "github.com/DQGriffin/labrador/internal/types"
	"github.com/DQGriffin/labrador/pkg/types"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/apigatewayv2"
)

func HandleDestroyCommand(projectConfig types.LabradorConfig, isDryRun bool, force bool, stageTypesMap *map[string]bool, env string) error {
	for _, stage := range projectConfig.Project.Stages {
		if isStageMarkedForDeletion(&stage, stageTypesMap, env) {
			if stage.Type == "lambda" {
				handleLambdaStage(&stage, isDryRun, force)
			} else if stage.Type == "s3" {
				handleS3Stage(&stage, isDryRun, force)
			} else if stage.Type == "api" {
				handleApiGatewayStage(&stage, isDryRun, force)
			}
		} else {
			console.Info("Skipping stage", stage.Name)
		}
	}

	return nil
}

func handleLambdaStage(stage *types.Stage, isDryRun bool, force bool) {
	console.Info("Stage", stage.Name)
	deletableLambdas, skippedLambdas := getDeletableLambdas(&stage.Functions, stage.Name)

	if isDryRun {
		handleDryRun(&deletableLambdas, &skippedLambdas)
	} else {
		destroyResources(&deletableLambdas, force)
	}
}

func handleS3Stage(stage *types.Stage, isDryRun bool, force bool) {
	console.Info("Stage", stage.Name)
	deletableBuckets, skippedBuckets := getDeletableBuckets(&stage.Buckets, stage.Name)

	if isDryRun {
		handleDryRun(&deletableBuckets, &skippedBuckets)
	} else {
		destroyResources(&deletableBuckets, force)
	}
}

func handleApiGatewayStage(stage *types.Stage, isDryRun bool, force bool) {
	console.Info("Stage", stage.Name)
	deletableGateways, skippedGateways := getDeletableApiGateways(&stage.Gateways, stage.Name)

	if isDryRun {
		handleDryRun(&deletableGateways, &skippedGateways)
	} else {
		destroyResources(&deletableGateways, force)
	}
}

func handleDryRun(forDeletion *[]internalTypes.UniversalResourceDefinition, skipped *[]internalTypes.UniversalResourceDefinition) {
	console.Info("Would delete:")
	for _, resource := range *forDeletion {
		console.Infof("- %s", resource.Name)
	}
	console.Info("Would skip:")
	for _, resource := range *skipped {
		console.Infof("- %s", resource.Name)
	}
	console.Info()
}

func isStageMarkedForDeletion(stage *types.Stage, stageTypesMap *map[string]bool, env string) bool {
	if len(*stageTypesMap) == 0 {
		console.Info("Stage types map is empty. Returning true")
		return true
	}

	return (*stageTypesMap)[stage.Type]
}

func getDeletableLambdas(config *[]types.LambdaData, stageName string) ([]internalTypes.UniversalResourceDefinition, []internalTypes.UniversalResourceDefinition) {
	var deletableLambdas []internalTypes.UniversalResourceDefinition
	var skippedLambdas []internalTypes.UniversalResourceDefinition
	for _, stageFuncs := range *config {
		for _, fn := range stageFuncs.Functions {
			if (fn.OnDelete == nil) || (fn.OnDelete != nil && *fn.OnDelete != "skip") {
				deletableLambdas = append(deletableLambdas, internalTypes.UniversalResourceDefinition{
					Name:         fn.Name,
					StageName:    stageName,
					Arn:          "",
					ResourceType: "lambda",
				})
			} else {
				skippedLambdas = append(skippedLambdas, internalTypes.UniversalResourceDefinition{
					Name:         fn.Name,
					StageName:    stageName,
					Arn:          "",
					ResourceType: "lambda",
				})
			}
		}
	}

	return deletableLambdas, skippedLambdas
}

func getDeletableBuckets(config *[]types.S3Config, stageName string) ([]internalTypes.UniversalResourceDefinition, []internalTypes.UniversalResourceDefinition) {
	var deletableBuckets []internalTypes.UniversalResourceDefinition
	var skippedBuckets []internalTypes.UniversalResourceDefinition

	for _, stageBuckets := range *config {
		for _, bucket := range stageBuckets.Buckets {
			if (bucket.OnDelete == nil) || (bucket.OnDelete != nil && *bucket.OnDelete != "skip") {
				deletableBuckets = append(deletableBuckets, internalTypes.UniversalResourceDefinition{
					Name:         *bucket.Name,
					StageName:    stageName,
					Arn:          "",
					ResourceType: "s3",
				})
			} else {
				skippedBuckets = append(skippedBuckets, internalTypes.UniversalResourceDefinition{
					Name:         *bucket.Name,
					StageName:    stageName,
					Arn:          "",
					ResourceType: "s3",
				})
			}
		}
	}

	return deletableBuckets, skippedBuckets
}

func getDeletableApiGateways(config *[]types.ApiGatewayConfig, stageName string) ([]internalTypes.UniversalResourceDefinition, []internalTypes.UniversalResourceDefinition) {
	var deletableGateways []internalTypes.UniversalResourceDefinition
	var skippedGateways []internalTypes.UniversalResourceDefinition

	for _, stageGateways := range *config {
		for _, gateway := range stageGateways.Gateways {
			if (gateway.OnDelete == nil) || (gateway.OnDelete != nil && *gateway.OnDelete != "skip") {
				deletableGateways = append(deletableGateways, internalTypes.UniversalResourceDefinition{
					Name:         *gateway.Name,
					StageName:    stageName,
					Arn:          "",
					ResourceType: "api",
					Region:       *gateway.Region,
				})
			} else {
				skippedGateways = append(skippedGateways, internalTypes.UniversalResourceDefinition{
					Name:         *gateway.Name,
					StageName:    stageName,
					Arn:          "",
					ResourceType: "api",
					Region:       *gateway.Region,
				})
			}
		}
	}

	return deletableGateways, skippedGateways
}

func destroyResources(resources *[]internalTypes.UniversalResourceDefinition, force bool) {
	for _, resource := range *resources {
		if resource.ResourceType == "lambda" {
			aws.DeleteLambda(resource.Name)
		} else if resource.ResourceType == "s3" {
			aws.DeleteBucket(resource.Name, force)
		} else if resource.ResourceType == "api" {

			ctx := context.TODO()
			cfg, _ := config.LoadDefaultConfig(ctx, config.WithRegion(resource.Region))
			client := apigatewayv2.NewFromConfig(cfg)
			err := aws.DestroyApiGateway(ctx, *client, resource.Name)
			if err != nil {
				console.Error(err.Error())
			}
		}
	}
}
