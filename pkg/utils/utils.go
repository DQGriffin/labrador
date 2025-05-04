package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/DQGriffin/labrador/internal/cli/console"
	"github.com/DQGriffin/labrador/internal/constants"
	"github.com/DQGriffin/labrador/internal/services/cognito"
	"github.com/DQGriffin/labrador/pkg/types"
	lambdaTypes "github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/urfave/cli/v2"
)

func ApplyDefaults[T any](target *T, defaults T) error {
	targetVal := reflect.ValueOf(target).Elem()
	defaultVal := reflect.ValueOf(defaults)

	if targetVal.Type() != defaultVal.Type() {
		console.Debugf("Target type: %s, defaults type: %s\n", targetVal.Type(), defaultVal.Type())
		return fmt.Errorf("target and defaults must be the same type")
	}

	for i := 0; i < targetVal.NumField(); i++ {
		// field := targetVal.Type().Field(i)

		// Only apply to exported fields
		if !targetVal.Field(i).CanSet() {
			console.Debug("Cannot set")
			continue
		}

		targetField := targetVal.Field(i)
		defaultField := defaultVal.Field(i)

		// Skip if already set
		if !targetField.IsZero() {
			continue
		}

		// Set the default value
		targetField.Set(defaultField)
	}

	return nil
}

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
		console.Error("Failed to read project config")
		console.Error(err.Error())
		return project, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&project); err != nil {
		console.Error("Failed to decode project config")
		console.Error(err.Error())
		return project, err
	}

	return project, nil
}

func ReadFunctionConfig(filepath string) (types.LambdaData, error) {
	var functionData types.LambdaData

	file, err := os.Open(filepath)
	if err != nil {
		console.Error("Failed to read function config")
		console.Error(err.Error())
		return functionData, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&functionData); err != nil {
		console.Error("Failed to decode function config")
		console.Error(err.Error())
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
				console.Error("Failed to read lambda config")
				return configs, err
			}
			stage.Functions = append(stage.Functions, data)
			configs = append(configs, data)
		}
	}

	return configs, nil
}

func ReadS3Configs(stages *[]types.Stage) ([]types.S3Config, error) {
	var configs []types.S3Config

	for i := range *stages {
		stage := &(*stages)[i]

		if stage.Type == "s3" {
			config, err := readS3Config(stage.ConfigFile)

			if err != nil {
				return configs, err
			}

			for i := range config.Buckets {
				ApplyDefaults(&config.Buckets[i], *config.Defaults)
			}

			configs = append(configs, config)
			stage.Buckets = append(stage.Buckets, config)
		}
	}

	return configs, nil
}

func readS3Config(filepath string) (types.S3Config, error) {
	var config types.S3Config

	file, err := os.Open(filepath)
	if err != nil {
		console.Error("Failed to read s3 config")
		console.Error(err.Error())
		return config, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		console.Error("Failed to decode s3 config")
		console.Error(err.Error())
		return config, err
	}

	return config, nil
}

func ReadApiGatewayConfigs(stages *[]types.Stage) ([]types.ApiGatewayConfig, error) {
	var configs []types.ApiGatewayConfig

	for i := range *stages {
		stage := &(*stages)[i]

		if stage.Type == "api" {
			config, err := readApiGatewayConfig(stage.ConfigFile)

			if err != nil {
				return configs, err
			}

			for i := range config.Gateways {
				ApplyDefaults(&config.Gateways[i], *config.Defaults)
			}

			configs = append(configs, config)
			stage.Gateways = append(stage.Gateways, config)
		}
	}

	return configs, nil
}

func readApiGatewayConfig(filepath string) (types.ApiGatewayConfig, error) {
	var config types.ApiGatewayConfig

	file, err := os.Open(filepath)
	if err != nil {
		console.Error("Failed to read API gateway config")
		console.Error(err.Error())
		return config, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		console.Error("Failed to decode API gateway config")
		console.Error(err.Error())
		return config, err
	}

	return config, nil
}

func ReadIamRoleConfigs(stages *[]types.Stage) ([]types.IamRoleConfig, error) {
	var configs []types.IamRoleConfig

	for i := range *stages {
		stage := &(*stages)[i]

		if stage.Type == "iam-role" {
			config, err := readIamRoleConfig(stage.ConfigFile)

			if err != nil {
				return configs, err
			}

			for i := range config.Roles {
				ApplyDefaults(&config.Roles[i], *config.Defaults)
			}

			configs = append(configs, config)
			stage.IamRoles = append(stage.IamRoles, config)
		}
	}

	return configs, nil
}

func readIamRoleConfig(filepath string) (types.IamRoleConfig, error) {
	var config types.IamRoleConfig

	file, err := os.Open(filepath)
	if err != nil {
		console.Error("Failed to read IAM role config")
		console.Error(err.Error())
		return config, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		console.Error("Failed to decode IAM role config")
		console.Error(err.Error())
		return config, err
	}

	return config, nil
}

func ReadCognitoConfigs(stages *[]types.Stage) ([]cognito.CognitoConfig, error) {
	var configs []cognito.CognitoConfig

	for i := range *stages {
		stage := &(*stages)[i]

		if stage.Type == constants.COGNITO_USER_POOL_STAGE {
			config, err := readCognitoConfig(stage.ConfigFile)

			if err != nil {
				return configs, err
			}

			for i := range config.Pools {
				ApplyDefaults(&config.Pools[i], *config.Defaults)
			}

			configs = append(configs, config)
			stage.UserPools = append(stage.UserPools, config)
		}
	}

	return configs, nil
}

func readCognitoConfig(filepath string) (cognito.CognitoConfig, error) {
	var config cognito.CognitoConfig

	file, err := os.Open(filepath)
	if err != nil {
		console.Error("Failed to read Cognito user pool config")
		console.Error(err.Error())
		return config, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		console.Error("Failed to decode Cognito user pool config")
		console.Error(err.Error())
		return config, err
	}

	return config, nil
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
