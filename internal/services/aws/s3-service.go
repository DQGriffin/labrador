package aws

import (
	"context"
	"fmt"

	"github.com/DQGriffin/labrador/pkg/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3Types "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func GetClient(cfg aws.Config) *s3.Client {
	client := s3.NewFromConfig(cfg)
	return client
}

func ListBuckets(ctx context.Context, client *s3.Client) (map[string]bool, error) {
	// A map isn't the best way to do this, but it'll work for now
	m := make(map[string]bool)

	output, err := client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return m, err
	}

	for _, bucket := range output.Buckets {
		m[*bucket.Name] = true
	}

	return m, nil
}

func CreateBucket(ctx context.Context, cfg aws.Config, client s3.Client, bucket types.S3Settings) error {
	input := &s3.CreateBucketInput{
		Bucket: aws.String(*bucket.Name),
	}

	// S3 quirk: us-east-1 cannot have LocationConstraint set
	if *bucket.Region != "us-east-1" {
		input.CreateBucketConfiguration = &s3Types.CreateBucketConfiguration{
			LocationConstraint: s3Types.BucketLocationConstraint(*bucket.Region),
		}
	}

	_, err := client.CreateBucket(ctx, input)
	if err != nil {
		return err
	}

	settingsErr := setBucketSettings(ctx, client, &bucket)
	if settingsErr != nil {
		return settingsErr
	}

	fmt.Printf("Created bucket: %s\n", *bucket.Name)

	return nil
}

func UpdateBucket(ctx context.Context, client s3.Client, bucket types.S3Settings) error {
	settingsErr := setBucketSettings(ctx, client, &bucket)
	if settingsErr != nil {
		return settingsErr
	}

	fmt.Printf("Updated bucket: %s\n", *bucket.Name)

	return nil
}

func DeleteBucket(bucketName string, force bool) error {
	ctx, cfg, err := GetConfig("us-east-2")

	if err != nil {
		return err
	}

	client := GetClient(cfg)

	if force {
		// Region is hard coded for now. Need to refactor
		EmptyBucket(ctx, bucketName, "us-east-2")
	}

	_, deleteErr := client.DeleteBucket(ctx, &s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	})

	if deleteErr != nil {
		return fmt.Errorf("failed to delete bucket %s: %w", bucketName, deleteErr)
	}

	fmt.Printf("Deleted bucket: %s\n", bucketName)
	return nil
}

func EmptyBucket(ctx context.Context, bucketName, region string) error {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := s3.NewFromConfig(cfg)

	paginator := s3.NewListObjectsV2Paginator(client, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("failed to list objects: %w", err)
		}

		if len(page.Contents) == 0 {
			break
		}

		// Prepare delete request for up to 1000 objects
		var objectsToDelete []s3Types.ObjectIdentifier
		for _, obj := range page.Contents {
			objectsToDelete = append(objectsToDelete, s3Types.ObjectIdentifier{
				Key: obj.Key,
			})
		}

		_, err = client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
			Bucket: aws.String(bucketName),
			Delete: &s3Types.Delete{
				Objects: objectsToDelete,
				Quiet:   aws.Bool(true),
			},
		})
		if err != nil {
			return fmt.Errorf("failed to delete objects: %w", err)
		}

		fmt.Printf("Deleted %d objects from %s\n", len(objectsToDelete), bucketName)
	}

	fmt.Printf("Bucket %s is now empty\n", bucketName)
	return nil
}

func setBucketSettings(ctx context.Context, client s3.Client, bucket *types.S3Settings) error {
	tagsErr := setTags(ctx, client, bucket)
	if tagsErr != nil {
		return tagsErr
	}

	versioningErr := setVersioning(ctx, client, *bucket)
	if versioningErr != nil {
		return versioningErr
	}

	publicAccessErr := blockPublicAccess(ctx, client, bucket)
	if publicAccessErr != nil {
		return publicAccessErr
	}

	return nil
}

func setTags(ctx context.Context, client s3.Client, bucket *types.S3Settings) error {
	_, err := client.PutBucketTagging(ctx, &s3.PutBucketTaggingInput{
		Bucket: aws.String(*bucket.Name),
		Tagging: &s3Types.Tagging{
			TagSet: ConvertTags(bucket.Tags),
		},
	})

	if err != nil {
		return err
	}

	return nil
}

func setVersioning(ctx context.Context, client s3.Client, bucket types.S3Settings) error {

	_, err := client.PutBucketVersioning(ctx, &s3.PutBucketVersioningInput{
		Bucket: aws.String(*bucket.Name),
		VersioningConfiguration: &s3Types.VersioningConfiguration{
			Status: s3Types.BucketVersioningStatus(getVersioningSettingString(*bucket.Versioning)),
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func blockPublicAccess(ctx context.Context, client s3.Client, bucket *types.S3Settings) error {
	_, err := client.PutPublicAccessBlock(ctx, &s3.PutPublicAccessBlockInput{
		Bucket: aws.String(*bucket.Name),
		PublicAccessBlockConfiguration: &s3Types.PublicAccessBlockConfiguration{
			BlockPublicAcls:       aws.Bool(*bucket.BlockPublicAccess),
			IgnorePublicAcls:      aws.Bool(*bucket.BlockPublicAccess),
			BlockPublicPolicy:     aws.Bool(*bucket.BlockPublicAccess),
			RestrictPublicBuckets: aws.Bool(*bucket.BlockPublicAccess),
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func getVersioningSettingString(versioningEnabled bool) string {
	switch versioningEnabled {
	case true:
		return "Enabled"
	default:
		return "Disabled"
	}
}

func ConvertTags(tagMap map[string]string) []s3Types.Tag {
	tags := make([]s3Types.Tag, 0, len(tagMap))
	for k, v := range tagMap {
		tags = append(tags, s3Types.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		})
	}
	return tags
}
