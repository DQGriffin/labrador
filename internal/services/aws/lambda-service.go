package aws

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/DQGriffin/labrador/internal/cli/console"
	internalTypes "github.com/DQGriffin/labrador/internal/types"
	"github.com/DQGriffin/labrador/pkg/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	lambdaTypes "github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

const (
	region = "us-east-1"
)

func ListLambdas() (map[string]lambdaTypes.FunctionConfiguration, error) {
	m := make(map[string]lambdaTypes.FunctionConfiguration)

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(os.Getenv("AWS_REGION")),
	)
	if err != nil {
		return m, err
	}

	client := lambda.NewFromConfig(cfg)

	paginator := lambda.NewListFunctionsPaginator(client, &lambda.ListFunctionsInput{})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			console.Fatal(err)
		}
		for _, fn := range page.Functions {
			m[*fn.FunctionName] = lambdaTypes.FunctionConfiguration{}
		}
	}

	return m, nil
}

// Should refactor this in the future. Currently we're creating a new client every time
// a function is created or update. Ideally we would reuse the client
func CreateLambda(lambdaConfig types.LambdaConfig) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(*lambdaConfig.Region))
	if err != nil {
		console.Fatalf("Unable to load AWS config: %v", err)
	}

	client := lambda.NewFromConfig(cfg)

	zipData, err := os.ReadFile(*lambdaConfig.Code)
	if err != nil {
		console.Fatalf("Failed to read scaffold zip: %v", err)
	}

	_, getErr := client.GetFunction(context.TODO(), &lambda.GetFunctionInput{
		FunctionName: aws.String(lambdaConfig.Name),
	})

	if getErr == nil {
		console.Infof("Lambda %q already exists. Skipping.", lambdaConfig.Name)
		// continue
		return
	}

	console.Infof("Creating Lambda %q...", lambdaConfig.Name)
	_, err = client.CreateFunction(context.TODO(), &lambda.CreateFunctionInput{
		FunctionName: aws.String(lambdaConfig.Name),
		Description:  aws.String(*lambdaConfig.Description),
		Timeout:      aws.Int32(int32(*lambdaConfig.Timeout)),
		Role:         aws.String(*lambdaConfig.RoleArn),
		Handler:      aws.String(*lambdaConfig.Handler),
		Runtime:      lambdaTypes.Runtime(*lambdaConfig.Runtime),
		MemorySize:   aws.Int32(int32(*lambdaConfig.MemorySize)),
		Code: &lambdaTypes.FunctionCode{
			ZipFile: zipData,
		},
		Environment: &lambdaTypes.Environment{
			Variables: lambdaConfig.Environment,
		},
		Tags:    lambdaConfig.Tags,
		Publish: true,
	})

	if err != nil {
		console.Errorf("Failed to create function %q: %v", lambdaConfig.Name, err)
		return
	}

	console.Infof("Created Lambda %q", lambdaConfig.Name)
}

func UpdateLambda(lambdaConfig types.LambdaConfig) {
	console.Infof("Updating lambda %q", lambdaConfig.Name)
	updateLambdaCode(lambdaConfig)
	time.Sleep(5 * time.Second)
	UpdateLambdaConfiguration(lambdaConfig)
	console.Infof("Finished updating lambda %q", lambdaConfig.Name)
}

func updateLambdaCode(lambdaConfig types.LambdaConfig) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(*lambdaConfig.Region))
	if err != nil {
		console.Fatalf("Unable to load AWS config: %v", err)
	}

	client := lambda.NewFromConfig(cfg)

	zipData, err := os.ReadFile(*lambdaConfig.Code)
	if err != nil {
		console.Fatalf("Failed to read zip: %v", err)
	}

	_, updateErr := client.UpdateFunctionCode(context.TODO(), &lambda.UpdateFunctionCodeInput{
		FunctionName: aws.String(lambdaConfig.Name),
		ZipFile:      zipData,
	})
	if updateErr != nil {
		console.Errorf("Failed to update function code for %s: %v", lambdaConfig.Name, err)
	}
}

func UpdateLambdaConfiguration(lambdaConfig types.LambdaConfig) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(*lambdaConfig.Region))
	if err != nil {
		console.Fatalf("Unable to load AWS config: %v", err)
	}

	client := lambda.NewFromConfig(cfg)

	_, err = client.UpdateFunctionConfiguration(context.TODO(), &lambda.UpdateFunctionConfigurationInput{
		FunctionName: aws.String(lambdaConfig.Name),
		Handler:      aws.String(*lambdaConfig.Handler),
		Runtime:      lambdaTypes.Runtime(*lambdaConfig.Runtime),
		MemorySize:   aws.Int32(int32(*lambdaConfig.MemorySize)),
		Timeout:      aws.Int32(int32(*lambdaConfig.Timeout)),
		Role:         aws.String(*lambdaConfig.RoleArn),
		Environment: &lambdaTypes.Environment{
			Variables: lambdaConfig.Environment,
		},
	})
	if err != nil {
		console.Errorf("Failed to update function config: %v", err)
	}
}

func GetLambda(ctx context.Context, cfg aws.Config, lambdaName string) (lambdaTypes.FunctionConfiguration, error) {
	client := lambda.NewFromConfig(cfg)

	output, err := client.GetFunction(ctx, &lambda.GetFunctionInput{
		FunctionName: aws.String(lambdaName),
	})

	if err != nil {
		return lambdaTypes.FunctionConfiguration{}, err
	}

	fn := *output.Configuration

	return fn, err
}

func DeleteLambda(lambdaName string) {
	console.Infof("Deleting lambda: %s", lambdaName)
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		console.Fatalf("Unable to load AWS config: %v", err)
	}

	client := lambda.NewFromConfig(cfg)

	_, deleteErr := client.DeleteFunction(context.TODO(), &lambda.DeleteFunctionInput{
		FunctionName: aws.String(lambdaName),
	})
	if deleteErr != nil {
		if strings.Contains(deleteErr.Error(), "404") {
			console.Infof("Lambda %s did not exist. No action taken", lambdaName)
			return
		} else {
			console.Errorf("failed to delete Lambda %s: %s", lambdaName, deleteErr.Error())
			return
		}
	}

	console.Infof("Deleted Lambda: %s", lambdaName)
}

func AddPermissionToLambda(ctx context.Context, cfg aws.Config, permission internalTypes.LambdaPermission) error {
	client := lambda.NewFromConfig(cfg)

	_, err := client.AddPermission(ctx, &lambda.AddPermissionInput{
		Action:       aws.String(permission.Action),
		FunctionName: aws.String(permission.FunctionName),
		Principal:    aws.String(permission.Principal),
		StatementId:  aws.String(permission.StatementId),
		SourceArn:    aws.String(permission.SourceArn),
	})
	if err != nil {
		return fmt.Errorf("failed to add permission to %s: %w", permission.FunctionName, err)
	}

	console.Infof("Added permission to lambda %s", permission.FunctionName)
	return nil
}
