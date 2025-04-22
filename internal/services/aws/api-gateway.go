package aws

import (
	"context"
	"fmt"

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
	fmt.Println("Created API:", apiID)

	return nil
}
