package refs

import (
	"errors"
	"fmt"

	"github.com/DQGriffin/labrador/internal/services/aws"
	"github.com/DQGriffin/labrador/pkg/types"
)

func ResolveTarget(target types.ResourceTarget, refMap map[string]string) (string, error) {
	// Prefer Labrador-managed reference
	if target.Ref != nil && *target.Ref != "" {
		return resolveRef(*target.Ref, refMap)
	}

	// Use explicit ARN if present
	if target.External != nil {
		if target.External.Arn != nil && *target.External.Arn != "" {
			return *target.External.Arn, nil
		}

		// Try dynamic lookup
		if target.External.Dynamic != nil {
			return lookupByNameAndType(*target.External.Dynamic)
		}
	}

	return "", errors.New("no valid target found")
}

func resolveRef(ref string, refMap map[string]string) (string, error) {
	result := refMap[ref]

	if result == "" {
		return "", errors.New("ref not found")
	}

	return result, nil
}

func lookupByNameAndType(resource types.DynamicResourceRefData) (string, error) {
	if resource.Type == "s3" {
		arn := lookupS3Bucket(resource.Name)
		return arn, nil
	} else if resource.Type == "lambda" {
		arn, err := lookupLambda(resource)
		return arn, err
	}
	return "", errors.New("lookupByNameAndType is not fully implemented")
}

// We're not actually going to lookup anything here.
// S3 ARNs are deterministic
func lookupS3Bucket(bucketName string) string {
	return fmt.Sprintf("arn:aws:s3:::%s", bucketName)
}

func lookupLambda(resource types.DynamicResourceRefData) (string, error) {
	ctx, cfg, err := aws.GetConfig(resource.Region)
	if err != nil {
		return "", nil
	}

	fn, fnErr := aws.GetLambda(ctx, cfg, resource.Name)
	if fnErr != nil {
		return "", fnErr
	}

	return *fn.FunctionArn, nil
}
