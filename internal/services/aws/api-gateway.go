package aws

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/DQGriffin/labrador/internal/cli/console"
	internalTypes "github.com/DQGriffin/labrador/internal/types"
	"github.com/DQGriffin/labrador/pkg/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/apigatewayv2"
	gatewayTypes "github.com/aws/aws-sdk-go-v2/service/apigatewayv2/types"
)

func CreateApiGateway(gateway *types.ApiGatewaySettings) error {
	ctx := context.TODO()
	cfg, _ := config.LoadDefaultConfig(ctx, config.WithRegion(*gateway.Region))
	client := apigatewayv2.NewFromConfig(cfg)

	apiOut, err := client.CreateApi(ctx, &apigatewayv2.CreateApiInput{
		Name:         aws.String(*gateway.Name),
		ProtocolType: gatewayTypes.ProtocolTypeHttp,
		Description:  aws.String(*gateway.Description),
		Tags:         gateway.Tags,
	})
	if err != nil {
		return fmt.Errorf("failed to create API: %w", err)
	}
	apiID := *apiOut.ApiId
	console.Info("Created API:", apiID)

	m := make(map[string]string)
	settingsErr := setApiGatewaySettings(gateway, &m, ctx, *client, apiID)

	if settingsErr != nil {
		console.Error(settingsErr.Error())
	}

	return nil
}

func setApiGatewaySettings(gateway *types.ApiGatewaySettings, refMap *map[string]string, ctx context.Context, client apigatewayv2.Client, apiId string) error {
	stageErr := createStages(gateway.Stages, ctx, client, apiId)

	if stageErr != nil {
		return stageErr
	}

	integrationRefs, err := addIntegrations(&gateway.Integrations, *gateway.Region, refMap, ctx, client, apiId)

	if err != nil {
		return err
	}

	routeErr := addRoutes(&gateway.Routes, &integrationRefs, ctx, client, apiId)
	if routeErr != nil {
		return routeErr
	}

	return nil
}

func addIntegrations(integrations *[]types.ApiGatewayIntegration, region string, refMap *map[string]string, ctx context.Context, client apigatewayv2.Client, apiId string) (map[string]string, error) {
	integrationRefMap := make(map[string]string)

	for _, integration := range *integrations {
		targetArn, arnErr := ResolveTarget(integration.Target, *refMap)

		if arnErr != nil {
			return integrationRefMap, arnErr
		}

		intOut, err := client.CreateIntegration(ctx, &apigatewayv2.CreateIntegrationInput{
			ApiId:                aws.String(apiId),
			IntegrationType:      gatewayTypes.IntegrationTypeAwsProxy,
			IntegrationUri:       aws.String(targetArn),
			IntegrationMethod:    aws.String("POST"),
			PayloadFormatVersion: aws.String("2.0"),
		})

		if err != nil {
			return integrationRefMap, fmt.Errorf("failed to create integration: %w", err)
		}

		integrationID := *intOut.IntegrationId

		if integration.Ref != "" {
			integrationRefMap[integration.Ref] = integrationID
		}

		accountId := os.Getenv("AWS_ACCOUNT_ID")
		arn := fmt.Sprintf("arn:aws:execute-api:%s:%s:%s/*/*/*", region, accountId, apiId)

		permission := &internalTypes.LambdaPermission{
			FunctionName: integration.Target.External.Dynamic.Name,
			Action:       "lambda:InvokeFunction",
			Principal:    "apigateway.amazonaws.com",
			StatementId:  fmt.Sprintf("apigateway-%s-invoke", apiId),
			SourceArn:    arn,
		}

		ctx, cfg, err := GetConfig(region)
		if err != nil {
			return integrationRefMap, fmt.Errorf("failed to add permission to lambda: %w", err)
		}

		AddPermissionToLambda(ctx, cfg, *permission)

		console.Info("Created integration:", integrationID)
	}

	return integrationRefMap, nil
}

func addRoutes(routes *[]types.ApiGatewayRoute, refMap *map[string]string, ctx context.Context, client apigatewayv2.Client, apiId string) error {
	for _, route := range *routes {
		integrationId := (*refMap)[*route.Target.Ref]

		if integrationId == "" {
			return errors.New("failed to add route to API Gateway. Integration ref could not be resolved")
		}

		_, err := client.CreateRoute(ctx, &apigatewayv2.CreateRouteInput{
			ApiId:    aws.String(apiId),
			RouteKey: aws.String(route.Method + " " + route.Route),
			Target:   aws.String("integrations/" + integrationId),
		})

		if err != nil {
			return fmt.Errorf("failed to create route: %w", err)
		}

		console.Info("Created route GET /users")
	}

	return nil
}

func createStages(stages *[]types.ApiGatewayStage, ctx context.Context, client apigatewayv2.Client, apiId string) error {
	for _, stage := range *stages {
		_, err := client.CreateStage(ctx, &apigatewayv2.CreateStageInput{
			ApiId:       aws.String(apiId),
			StageName:   aws.String(stage.Name),
			Description: aws.String(stage.Description),
			AutoDeploy:  aws.Bool(stage.AutoDeploy),
			Tags:        stage.Tags,
		})
		if err != nil {
			return fmt.Errorf("failed to create stage %q: %w", stage.Name, err)
		}

		console.Infof("Created stage: %s\n", stage.Name)
	}
	return nil
}

func GetApiIDByName(ctx context.Context, client *apigatewayv2.Client, targetName string) (string, error) {
	output, err := client.GetApis(ctx, &apigatewayv2.GetApisInput{})
	if err != nil {
		return "", fmt.Errorf("failed to list APIs: %w", err)
	}

	for _, api := range output.Items {
		if api.Name != nil && *api.Name == targetName {
			return *api.ApiId, nil
		}
	}

	return "", fmt.Errorf("API with name %q not found", targetName)
}

func DestroyApiGateway(ctx context.Context, client apigatewayv2.Client, gatewayName string) error {
	console.Infof("Deleting API Gateway %s\n", gatewayName)
	apiId, err := GetApiIDByName(ctx, &client, gatewayName)
	if err != nil {
		return err
	}

	_, deleteErr := client.DeleteApi(ctx, &apigatewayv2.DeleteApiInput{
		ApiId: &apiId,
	})

	if deleteErr != nil {
		return deleteErr
	}

	return nil
}
