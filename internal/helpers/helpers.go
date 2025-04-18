package helpers

import (
	"fmt"
	"os"

	"github.com/DQGriffin/labrador/internal/validation"
	"github.com/DQGriffin/labrador/pkg/interpolation"
	"github.com/DQGriffin/labrador/pkg/types"
	"github.com/DQGriffin/labrador/pkg/utils"
	"github.com/joho/godotenv"
)

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
		fmt.Printf("Unable to read project config from %s\n", filepath)
		fmt.Println(err.Error())
		os.Exit(1)
	}

	errs := validation.ValidateProject(project)
	if len(errs) > 0 {
		fmt.Println("Errors validating project config")
		for _, err := range errs {
			fmt.Println(err)
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

	config.Project = project
	return config, nil
}
