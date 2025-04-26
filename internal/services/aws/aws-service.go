package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

func GetConfig(awsRegion string) (context.Context, aws.Config, error) {
	ctx := context.TODO()

	// Load AWS config (credentials + region)
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(awsRegion))
	if err != nil {
		return ctx, cfg, err
	}

	return ctx, cfg, nil
}
