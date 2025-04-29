package add

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/DQGriffin/labrador/internal/cli/console"
	"github.com/DQGriffin/labrador/internal/helpers"
	"github.com/DQGriffin/labrador/pkg/types"
)

func HandleAddStage(projectPath, stageType, stageName, outputPath string) error {
	console.Debug("HandleAddLambdaStage")
	console.Debugf("Project: %s, Type: %s, Name: %s, Output: %s", projectPath, stageType, stageName, outputPath)

	config, err := helpers.LoadProject(projectPath)

	if err != nil {
		console.Error("Could not load project configuration")
		console.Fatal(err.Error())
	}

	stageErr := dispatchCommand(&config.Project, projectPath, stageName, stageType, outputPath)
	if stageErr != nil {
		return stageErr
	}

	return nil
}

func dispatchCommand(project *types.Project, projectPath, stageName, stageType, outputPath string) error {
	switch stageType {
	case "lambda":
		return handleAddLambdaStage(project, projectPath, stageName, outputPath)
	case "s3":
		return handleAddS3Stage(project, projectPath, stageName, outputPath)
	default:
		console.Fatalf("Cannot add stage of unknown type: %s", stageType)

		// This will never run as console.Fatalf() exits
		return fmt.Errorf("Cannot add stage of unknown type: %s", stageType)
	}
}

func handleAddLambdaStage(project *types.Project, projectPath, stageName, outputPath string) error {
	lambdaData := types.LambdaData{
		Defaults: &types.LambdaDefaults{
			Handler:    helpers.AsPtr("index.handler"),
			Runtime:    helpers.AsPtr("nodejs22.x"),
			Region:     helpers.AsPtr("us-east-1"),
			MemorySize: helpers.AsPtr(uint16(128)),
			Timeout:    helpers.AsPtr(uint16(3)),
			RoleArn:    helpers.AsPtr("[ROLE_ARN]"),
		},
		Functions: []types.LambdaConfig{
			{
				Name:        "{{env}}-{{project_name}}-func",
				Code:        helpers.AsPtr("sample-code.zip"),
				Description: helpers.AsPtr("{{project_name}} function"),
				OnDelete:    helpers.AsPtr("delete"),
				Tags: map[string]string{
					"app": "{{project_name}}",
				},
				Environment: map[string]string{
					"MY_ENV_VARIABLE": "my_value",
				},
			},
		},
	}

	data, err := json.MarshalIndent(lambdaData, "", "  ")
	if err != nil {
		console.Debug("Failed to marshal lambda data")
		return err
	}

	err = os.WriteFile(outputPath, data, 0644)
	if err != nil {
		console.Debug("Failed to write lambda config file")
		return err
	}

	stage := types.Stage{
		Name:         stageName,
		Type:         "lambda",
		Enabled:      true,
		OnConflict:   "stop",
		OnError:      "stop",
		ConfigFile:   outputPath,
		Environments: []string{"prod"},
	}

	project.Stages = append(project.Stages, stage)

	projectData, err := json.MarshalIndent(project, "", "  ")
	if err != nil {
		console.Debug("Failed to marshal project data")
		return err
	}

	err = os.WriteFile(projectPath, projectData, 0644)
	if err != nil {
		console.Debug("Failed to write project config file")
		return err
	}

	console.Infof("Added lambda stage %s to project %s", stageName, project.Name)
	console.Infof("Stage configuration saved to %s", outputPath)

	return nil
}

func handleAddS3Stage(project *types.Project, projectPath, stageName, outputPath string) error {
	bucketConfig := types.S3Config{
		Defaults: &types.S3Settings{
			Region:            helpers.AsPtr("us-east-1"),
			Versioning:        helpers.AsPtr(false),
			OnDelete:          helpers.AsPtr("delete"),
			BlockPublicAccess: helpers.AsPtr(true),
			StaticHosting: &types.StaticHostingSettings{
				Enabled: false,
			},
			Tags: map[string]string{
				"app": "{{project_name}}",
			},
		},
		Buckets: []types.S3Settings{
			{
				Name: helpers.AsPtr("{{env}}-{{project_name}}-assets"),
			},
		},
	}

	data, err := json.MarshalIndent(bucketConfig, "", "  ")
	if err != nil {
		console.Debug("Failed to marshal bucket data")
		return err
	}

	err = os.WriteFile(outputPath, data, 0644)
	if err != nil {
		console.Debug("Failed to write bucket config file")
		return err
	}

	stage := types.Stage{
		Name:         stageName,
		Type:         "s3",
		Enabled:      true,
		OnConflict:   "stop",
		OnError:      "stop",
		ConfigFile:   outputPath,
		Environments: []string{"prod"},
	}

	project.Stages = append(project.Stages, stage)

	projectData, err := json.MarshalIndent(project, "", "  ")
	if err != nil {
		console.Debug("Failed to marshal project data")
		return err
	}

	err = os.WriteFile(projectPath, projectData, 0644)
	if err != nil {
		console.Debug("Failed to write project config file")
		return err
	}

	console.Infof("Added s3 stage %s to project %s", stageName, project.Name)
	console.Infof("Stage configuration saved to %s", outputPath)

	return nil
}
