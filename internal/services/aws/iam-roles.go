package aws

import (
	"context"
	"encoding/json"
	"os"
	"strings"

	"github.com/DQGriffin/labrador/internal/cli/console"
	"github.com/DQGriffin/labrador/internal/cli/styles"
	"github.com/DQGriffin/labrador/internal/helpers"
	"github.com/DQGriffin/labrador/internal/refs"
	internalTypes "github.com/DQGriffin/labrador/internal/types"
	"github.com/DQGriffin/labrador/pkg/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamTypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
)

func generateAssumedRole(trustPolicy types.IamTrustPolicy) (string, error) {
	assumedrole := internalTypes.IamAssumedRolePolicy{
		Version: "2012-10-17",
		Statement: []internalTypes.IamAssumedRoleStatement{
			{
				Effect: "Allow",
				Principal: internalTypes.IamAssumedRolePrincipal{
					Service: trustPolicy.Principals.Services,
				},
				Action: "sts:AssumeRole",
			},
		},
	}

	data, err := json.MarshalIndent(assumedrole, "", "\t")
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func generateInlinePolicyDocument(policy types.IamInlinePolicy) (string, error) {
	inlinePolicy := internalTypes.IamInlinePolicy{
		Version: "2012-10-17",
		Statement: []internalTypes.IamInlinePolicyStatement{
			{
				Effect:   *policy.Effect,
				Action:   policy.Actions,
				Resource: policy.Resources,
			},
		},
	}

	data, err := json.MarshalIndent(inlinePolicy, "", "\t")
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func getTagsForRole(role *types.IamRoleSettings) []iamTypes.Tag {
	var tags []iamTypes.Tag

	for key, value := range role.Tags {
		tags = append(tags, iamTypes.Tag{Key: helpers.AsPtr(key), Value: helpers.AsPtr(value)})
	}

	return tags
}

func CreateIamRole(role *types.IamRoleSettings) error {
	console.Styledf(&styles.PrimaryStyle, "[%s]", *role.Name)
	console.Infof("Creating IAM role %s", *role.Name)

	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		console.Fatalf("Unable to load AWS config: %v", err)
	}

	client := iam.NewFromConfig(cfg)
	var document string

	if role.TrustPolicy.FilePath != nil && *role.TrustPolicy.FilePath != "" {
		console.Verbosef("Reading trust policy document from %s", *role.TrustPolicy.FilePath)
		documentData, docErr := os.ReadFile(*role.TrustPolicy.FilePath)
		if docErr != nil {
			console.Debug("Failed to read trust policy document")
			return docErr
		}

		document = string(documentData)
	} else {
		policyDocument, trustPolicyErr := generateAssumedRole(*role.TrustPolicy)
		if trustPolicyErr != nil {
			console.Debug("Failed to generate trust policy document")
			return trustPolicyErr
		}

		document = policyDocument
	}

	output, err := client.CreateRole(ctx, &iam.CreateRoleInput{
		RoleName:                 aws.String(*role.Name),
		AssumeRolePolicyDocument: aws.String(document),
		Description:              aws.String(*role.Description),
		Tags:                     getTagsForRole(role),
	})
	if err != nil {
		if strings.Contains(err.Error(), "409") {
			console.Infof("IAM role %s already exists\n", *role.Name)
			return nil
		}
		console.Errorf("failed to create IAM role: %v", err.Error())
	}

	for _, policy := range role.PolicyArns {
		err := attachPolicy(*role.Name, policy, &ctx, client)
		if err != nil {
			return err
		}
	}

	for _, inlinePolicy := range role.InlinePolicies {
		err := attachInlinePolicy(*role.Name, &inlinePolicy, &ctx, client)
		if err != nil {
			return err
		}
	}

	console.Debugf("Role ARN: %s", *output.Role.Arn)
	if role.Ref != nil && *role.Ref != "" {
		refs.SetRef(*role.Ref, *output.Role.Arn)
	}
	console.Infof("Finished creating IAM role %s\n", *role.Name)
	return nil
}

func attachPolicy(roleName, policyArn string, ctx *context.Context, client *iam.Client) error {
	console.Verbosef("Attaching policy ARN to IAM role %s, %s", roleName, policyArn)

	_, err := client.AttachRolePolicy(*ctx, &iam.AttachRolePolicyInput{
		RoleName:  aws.String(roleName),
		PolicyArn: aws.String(policyArn),
	})
	if err != nil {
		return err
	}

	console.Verbosef("Finished attaching managed policy to IAM role %s", roleName)
	return nil
}

func attachInlinePolicy(roleName string, policy *types.IamInlinePolicy, ctx *context.Context, client *iam.Client) error {
	console.Verbosef("Attaching inline policy %s to role %s", policy.Name, roleName)
	var document string

	if policy.FilePath != nil && *policy.FilePath != "" {
		console.Debugf("Reading inline policy document from %s", *policy.FilePath)
		documentData, docErr := os.ReadFile(*policy.FilePath)
		if docErr != nil {
			console.Debug("Failed to read inline policy document")
			return docErr
		}

		document = string(documentData)
		console.Debug(document)
	} else {
		generatedDoc, docErr := generateInlinePolicyDocument(*policy)
		if docErr != nil {
			return docErr
		}
		document = generatedDoc
	}

	_, err := client.PutRolePolicy(*ctx, &iam.PutRolePolicyInput{
		RoleName:       aws.String(roleName),
		PolicyName:     aws.String(policy.Name),
		PolicyDocument: aws.String(document),
	})
	if err != nil {
		return err
	}

	console.Verbosef("Finished adding inline policy %s to role %s", policy.Name, roleName)
	return nil
}
