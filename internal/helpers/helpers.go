package helpers

import (
	"fmt"
	"os"

	"github.com/DQGriffin/labrador/internal/cli/console"
	"github.com/DQGriffin/labrador/internal/validation"
	"github.com/DQGriffin/labrador/pkg/interpolation"
	"github.com/DQGriffin/labrador/pkg/types"
	"github.com/DQGriffin/labrador/pkg/utils"
	"github.com/joho/godotenv"
)

func IsStageActionable(stage *types.Stage, stageTypesMap *map[string]bool) bool {
	if len(*stageTypesMap) == 0 {
		return true
	}

	return (*stageTypesMap)[stage.Type]
}

func PtrOrDefault[T any](ptr *T, fallback T) T {
	if ptr != nil {
		return *ptr
	}
	return fallback
}

func FirstNonNilOrDefault[T any](fallback T, ptrs ...*T) T {
	for _, ptr := range ptrs {
		if ptr != nil {
			return *ptr
		}
	}
	return fallback
}

func Filter[T any](input []T, predicate func(T) bool) []T {
	var result []T
	for _, item := range input {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

func LoadEnvFile(path string) error {
	if _, err := os.Stat(path); err == nil {
		return godotenv.Load(path)
	}
	return nil // silently skip if file not found
}

func LoadProject(filepath string) (types.LabradorConfig, error) {
	var config types.LabradorConfig

	project, err := utils.ReadProjectData(filepath)
	if err != nil {
		console.Errorf("Unable to read project config from %s\n", filepath)
		console.Fatal(err.Error())
	}

	errs := validation.ValidateProject(project)
	if len(errs) > 0 {
		console.Error("Errors validating project config")
		for _, err := range errs {
			console.Info(err)
		}
		os.Exit(1)
	}

	interpolation.InterpolateProjectVariables(&project)
	project.Variables["project_name"] = project.Name
	project.Variables["env"] = project.Environment

	functionData, readErr := utils.ReadFunctionConfigs(&project.Stages)

	if readErr != nil {
		fmt.Println(readErr)
		os.Exit(1)
	}

	for i := range functionData {
		interpolation.Interpolate(&functionData[i], project.Variables)
		utils.ApplyDefaultsToFunctions(&functionData[i])

		for functionIndex := range functionData[i].Functions {
			project.Variables["name"] = functionData[i].Functions[functionIndex].Name
			interpolation.Interpolate(&functionData[i].Functions[functionIndex], project.Variables)
		}

		config.FunctionData = append(config.FunctionData, functionData[i])
	}

	s3Configs, s3Err := utils.ReadS3Configs(&project.Stages)

	if s3Err != nil {
		fmt.Println(s3Err)
		os.Exit(1)
	}

	for i := range s3Configs {
		interpolation.Interpolate(&s3Configs[i], project.Variables)

		// for functionIndex := range functionData[i] {
		// 	project.Variables["name"] = functionData[i].Functions[functionIndex].Name
		// 	interpolation.Interpolate(&functionData[i].Functions[functionIndex], project.Variables)
		// }
	}

	gatewayConfigs, gatewayErr := utils.ReadApiGatewayConfigs(&project.Stages)
	if gatewayErr != nil {
		fmt.Println(gatewayErr)
		os.Exit(1)
	}

	for i := range gatewayConfigs {
		interpolation.Interpolate(&gatewayConfigs[i], project.Variables)
	}

	config.Project = project
	return config, nil
}
