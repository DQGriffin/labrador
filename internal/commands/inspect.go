package commands

import (
	"fmt"
	"sort"

	"github.com/DQGriffin/labrador/internal/cli/styles"
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
	default:
		printPlainText(config, stageTypesMap, verbose)
	}
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

		if stage.Type == "lambda" {
			for _, fnConfig := range stage.Functions {
				for _, fn := range fnConfig.Functions {
					childNodes := generateLambdaNodes(&fn, verbose)
					for i := range childNodes {
						node.Child(childNodes[i])
					}
				}
			}
		} else if stage.Type == "s3" {
			for _, s3Config := range stage.Buckets {
				for _, bucket := range s3Config.Buckets {
					childNodes := generateS3Nodes(&bucket, verbose)
					for i := range childNodes {
						node.Child(childNodes[i])
					}
				}
			}
		} else if stage.Type == "api" {
			for _, gatewayConfig := range stage.Gateways {
				for _, gateway := range gatewayConfig.Gateways {
					childNodes := generateApiGatewayNodes(&gateway, verbose)
					for i := range childNodes {
						node.Child(childNodes[i])
					}
				}
			}
		}

		nodes = append(nodes, node)
	}

	return nodes
}

func generateLambdaNodes(lambda *types.LambdaConfig, verbose bool) []*tree.Tree {
	var nodes []*tree.Tree
	node := tree.New().Root(styles.Tertiary.Bold(true).Render(lambda.Name))

	if verbose {
		node.Child(styles.Primary.Render("Region:     ") + styles.Secondary.Render(*lambda.Region))
		node.Child(styles.Primary.Render("Code:       ") + styles.Secondary.Render(*lambda.Code))
		node.Child(styles.Primary.Render("Handler:    ") + styles.Secondary.Render(*lambda.Handler))
		node.Child(styles.Primary.Render("Runtime:    ") + styles.Secondary.Render(*lambda.Runtime))
		node.Child(styles.Primary.Render("Role ARN:   ") + styles.Secondary.Render(*lambda.RoleArn))
		node.Child(styles.Primary.Render("Memory:     ") + styles.Secondary.Render(fmt.Sprintf("%dmb", *lambda.MemorySize)))
		node.Child(styles.Primary.Render("Timeout:    ") + styles.Secondary.Render(fmt.Sprintf("%ds", *lambda.Timeout)))
		node.Child(styles.Primary.Render("On Delete:  ") + styles.Secondary.Render(helpers.PtrOrDefault(lambda.OnDelete, "delete")))
	}

	nodes = append(nodes, node)

	return nodes
}

func generateS3Nodes(s3 *types.S3Settings, verbose bool) []*tree.Tree {
	var nodes []*tree.Tree
	node := tree.New().Root(styles.Tertiary.Bold(true).Render(*s3.Name))

	if verbose {
		node.Child(styles.Primary.Render("Region:               ") + styles.Secondary.Render(helpers.PtrOrDefault(s3.Region, "[region not set]")))
		node.Child(styles.Primary.Render("Versioning:           ") + styles.Secondary.Render(fmt.Sprintf("%t", helpers.PtrOrDefault(s3.Versioning, false))))
		node.Child(styles.Primary.Render("Block Public Access:  ") + styles.Secondary.Render(fmt.Sprintf("%t", helpers.PtrOrDefault(s3.BlockPublicAccess, true))))
		node.Child(styles.Primary.Render("On Delete:            ") + styles.Secondary.Render(helpers.PtrOrDefault(s3.OnDelete, "delete")))
	}

	nodes = append(nodes, node)
	return nodes
}

func generateApiGatewayNodes(gateway *types.ApiGatewaySettings, verbose bool) []*tree.Tree {
	var nodes []*tree.Tree
	node := tree.New().Root(styles.Tertiary.Bold(true).Render(*gateway.Name))

	if verbose {
		node.Child(styles.Primary.Render("Region:       ") + styles.Secondary.Render(helpers.PtrOrDefault(gateway.Region, "[region not set]")))
		node.Child(styles.Primary.Render("Protocol:     ") + styles.Secondary.Render(helpers.PtrOrDefault(gateway.Protocol, "[protocol not set]")))
		node.Child(styles.Primary.Render("Description:  ") + styles.Secondary.Render(helpers.PtrOrDefault(gateway.Description, "[description not set]")))
		stagesNode := tree.New().Root(styles.Primary.Render("Stages"))

		for _, stage := range *gateway.Stages {
			childNodes := generateApiGatewayStageNodes(&stage)
			for i := range childNodes {
				stagesNode.Child(childNodes[i])
			}
		}

		integrationsNode := tree.New().Root(styles.Primary.Render("Integrations"))
		for integationIndex, integration := range gateway.Integrations {
			childNodes := generateApiGatewayIntegrationNodes(&integration, integationIndex)
			for i := range childNodes {
				integrationsNode.Child(childNodes[i])
			}
		}

		routesNode := tree.New().Root(styles.Primary.Render("Routes"))
		for _, route := range gateway.Routes {
			childNodes := generateApiGatewayRouteNodes(&route)
			for i := range childNodes {
				routesNode.Child(childNodes[i])
			}
		}

		node.Child(stagesNode)
		node.Child(integrationsNode)
		node.Child(routesNode)
	}

	nodes = append(nodes, node)
	return nodes
}

func generateApiGatewayStageNodes(stage *types.ApiGatewayStage) []*tree.Tree {
	var nodes []*tree.Tree
	node := tree.New().Root(styles.Primary.Render(stage.Name))
	node.Child(styles.Primary.Render("Description:  ") + styles.Secondary.Render(helpers.PtrOrDefault(&stage.Description, "[description not set]")))
	node.Child(styles.Primary.Render("Auto-Deploy:  ") + styles.Secondary.Render(fmt.Sprintf("%t", stage.AutoDeploy)))

	nodes = append(nodes, node)
	return nodes
}

func generateApiGatewayIntegrationNodes(integration *types.ApiGatewayIntegration, iteration int) []*tree.Tree {
	var nodes []*tree.Tree
	node := tree.New().Root(styles.Primary.Render(fmt.Sprintf("Integration %d", iteration+1)))

	m := make(map[string]string)
	arn, err := aws.ResolveTarget(integration.Target, m)

	node.Child(styles.Primary.Render("Type:                ") + styles.Secondary.Render(integration.Type))
	if err != nil {
		node.Child(styles.Primary.Render("Target:              ") + styles.Secondary.Render("[unresolved]"))
	} else {
		node.Child(styles.Primary.Render("Target:              ") + styles.Secondary.Render(arn))
	}
	node.Child(styles.Primary.Render("Payload Version:     ") + styles.Secondary.Render(integration.PayloadVersion))
	node.Child(styles.Primary.Render("Integration Method:  ") + styles.Secondary.Render(integration.IntegrationMethod))

	nodes = append(nodes, node)
	return nodes
}

func generateApiGatewayRouteNodes(route *types.ApiGatewayRoute) []*tree.Tree {
	var nodes []*tree.Tree
	node := tree.New().Root(styles.Primary.Render(route.Method + " " + route.Route))
	node.Child(styles.Primary.Render("Target:  ") + styles.Secondary.Render(*route.Target.Ref))

	nodes = append(nodes, node)
	return nodes
}

func printPlainText(config *types.LabradorConfig, stageTypesMap *map[string]bool, verbose bool) {
	fmt.Println("============================")
	plainPrintProject(&config.Project, verbose)
	plainPrintStages(&config.Project.Stages, stageTypesMap, verbose)

	if !verbose {
		fmt.Println("\nRun with --verbose to view detailed resource configuration.")
	}
}

func plainPrintProject(project *types.Project, verbose bool) {
	fmt.Printf("Project: %s\n", project.Name)
	fmt.Printf("Environment: %s\n", project.Environment)
}

func plainPrintStages(stages *[]types.Stage, stageTypesMap *map[string]bool, verbose bool) {
	fmt.Println("\nStages:")
	for _, stage := range *stages {
		if isStageActionable(&stage, stageTypesMap) {
			fmt.Printf("- %s (%s)\n", stage.Name, stage.Type)

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
	fmt.Printf("  - %-25s -> %s\n", lambda.Name, *lambda.Code)
	if verbose {
		fmt.Printf("    - Region      : %s\n", *lambda.Region)
		fmt.Printf("    - Handler     : %s\n", *lambda.Handler)
		fmt.Printf("    - Runtime     : %s\n", *lambda.Runtime)
		fmt.Printf("    - Role ARN    : %s\n", *lambda.RoleArn)
		fmt.Printf("    - Memory      : %dmb\n", *lambda.MemorySize)
		fmt.Printf("    - Timeout     : %ds\n", *lambda.Timeout)
		fmt.Printf("    - On Delete   : %s\n", helpers.PtrOrDefault(lambda.OnDelete, "delete"))
		fmt.Println("    - Environment :")
		PrintMapAligned("      - ", lambda.Environment)
		fmt.Println("    - Tags :")
		PrintMapAligned("      - ", lambda.Tags)
		fmt.Println()
	}
}

func plainPrintS3(s3 *types.S3Settings, verbose bool) {
	fmt.Printf("  - %s \n", helpers.PtrOrDefault(s3.Name, "[Name not set]"))
	if verbose {
		fmt.Printf("    - Region               : %s\n", helpers.PtrOrDefault(s3.Region, "[region not set]"))
		fmt.Printf("    - Versioning           : %t\n", helpers.PtrOrDefault(s3.Versioning, false))
		fmt.Printf("    - Block Public Access  : %t\n", helpers.PtrOrDefault(s3.BlockPublicAccess, true))
		fmt.Printf("    - On Delete            : %s\n", helpers.PtrOrDefault(s3.OnDelete, "delete"))
		fmt.Println("    - Tags                 :")
		PrintMapAligned("      - ", s3.Tags)
		fmt.Println()
	}
}

func plainPrintApiGateway(gateway *types.ApiGatewaySettings, verbose bool) {
	fmt.Printf("  - %s \n", helpers.PtrOrDefault(gateway.Name, "[Name not set]"))
	if verbose {
		fmt.Printf("    - Protocol     : %s\n", helpers.PtrOrDefault(gateway.Protocol, "[protocol not set]"))
		fmt.Printf("    - Description  : %s\n", helpers.PtrOrDefault(gateway.Description, "[description not set]"))
		plainPrintApiGatewayStages(gateway.Stages)
		plainPrintApiGatewayIntegrations(&gateway.Integrations)
		plainPrintApiGatewayRoutes(&gateway.Routes)
		fmt.Println("    - Tags         :")
		PrintMapAligned("      - ", gateway.Tags)
		fmt.Println()
	}
}

func plainPrintApiGatewayIntegrations(integrations *[]types.ApiGatewayIntegration) {
	fmt.Println("    - Integrations")
	for _, integration := range *integrations {
		fmt.Printf("      - Type                : %s\n", integration.Type)
		m := make(map[string]string)
		arn, err := aws.ResolveTarget(integration.Target, m)
		if err != nil {
			fmt.Printf("      - Target              : %s\n", "[unresolved]")
		} else {
			fmt.Printf("      - Target              : %s\n", arn)
		}
		fmt.Printf("      - Payload version     : %s\n", integration.PayloadVersion)
		fmt.Printf("      - Integration method  : %s\n", integration.IntegrationMethod)
	}
}

func plainPrintApiGatewayRoutes(routes *[]types.ApiGatewayRoute) {
	fmt.Println("    - Routes")
	for _, route := range *routes {
		fmt.Printf("      - Method  : %s\n", route.Method)
		fmt.Printf("      - Route   : %s\n", route.Route)
		fmt.Printf("      - Target  : %s\n", *route.Target.Ref)
	}
}

func plainPrintApiGatewayStages(stages *[]types.ApiGatewayStage) {
	fmt.Println("    - Stages")
	for _, stage := range *stages {
		fmt.Printf("      - Name  : %s\n", stage.Name)
		fmt.Printf("      - Description   : %s\n", stage.Description)
		fmt.Printf("      - Auto-deploy   : %t\n", stage.AutoDeploy)
		fmt.Println("      - Tags        :")
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
	format := fmt.Sprintf("%s%%-%ds : %%s\n", prefix, maxKeyLen)

	// Print each key-value pair
	for _, key := range keys {
		fmt.Printf(format, key, m[key])
	}
}
