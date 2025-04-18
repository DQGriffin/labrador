package utils

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/DQGriffin/labrador/pkg/types"
	lambdaTypes "github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/urfave/cli/v2"
)

func ReadCliArgs(c *cli.Context) {
	if c.String("aws-access-key-id") != "" {
		os.Setenv("AWS_ACCESS_KEY_ID", c.String("aws-access-key-id"))
	}

	if c.String("aws-secret-access-key") != "" {
		os.Setenv("AWS_SECRET_ACCESS_KEY", c.String("aws-secret-access-key"))
	}

	if c.String("aws-region") != "" {
		os.Setenv("AWS_REGION", c.String("aws-region"))
	}
}

func BuildFunctionMap(fns []lambdaTypes.FunctionConfiguration) map[string]lambdaTypes.FunctionConfiguration {
	m := make(map[string]lambdaTypes.FunctionConfiguration, len(fns))
	for _, fn := range fns {
		if fn.FunctionName != nil {
			m[*fn.FunctionName] = fn
		}
	}
	return m
}

func ReadProjectData(filepath string) (types.Project, error) {
	var project types.Project

	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Failed to read project config")
		fmt.Println(err.Error())
		return project, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&project); err != nil {
		fmt.Println("Failed to decode project config")
		fmt.Println(err.Error())
		return project, err
	}

	return project, nil
}

func ReadFunctionConfig(filepath string) (types.LambdaData, error) {
	var functionData types.LambdaData

	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Failed to read function config")
		fmt.Println(err.Error())
		return functionData, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&functionData); err != nil {
		fmt.Println("Failed to decode function config")
		fmt.Println(err.Error())
		return functionData, err
	}

	return functionData, nil
}

func ReadFunctionConfigs(stages *[]types.Stage) ([]types.LambdaData, error) {
	var configs []types.LambdaData

	for i := range *stages {
		stage := &(*stages)[i]

		if stage.Type == "lambda" {
			data, err := ReadFunctionConfig(stage.ConfigFile)

			if err != nil {
				fmt.Println("Failed to read lambda config")
				return configs, err
			}
			stage.Functions = append(stage.Functions, data)
			configs = append(configs, data)
		}
	}

	return configs, nil
}

func ApplyDefaultsToFunctions(functionData *types.LambdaData) {
	for i := range functionData.Functions {
		applyDefaultsToFunction(&functionData.Functions[i], *functionData.Defaults)
	}
}

func applyDefaultsToFunction(function *types.LambdaConfig, defaults types.LambdaDefaults) {
	if (function.Code == nil || *function.Code == "") && defaults.Code != nil {
		function.Code = defaults.Code
	}

	if function.Environment == nil && defaults.Environment != nil {
		function.Environment = defaults.Environment
	}

	if function.Tags == nil && defaults.Tags != nil {
		function.Tags = defaults.Tags
	}

	if (function.Handler == nil || *function.Handler == "") && defaults.Handler != nil {
		function.Handler = defaults.Handler
	}

	if (function.Runtime == nil || *function.Runtime == "") && defaults.Runtime != nil {
		function.Runtime = defaults.Runtime
	}

	if (function.Region == nil || *function.Region == "") && defaults.Region != nil {
		function.Region = defaults.Region
	}

	if (function.RoleArn == nil || *function.RoleArn == "") && defaults.RoleArn != nil {
		function.RoleArn = defaults.RoleArn
	}

	if (function.MemorySize == nil || *function.MemorySize == 0) && defaults.MemorySize != nil {
		function.MemorySize = defaults.MemorySize
	}

	if (function.Timeout == nil || *function.Timeout == 0) && defaults.Timeout != nil {
		function.Timeout = defaults.Timeout
	}

	if (function.Description == nil || *function.Description == "") && defaults.Description != nil {
		function.Description = defaults.Description
	}
}
