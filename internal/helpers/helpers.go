package helpers

import (
	"os"
	"os/exec"

	"github.com/DQGriffin/labrador/internal/cli/console"
	"github.com/DQGriffin/labrador/internal/validation"
	"github.com/DQGriffin/labrador/pkg/interpolation"
	"github.com/DQGriffin/labrador/pkg/types"
	"github.com/DQGriffin/labrador/pkg/utils"
	"github.com/joho/godotenv"
)

func RunHooks(hookType, workingDir string, commands *[]string, suppressStdout, suppressStderr, stopOnError bool) {
	totalCommands := len(*commands)
	if totalCommands == 0 {
		return
	}

	console.Headingf("[%s hooks]", hookType)
	successfulCommands := 0

	for _, command := range *commands {
		execCmd := exec.Command("sh", "-c", command)

		// Set the working dir if specifiec, otherwise use cwd
		wd := ""
		if workingDir != "" {
			console.Debugf("Setting working directory to %s", workingDir)
			wd = workingDir
		} else {
			console.Debug("Working directory not specified. Setting to current working directory")
			cwd, err := os.Getwd()
			if err != nil {
				console.Errorf("hook cannot be executed. working directory could not be set. %s", err.Error())
				return
			}
			wd = cwd
		}

		if !suppressStdout {
			execCmd.Stdout = os.Stdout
		}

		if !suppressStderr {
			execCmd.Stderr = os.Stderr
		}

		console.Infof("> %s", command)

		execCmd.Dir = wd
		err := execCmd.Run()
		if err != nil {
			if stopOnError {
				console.Fatalf("%s hook %q failed: %s", hookType, command, err.Error())
			} else {
				console.Warnf("%s hook %q failed: %s", hookType, command, err.Error())
			}
		}
		successfulCommands += 1
		console.Info()
	}

	console.Infof("Finished running %s hooks", hookType)
	console.Infof("%d total, %d successful, %d failed", totalCommands, successfulCommands, totalCommands-successfulCommands)
}

func AsPtr[T any](v T) *T {
	return &v
}

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
		console.Fatal(readErr)
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
		console.Fatal(s3Err)
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
		console.Fatal(gatewayErr)
	}

	for i := range gatewayConfigs {
		interpolation.Interpolate(&gatewayConfigs[i], project.Variables)
	}

	iamRoleConfigs, roleErr := utils.ReadIamRoleConfigs(&project.Stages)
	if roleErr != nil {
		console.Fatal(roleErr)
	}

	for i := range iamRoleConfigs {
		interpolation.Interpolate(&iamRoleConfigs[i], project.Variables)
	}

	config.Project = project
	return config, nil
}
