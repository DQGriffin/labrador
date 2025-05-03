package commands

import (
	"context"

	"github.com/DQGriffin/labrador/internal/cli/console"
	"github.com/DQGriffin/labrador/internal/helpers"
	"github.com/DQGriffin/labrador/internal/services/aws"
	internalTypes "github.com/DQGriffin/labrador/internal/types"
	"github.com/DQGriffin/labrador/pkg/types"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/apigatewayv2"
)

func HandleDestroyCommand(projectConfig types.LabradorConfig, isDryRun bool, force bool, stageTypesMap *map[string]bool, env string) error {
	for _, stage := range projectConfig.Project.Stages {
		if isStageMarkedForDeletion(&stage, stageTypesMap, env) {
			if stage.Hooks != nil {
				helpers.RunHooks("preDestroy", stage.Hooks.WorkingDir, &stage.Hooks.PreDestroy, stage.Hooks.SuppressStdout, stage.Hooks.SuppressStderr, stage.Hooks.StopOnError)
			}

			if stage.Type == "lambda" {
				handleLambdaStage(&stage, isDryRun, force)
			} else if stage.Type == "s3" {
				handleS3Stage(&stage, isDryRun, force)
			} else if stage.Type == "api" {
				handleApiGatewayStage(&stage, isDryRun, force)
			} else if stage.Type == "iam-role" {
				handleIamRoleStage(&stage, isDryRun, force)
			}

			if stage.Hooks != nil {
				helpers.RunHooks("postDestroy", stage.Hooks.WorkingDir, &stage.Hooks.PostDestroy, stage.Hooks.SuppressStdout, stage.Hooks.SuppressStderr, stage.Hooks.StopOnError)
			}
		} else {
			console.Debug("Skipping stage: ", stage.Name)
		}
	}

	console.Info("\nDone")
	return nil
}

func handleLambdaStage(stage *types.Stage, isDryRun bool, force bool) {
	console.Headingf("[Stage - %s - %s]", stage.Name, stage.Type)
	deletableLambdas, skippedLambdas := getDeletableLambdas(&stage.Functions, stage.Name)

	if isDryRun {
		handleDryRun(&deletableLambdas, &skippedLambdas)
	} else {
		destroyResources(&deletableLambdas, force)
	}
}

func handleS3Stage(stage *types.Stage, isDryRun bool, force bool) {
	console.Headingf("[Stage - %s - %s]", stage.Name, stage.Type)
	deletableBuckets, skippedBuckets := getDeletableBuckets(&stage.Buckets, stage.Name)

	if isDryRun {
		handleDryRun(&deletableBuckets, &skippedBuckets)
	} else {
		destroyResources(&deletableBuckets, force)
	}
}

func handleApiGatewayStage(stage *types.Stage, isDryRun bool, force bool) {
	console.Headingf("[Stage - %s - %s]", stage.Name, stage.Type)
	deletableGateways, skippedGateways := getDeletableApiGateways(&stage.Gateways, stage.Name)

	if isDryRun {
		handleDryRun(&deletableGateways, &skippedGateways)
	} else {
		destroyResources(&deletableGateways, force)
	}
}

func handleIamRoleStage(stage *types.Stage, isDryRun bool, force bool) {
	console.Headingf("[Stage - %s - %s]", stage.Name, stage.Type)
	deletableRoles, skippedRoles := getDeletableRoles(&stage.IamRoles, stage.Name)

	if isDryRun {
		handleDryRun(&deletableRoles, &skippedRoles)
	} else {
		destroyResources(&deletableRoles, force)
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
					Region:       *bucket.Region,
				})
			} else {
				skippedBuckets = append(skippedBuckets, internalTypes.UniversalResourceDefinition{
					Name:         *bucket.Name,
					StageName:    stageName,
					Arn:          "",
					ResourceType: "s3",
					Region:       *bucket.Region,
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

func getDeletableRoles(config *[]types.IamRoleConfig, stageName string) ([]internalTypes.UniversalResourceDefinition, []internalTypes.UniversalResourceDefinition) {
	var deletableRoles []internalTypes.UniversalResourceDefinition
	var skippedRoles []internalTypes.UniversalResourceDefinition

	for _, roleConfigs := range *config {
		for _, role := range roleConfigs.Roles {
			if (role.OnDelete == nil) || (role.OnDelete != nil && *role.OnDelete != "skip") {
				deletableRoles = append(deletableRoles, internalTypes.UniversalResourceDefinition{
					Name:         *role.Name,
					StageName:    stageName,
					Arn:          "",
					ResourceType: "iam-role",
					Region:       "",
				})
			} else {
				skippedRoles = append(skippedRoles, internalTypes.UniversalResourceDefinition{
					Name:         *role.Name,
					StageName:    stageName,
					Arn:          "",
					ResourceType: "iam-role",
					Region:       "",
				})
			}
		}
	}

	return deletableRoles, skippedRoles
}

func destroyResources(resources *[]internalTypes.UniversalResourceDefinition, force bool) {
	for _, resource := range *resources {
		if resource.ResourceType == "lambda" {
			aws.DeleteLambda(resource.Name)
		} else if resource.ResourceType == "s3" {
			err := aws.DeleteBucket(resource.Name, resource.Region, force)
			if err != nil {
				console.Error(err.Error())
			}
		} else if resource.ResourceType == "api" {

			ctx := context.TODO()
			cfg, _ := config.LoadDefaultConfig(ctx, config.WithRegion(resource.Region))
			client := apigatewayv2.NewFromConfig(cfg)
			err := aws.DestroyApiGateway(ctx, *client, resource.Name)
			if err != nil {
				console.Error(err.Error())
			}
		} else if resource.ResourceType == "iam-role" {
			err := aws.DeleteRole(resource.Name)
			if err != nil {
				console.Error(err.Error())
			}
		}
	}
}
