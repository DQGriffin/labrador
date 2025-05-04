package cognito

import (
	"context"
	"log"

	"github.com/DQGriffin/labrador/internal/cli/console"
	"github.com/DQGriffin/labrador/internal/refs"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	awsCognito "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

func CreateUserPool(pool *CognitoSettings) error {
	console.Infof("Creating Cognito user pool %s", *pool.ApplicationName)

	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	client := awsCognito.NewFromConfig(cfg)

	var schema []types.SchemaAttributeType
	for _, attribute := range *pool.SignUpAttributes {
		schema = append(schema, types.SchemaAttributeType{
			Name:              aws.String(attribute),
			Required:          aws.Bool(true),
			Mutable:           aws.Bool(true),                // or false if you'd like to prevent updates
			AttributeDataType: types.AttributeDataTypeString, // assuming strings
		})
	}

	output, createErr := client.CreateUserPool(ctx, &awsCognito.CreateUserPoolInput{
		PoolName: aws.String(*pool.ApplicationName),
		Policies: &types.UserPoolPolicyType{
			PasswordPolicy: &types.PasswordPolicyType{
				MinimumLength:    aws.Int32(int32(pool.PasswordRequirements.MinLength)),
				RequireSymbols:   *aws.Bool(pool.PasswordRequirements.RequireSymbols),
				RequireNumbers:   *aws.Bool(pool.PasswordRequirements.RequireNumbers),
				RequireLowercase: *aws.Bool(pool.PasswordRequirements.RequireLowercase),
				RequireUppercase: *aws.Bool(pool.PasswordRequirements.RequireUppercase),
			},
		},
		Schema:       schema,
		UserPoolTags: pool.Tags,
	})
	if createErr != nil {
		log.Fatalf("Failed to create user pool: %v", createErr)
	}

	if pool.Ref != nil && *pool.Ref != "" {
		refs.SetRef(*pool.Ref, *output.UserPool.Id)
	}

	for _, appClient := range *pool.AppClients {
		createAppClient(&appClient, *output.UserPool.Id, &ctx, client)
	}

	createUserPoolDomain(*output.UserPool.Id, pool.DomainPrefix, &ctx, client)

	console.Infof("Finished creating Cognito user pool %s", *pool.ApplicationName)
	return nil
}

func createAppClient(appClient *CognitoAppClient, userPoolId string, ctx *context.Context, client *awsCognito.Client) error {
	console.Verbosef("Creating cognito app client %s", appClient.Name)
	_, err := client.CreateUserPoolClient(*ctx, &awsCognito.CreateUserPoolClientInput{
		UserPoolId:     aws.String(userPoolId),
		ClientName:     aws.String(appClient.Name),
		GenerateSecret: false, // use true for server-side apps

		// 🔐 Required for OAuth
		AllowedOAuthFlowsUserPoolClient: true,
		AllowedOAuthFlows: []types.OAuthFlowType{
			types.OAuthFlowTypeCode, // or "implicit", "client_credentials"
		},
		AllowedOAuthScopes: []string{
			"openid", "email", "profile",
		},
		CallbackURLs: *appClient.ReturnUrls,
		LogoutURLs:   *appClient.LogoutUrls,
		SupportedIdentityProviders: []string{
			"COGNITO", // Required to allow signing in with the user pool itself
		},
	})

	if err != nil {
		return err
	}

	console.Verbosef("Finished creating cognito app client %s", appClient.Name)
	return nil
}

func createUserPoolDomain(userPoolId string, domainPrefix *string, ctx *context.Context, client *awsCognito.Client) {
	console.Verbose("Creating domain")

	prefix := userPoolId
	if domainPrefix != nil && *domainPrefix != "" {
		prefix = *domainPrefix
	}

	_, err := client.CreateUserPoolDomain(*ctx, &awsCognito.CreateUserPoolDomainInput{
		Domain:     aws.String(prefix),
		UserPoolId: aws.String(userPoolId),
	})

	if err != nil {
		console.Warnf("failed to create user pool domain: %s", err.Error())
		return
	}

	console.Verbose("Finished creating domain")
}
