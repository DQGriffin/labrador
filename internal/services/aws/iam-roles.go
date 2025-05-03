package aws

import (
	"context"
	"encoding/json"
	"fmt"
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

func DeleteRole(roleName string) error {
	console.Infof("Deleting IAM role %s", roleName)
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		console.Fatalf("Unable to load AWS config: %v", err)
	}

	client := iam.NewFromConfig(cfg)

	policyErr := detachAllManagedPolicies(ctx, client, roleName)
	if policyErr != nil {
		return policyErr
	}

	inlinePolicyErr := DeleteAllInlinePolicies(ctx, client, roleName)
	if inlinePolicyErr != nil {
		return inlinePolicyErr
	}

	_, deleteErr := client.DeleteRole(ctx, &iam.DeleteRoleInput{
		RoleName: &roleName,
	})

	if deleteErr != nil {
		if strings.Contains(deleteErr.Error(), "404") {
			console.Infof("IAM role %s did not exist. No action taken", roleName)
			return nil
		}

		return deleteErr
	}

	console.Infof("Finished deleting IAM role %s", roleName)
	return nil
}

func detachAllManagedPolicies(ctx context.Context, client *iam.Client, roleName string) error {
	console.Verbosef("Detatching managed policies from %s", roleName)
	var marker *string

	for {
		output, err := client.ListAttachedRolePolicies(ctx, &iam.ListAttachedRolePoliciesInput{
			RoleName: aws.String(roleName),
			Marker:   marker,
		})
		if err != nil {
			if strings.Contains(err.Error(), "404") {
				console.Verbose("Role not found. No managed policies to detatch")
				return nil
			}
			return fmt.Errorf("listing attached policies: %w", err)
		}

		for _, policy := range output.AttachedPolicies {
			_, err := client.DetachRolePolicy(ctx, &iam.DetachRolePolicyInput{
				RoleName:  aws.String(roleName),
				PolicyArn: policy.PolicyArn,
			})
			if err != nil {
				return fmt.Errorf("detaching policy %s: %w", aws.ToString(policy.PolicyArn), err)
			}
		}

		if !output.IsTruncated {
			break
		}
		marker = output.Marker
	}

	console.Verbosef("Finished detatching managed policies from %s", roleName)
	return nil
}

func DeleteAllInlinePolicies(ctx context.Context, client *iam.Client, roleName string) error {
	console.Verbosef("Deleting inline policies from IAM role %s", roleName)
	var marker *string

	for {
		output, err := client.ListRolePolicies(ctx, &iam.ListRolePoliciesInput{
			RoleName: aws.String(roleName),
			Marker:   marker,
		})
		if err != nil {
			if strings.Contains(err.Error(), "404") {
				console.Verbose("Role not found. No managed policies to detatch")
				return nil
			}
			return fmt.Errorf("listing inline policies: %w", err)
		}

		for _, policyName := range output.PolicyNames {
			_, err := client.DeleteRolePolicy(ctx, &iam.DeleteRolePolicyInput{
				RoleName:   aws.String(roleName),
				PolicyName: aws.String(policyName),
			})
			if err != nil {
				return fmt.Errorf("deleting inline policy %s: %w", policyName, err)
			}
		}

		if !output.IsTruncated {
			break
		}
		marker = output.Marker
	}

	console.Verbosef("Finished deleting inline policies from IAM role %s", roleName)
	return nil
}

func ListAllRoleNames() ([]string, error) {
	var (
		roleNames []string
		marker    *string
	)

	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		console.Fatalf("Unable to load AWS config: %v", err)
	}

	client := iam.NewFromConfig(cfg)

	for {
		output, err := client.ListRoles(ctx, &iam.ListRolesInput{
			Marker: marker,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list roles: %w", err)
		}

		for _, role := range output.Roles {
			roleNames = append(roleNames, aws.ToString(role.RoleName))
		}

		if !output.IsTruncated {
			break
		}

		marker = output.Marker
	}

	return roleNames, nil
}
