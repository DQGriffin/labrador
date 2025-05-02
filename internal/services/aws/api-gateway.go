package aws

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

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
	console.Info("Created API: ", apiID)

	m := make(map[string]string)
	settingsErr := setApiGatewaySettings(gateway, &m, ctx, *client, apiID)

	if settingsErr != nil {
		console.Error(settingsErr.Error())
	}

	return nil
}

func UpdateApiGateway(gateway *types.ApiGatewaySettings, apiId string) error {
	console.Infof("Updating API Gateway %s", *gateway.Name)
	ctx := context.TODO()
	cfg, _ := config.LoadDefaultConfig(ctx, config.WithRegion(*gateway.Region))
	client := apigatewayv2.NewFromConfig(cfg)

	_, err := client.UpdateApi(ctx, &apigatewayv2.UpdateApiInput{
		ApiId:       aws.String(apiId),
		Description: aws.String(*gateway.Description),
	})

	if err != nil {
		return err
	}

	refMap := make(map[string]string)

	existingIntegrations, intErr := listIntegrations(&ctx, client, apiId)
	if intErr != nil {
		console.Debug("Something went wrong listing integrations")
		return intErr
	}

	existingRoutes, routeErr := ListRoutes(&ctx, client, apiId)
	if routeErr != nil {
		console.Debug("Something went wrong listing routes")
		return routeErr
	}

	deleteRoutes(&gateway.Routes, existingRoutes, ctx, client, apiId)
	deleteIntegrations(&existingIntegrations, &ctx, client, apiId)

	integrationRefs, err := addIntegrations(&gateway.Integrations, *gateway.Region, &refMap, ctx, *client, apiId)

	if err != nil {
		return err
	}

	routesErr := addRoutes(&gateway.Routes, &integrationRefs, ctx, *client, apiId)
	if routesErr != nil {
		return routesErr
	}

	console.Infof("Finished updating API Gateway %s", *gateway.Name)
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
	console.Verbose("Creating integrations")
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

		permErr := AddPermissionToLambda(ctx, cfg, *permission)
		if permErr != nil {
			if strings.Contains(permErr.Error(), "409") {
				console.Verbosef("Permission already exists for target %s", targetArn)
			} else {
				console.Error("failed to add permission to lambda: ", permErr.Error())
			}
		}

		console.Verbose("Created integration: ", integrationID)
	}

	console.Info("Finished creating integrations")
	return integrationRefMap, nil
}

func addRoutes(routes *[]types.ApiGatewayRoute, refMap *map[string]string, ctx context.Context, client apigatewayv2.Client, apiId string) error {
	console.Verbose("Creating routes")
	for _, route := range *routes {
		integrationId := (*refMap)[*route.Target.Ref]

		if integrationId == "" {
			return errors.New("failed to add route to API Gateway. Integration ref could not be resolved")
		}

		routeKey := fmt.Sprintf("%s %s", route.Method, route.Route)
		_, err := client.CreateRoute(ctx, &apigatewayv2.CreateRouteInput{
			ApiId:    aws.String(apiId),
			RouteKey: aws.String(routeKey),
			Target:   aws.String("integrations/" + integrationId),
		})

		if err != nil {
			return fmt.Errorf("failed to create route: %w", err)
		}

		console.Verbosef("Created route %s", routeKey)
	}

	console.Verbose("Finished creating routes")
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

		console.Infof("Created stage: %s", stage.Name)
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
	console.Infof("Deleting API Gateway: %s", gatewayName)
	apiId, err := GetApiIDByName(ctx, &client, gatewayName)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			console.Infof("API gateway %s did not exist. No action taken", gatewayName)
			return nil
		}
		return err
	}

	_, deleteErr := client.DeleteApi(ctx, &apigatewayv2.DeleteApiInput{
		ApiId: &apiId,
	})

	if deleteErr != nil {
		if strings.Contains(deleteErr.Error(), "409") {
			console.Infof("API gateway %s did not exist. No action taken", gatewayName)
			return nil
		}
		return deleteErr
	}

	console.Infof("Deleted API Gateway: %s", gatewayName)
	return nil
}

func ListApiGateways(region string) (map[string]string, error) {
	var existingApiGateways = make(map[string]string)

	ctx := context.TODO()
	cfg, cfgErr := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if cfgErr != nil {
		console.Error("Failed to load AWS configuration")
		return nil, cfgErr
	}

	client := apigatewayv2.NewFromConfig(cfg)
	input := &apigatewayv2.GetApisInput{}

	resp, err := client.GetApis(ctx, input)
	if err != nil {
		return nil, err
	}

	for _, api := range resp.Items {
		if api.Name == nil || api.ApiId == nil {
			console.Debug("An API Gateway could not be added to the existing gateway map because it's missing a name and/or API ID")
			continue
		}
		existingApiGateways[*api.Name] = *api.ApiId
	}

	// Handle pagination if needed
	for resp.NextToken != nil {
		input.NextToken = resp.NextToken
		resp, err = client.GetApis(ctx, input)
		if err != nil {
			return nil, err
		}

		for _, api := range resp.Items {
			if api.Name == nil || api.ApiId == nil {
				console.Debug("An API Gateway could not be added to the existing gateway map because it's missing a name and/or API ID")
				continue
			}
			existingApiGateways[*api.Name] = *api.ApiId
		}
	}

	return existingApiGateways, nil
}

func ListRoutes(ctx *context.Context, client *apigatewayv2.Client, apiID string) (map[string]gatewayTypes.Route, error) {
	var routes = make(map[string]gatewayTypes.Route)

	input := &apigatewayv2.GetRoutesInput{
		ApiId: &apiID,
	}

	resp, err := client.GetRoutes(*ctx, input)
	if err != nil {
		return nil, err
	}

	for _, route := range resp.Items {
		routes[*route.RouteKey] = route
		console.Debugf("Route Target: %s", *route.Target)
	}

	// Handle pagination if needed
	for resp.NextToken != nil {
		input.NextToken = resp.NextToken
		resp, err = client.GetRoutes(*ctx, input)
		if err != nil {
			return nil, err
		}

		for _, route := range resp.Items {
			routes[*route.RouteKey] = route
		}
	}

	return routes, nil
}

func listIntegrations(ctx *context.Context, client *apigatewayv2.Client, apiId string) (map[string]gatewayTypes.Integration, error) {
	var integrations = make(map[string]gatewayTypes.Integration)

	input := &apigatewayv2.GetIntegrationsInput{
		ApiId: &apiId,
	}

	resp, err := client.GetIntegrations(*ctx, input)
	if err != nil {
		return nil, err
	}

	for _, integration := range resp.Items {
		integrations[*integration.IntegrationId] = integration
	}

	// Handle pagination if needed
	for resp.NextToken != nil {
		input.NextToken = resp.NextToken
		resp, err = client.GetIntegrations(*ctx, input)
		if err != nil {
			return nil, err
		}

		for _, integration := range resp.Items {
			integrations[*integration.IntegrationId] = integration
		}
	}

	return integrations, nil
}

func deleteRoutes(routes *[]types.ApiGatewayRoute, existingRoutes map[string]gatewayTypes.Route, ctx context.Context, client *apigatewayv2.Client, apiID string) error {
	console.Verbose("Deleting routes...")
	for _, route := range *routes {
		routeKey := fmt.Sprintf("%s %s", route.Method, route.Route)
		console.Infof("Deleting route %s", routeKey)
		targetRoute := existingRoutes[routeKey]

		if targetRoute.RouteId == nil {
			console.Warnf("route %s not found, skipping delete", routeKey)
			continue
		}

		err := deleteRoute(ctx, client, apiID, *targetRoute.RouteId)
		if err != nil {
			return err
		}

		console.Verbosef("Finished deleting route %s", routeKey)
	}

	return nil
}

func deleteRoute(ctx context.Context, client *apigatewayv2.Client, apiID, routeID string) error {
	_, err := client.DeleteRoute(ctx, &apigatewayv2.DeleteRouteInput{
		ApiId:   aws.String(apiID),
		RouteId: aws.String(routeID),
	})
	return err
}

func deleteIntegrations(existingIntegrations *map[string]gatewayTypes.Integration, ctx *context.Context, client *apigatewayv2.Client, apiID string) {
	for _, integration := range *existingIntegrations {
		console.Verbosef("Deleting integration %s", *integration.IntegrationId)
		err := deleteIntegration(ctx, client, apiID, *integration.IntegrationId)
		if err != nil {
			console.Warnf("could not delete integration %s", *integration.IntegrationId)
			continue
		}
		console.Verbosef("Deleted integration %s", *integration.IntegrationId)
	}
}

func deleteIntegration(ctx *context.Context, client *apigatewayv2.Client, apiID, integrationID string) error {
	_, err := client.DeleteIntegration(*ctx, &apigatewayv2.DeleteIntegrationInput{
		ApiId:         aws.String(apiID),
		IntegrationId: aws.String(integrationID),
	})
	return err
}
