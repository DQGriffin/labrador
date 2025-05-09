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
	case "api":
		return handleAddApiGatewayStage(project, projectPath, stageName, outputPath)
	default:
		console.Fatalf("Cannot add stage of unknown type: %s", stageType)

		// This will never run as console.Fatalf() exits
		return fmt.Errorf("Cannot add stage of unknown type: %s", stageType)
	}
}

func handleAddLambdaStage(project *types.Project, projectPath, stageName, outputPath string) error {
	console.Debug("Adding lambda stage")

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

	data, err := json.MarshalIndent(lambdaData, "", "\t")
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

	projectData, err := json.MarshalIndent(project, "", "\t")
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
	console.Debug("Finished adding lambda stage")

	return nil
}

func handleAddS3Stage(project *types.Project, projectPath, stageName, outputPath string) error {
	console.Debug("Adding s3 stage")

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

	data, err := json.MarshalIndent(bucketConfig, "", "\t")
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

	projectData, err := json.MarshalIndent(project, "", "\t")
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
	console.Debug("Finished adding s3 stage")

	return nil
}

func handleAddApiGatewayStage(project *types.Project, projectPath, stageName, outputPath string) error {
	console.Debug("Adding api stage")

	gateway := types.ApiGatewayConfig{
		Defaults: &types.ApiGatewaySettings{
			OnDelete: helpers.AsPtr("delete"),
			Region:   helpers.AsPtr("us-east-1"),
			Protocol: helpers.AsPtr("http"),
			Tags: map[string]string{
				"app": "{{project_name}}",
			},
		},
		Gateways: []types.ApiGatewaySettings{
			{
				Name:        helpers.AsPtr("{{env}}-{{project_name}}-api"),
				Description: helpers.AsPtr("{{project_name}} API gateway"),
				Stages: &[]types.ApiGatewayStage{
					{
						Name:        "$default",
						Description: "Default stage",
						AutoDeploy:  true,
					},
				},
				Integrations: []types.ApiGatewayIntegration{
					{
						Type:              "proxy",
						PayloadVersion:    "2.0",
						IntegrationMethod: "POST",
						Ref:               "my-integration-reference",
						Target: types.ResourceTarget{
							External: &types.ExternalReference{
								Dynamic: &types.DynamicResourceRefData{
									Name:   "{{env}}-{{project_name}}-func",
									Region: "us-east-1",
									Type:   "lambda",
								},
							},
						},
					},
				},
				Routes: []types.ApiGatewayRoute{
					{
						Method: "GET",
						Route:  "/users",
						Target: types.ResourceTarget{
							Ref: helpers.AsPtr("my-integration-reference"),
						},
					},
				},
			},
		},
	}

	data, err := json.MarshalIndent(gateway, "", "\t")
	if err != nil {
		console.Debug("Failed to marshal api gateway data data")
		return err
	}

	err = os.WriteFile(outputPath, data, 0644)
	if err != nil {
		console.Debug("Failed to write api gateway config file")
		return err
	}

	stage := types.Stage{
		Name:         stageName,
		Type:         "api",
		Enabled:      true,
		OnConflict:   "stop",
		OnError:      "stop",
		ConfigFile:   outputPath,
		Environments: []string{"prod"},
	}

	project.Stages = append(project.Stages, stage)

	projectData, err := json.MarshalIndent(project, "", "\t")
	if err != nil {
		console.Debug("Failed to marshal project data")
		return err
	}

	err = os.WriteFile(projectPath, projectData, 0644)
	if err != nil {
		console.Debug("Failed to write project config file")
		return err
	}

	console.Infof("Added api stage %s to project %s", stageName, project.Name)
	console.Infof("Stage configuration saved to %s", outputPath)
	console.Debug("Finished adding api stage")

	return nil
}
