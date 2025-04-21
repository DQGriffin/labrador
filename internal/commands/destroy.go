package commands

import (
	"fmt"

	"github.com/DQGriffin/labrador/internal/services/aws"
	internalTypes "github.com/DQGriffin/labrador/internal/types"
	"github.com/DQGriffin/labrador/pkg/types"
)

func HandleDestroyCommand(projectConfig types.LabradorConfig, isDryRun bool, env string) error {
	for _, stage := range projectConfig.Project.Stages {
		if isStageMarkedForDeletion(&stage, env) {
			fmt.Println("Stage", stage.Name)
			deletableLambdas, skippedLambdas := getDeletableLambdas(&stage.Functions, stage.Name)
			deletableBuckets, skippedBuckets := getDeletableBuckets(&stage.Buckets, stage.Name)

			if isDryRun {
				handleDryRun(&deletableLambdas, &skippedLambdas)
				handleDryRun(&deletableBuckets, &skippedBuckets)
			} else {
				destroyResources(&deletableLambdas)
				destroyResources(&deletableBuckets)
			}
		} else {
			fmt.Println("Skipping stage", stage.Name)
		}
	}

	return nil
}

func handleDryRun(forDeletion *[]internalTypes.UniversalResourceDefinition, skipped *[]internalTypes.UniversalResourceDefinition) {
	fmt.Println("Would delete:")
	for _, resource := range *forDeletion {
		fmt.Printf("- %s\n", resource.Name)
	}
	fmt.Println("Would skip:")
	for _, resource := range *skipped {
		fmt.Printf("- %s\n", resource.Name)
	}
	fmt.Println()
}

func isStageMarkedForDeletion(stage *types.Stage, env string) bool {
	// Not yet implemented.
	return true
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
					ResourceType: "lambda",
				})
			} else {
				skippedBuckets = append(skippedBuckets, internalTypes.UniversalResourceDefinition{
					Name:         *bucket.Name,
					StageName:    stageName,
					Arn:          "",
					ResourceType: "lambda",
				})
			}
		}
	}

	return deletableBuckets, skippedBuckets
}

func destroyResources(resources *[]internalTypes.UniversalResourceDefinition) {
	for _, resource := range *resources {
		if resource.ResourceType == "lambda" {
			aws.DeleteLambda(resource.Name)
		} else if resource.ResourceType == "s3" {
			aws.DeleteBucket(resource.Name)
		}
	}
}
