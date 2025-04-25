package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func GetAccountID() (string, error) {
	_, cfg, err := GetConfig("us-east-1")

	if err != nil {
		return "", err
	}

	client := sts.NewFromConfig(cfg)

	resp, err := client.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})
	if err != nil {
		return "", fmt.Errorf("failed to get caller identity: %w", err)
	}

	return *resp.Account, nil
}
