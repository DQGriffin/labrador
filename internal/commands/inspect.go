package commands

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/DQGriffin/labrador/internal/cli/console"
	"github.com/DQGriffin/labrador/internal/helpers"
	"github.com/DQGriffin/labrador/internal/services/aws"
	"github.com/DQGriffin/labrador/pkg/types"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
)

func HandleInspectCommand(config *types.LabradorConfig, format string, stageTypesMap *map[string]bool, verbose bool) {
	switch format {
	case "plain":
		printPlainText(config, stageTypesMap, verbose)
	case "tree":
		printTree(config, stageTypesMap, verbose)
	case "json":
		printJson(config)
	default:
		printPlainText(config, stageTypesMap, verbose)
	}

	if len(config.Project.Stages) == 0 {
		console.Info()
		console.Info("No stages defined yet.")
		console.Info("Run 'labrador add stage --help' to see available options")
	}
}

func printJson(config *types.LabradorConfig) {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		console.Error("Failed to marshal project config to JSON:", err)
		return
	}
	fmt.Println(string(data))
}

func printTree(config *types.LabradorConfig, stageTypesMap *map[string]bool, verbose bool) {
	rootStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	t := tree.Root(config.Project.Name).RootStyle(rootStyle)

	nodes := generateStageNodes(&config.Project.Stages, stageTypesMap, verbose)

	for i := range nodes {
		t.Child(nodes[i])
	}

	fmt.Println(t)
}

func generateStageNodes(stages *[]types.Stage, stageTypesMap *map[string]bool, verbose bool) []*tree.Tree {
	var nodes []*tree.Tree

	for _, stage := range *stages {
		if !helpers.IsStageActionable(&stage, stageTypesMap) {
			continue
		}

		node := tree.New().
			Root(stage.Name)

		childNodes := stage.ToTreeNodes(verbose)
		for i := range childNodes {
			node.Child(childNodes[i])
		}

		nodes = append(nodes, node)
	}

	return nodes
}

func printPlainText(config *types.LabradorConfig, stageTypesMap *map[string]bool, verbose bool) {
	console.Info("============================")
	plainPrintProject(&config.Project, verbose)
	plainPrintStages(&config.Project.Stages, stageTypesMap, verbose)

	if !verbose {
		console.Info("\nRun with --full to view detailed resource configuration.")
	}
}

func plainPrintProject(project *types.Project, verbose bool) {
	console.Infof("Project: %s", project.Name)
	console.Infof("Environment: %s", project.Environment)
}

func plainPrintStages(stages *[]types.Stage, stageTypesMap *map[string]bool, verbose bool) {
	console.Info("\nStages:")
	for _, stage := range *stages {
		if isStageActionable(&stage, stageTypesMap) {
			console.Infof("- %s (%s)", stage.Name, stage.Type)

			for _, fnConfig := range stage.Functions {
				for _, fn := range fnConfig.Functions {
					plainPrintLambda(&fn, verbose)
				}
			}

			for _, s3Config := range stage.Buckets {
				for _, bucket := range s3Config.Buckets {
					plainPrintS3(&bucket, verbose)
				}
			}

			for _, gatewayConfig := range stage.Gateways {
				for _, gateway := range gatewayConfig.Gateways {
					plainPrintApiGateway(&gateway, verbose)
				}
			}
		}
	}
}

func isStageActionable(stage *types.Stage, stageTypesMap *map[string]bool) bool {
	if len(*stageTypesMap) == 0 {
		return true
	}

	return (*stageTypesMap)[stage.Type]
}

func plainPrintLambda(lambda *types.LambdaConfig, verbose bool) {
	console.Infof("  - %-25s -> %s", lambda.Name, *lambda.Code)
	if verbose {
		console.Infof("    - Region      : %s", *lambda.Region)
		console.Infof("    - Handler     : %s", *lambda.Handler)
		console.Infof("    - Runtime     : %s", *lambda.Runtime)
		console.Infof("    - Role ARN    : %s", *lambda.RoleArn)
		console.Infof("    - Memory      : %dmb", *lambda.MemorySize)
		console.Infof("    - Timeout     : %ds", *lambda.Timeout)
		console.Infof("    - On Delete   : %s", helpers.PtrOrDefault(lambda.OnDelete, "delete"))
		console.Info("    - Environment :")
		PrintMapAligned("      - ", lambda.Environment)
		console.Info("    - Tags :")
		PrintMapAligned("      - ", lambda.Tags)
		console.Info()
	}
}

func plainPrintS3(s3 *types.S3Settings, verbose bool) {
	console.Infof("  - %s ", helpers.PtrOrDefault(s3.Name, "[Name not set]"))
	if verbose {
		console.Infof("    - Region               : %s", helpers.PtrOrDefault(s3.Region, "[region not set]"))
		console.Infof("    - Versioning           : %t", helpers.PtrOrDefault(s3.Versioning, false))
		console.Infof("    - Block Public Access  : %t", helpers.PtrOrDefault(s3.BlockPublicAccess, true))
		console.Infof("    - On Delete            : %s", helpers.PtrOrDefault(s3.OnDelete, "delete"))
		console.Info("    - Tags                 :")
		PrintMapAligned("      - ", s3.Tags)
		console.Info()
	}
}

func plainPrintApiGateway(gateway *types.ApiGatewaySettings, verbose bool) {
	console.Infof("  - %s ", helpers.PtrOrDefault(gateway.Name, "[Name not set]"))
	if verbose {
		console.Infof("    - Protocol     : %s", helpers.PtrOrDefault(gateway.Protocol, "[protocol not set]"))
		console.Infof("    - Description  : %s", helpers.PtrOrDefault(gateway.Description, "[description not set]"))
		plainPrintApiGatewayStages(gateway.Stages)
		plainPrintApiGatewayIntegrations(&gateway.Integrations)
		plainPrintApiGatewayRoutes(&gateway.Routes)
		console.Info("    - Tags         :")
		PrintMapAligned("      - ", gateway.Tags)
		console.Info()
	}
}

func plainPrintApiGatewayIntegrations(integrations *[]types.ApiGatewayIntegration) {
	console.Info("    - Integrations")
	for _, integration := range *integrations {
		console.Infof("      - Type                : %s", integration.Type)
		m := make(map[string]string)
		arn, err := aws.ResolveTarget(integration.Target, m)
		if err != nil {
			console.Infof("      - Target              : %s", "[unresolved]")
		} else {
			console.Infof("      - Target              : %s", arn)
		}
		console.Infof("      - Payload version     : %s", integration.PayloadVersion)
		console.Infof("      - Integration method  : %s", integration.IntegrationMethod)
	}
}

func plainPrintApiGatewayRoutes(routes *[]types.ApiGatewayRoute) {
	console.Info("    - Routes")
	for _, route := range *routes {
		console.Infof("      - Method  : %s", route.Method)
		console.Infof("      - Route   : %s", route.Route)
		console.Infof("      - Target  : %s", *route.Target.Ref)
	}
}

func plainPrintApiGatewayStages(stages *[]types.ApiGatewayStage) {
	console.Info("    - Stages")
	for _, stage := range *stages {
		console.Infof("      - Name  : %s", stage.Name)
		console.Infof("      - Description   : %s", stage.Description)
		console.Infof("      - Auto-deploy   : %t", stage.AutoDeploy)
		console.Info("      - Tags        :")
		PrintMapAligned("        - ", stage.Tags)
	}
}

func PrintMapAligned(prefix string, m map[string]string) {
	// First, find the longest key
	maxKeyLen := 0
	for key := range m {
		if len(key) > maxKeyLen {
			maxKeyLen = len(key)
		}
	}

	// Optional: print in sorted order
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Format string for aligned printing
	format := fmt.Sprintf("%s%%-%ds : %%s", prefix, maxKeyLen)

	// Print each key-value pair
	for _, key := range keys {
		console.Infof(format, key, m[key])
	}
}
