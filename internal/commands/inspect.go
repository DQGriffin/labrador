package commands

import (
	"fmt"
	"sort"

	"github.com/DQGriffin/labrador/internal/helpers"
	"github.com/DQGriffin/labrador/pkg/refs"
	"github.com/DQGriffin/labrador/pkg/types"
)

func HandleInspectCommand(config *types.LabradorConfig, format string, stageTypesMap *map[string]bool, verbose bool) {
	switch format {
	case "plain":
		printPlainText(config, stageTypesMap, verbose)
	default:
		printPlainText(config, stageTypesMap, verbose)
	}
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
		arn, err := refs.ResolveTarget(integration.Target, m)
		if err != nil {
			fmt.Printf("      - Target              : %s\n", "[unresolved]")
		}
		fmt.Printf("      - Target              : %s\n", arn)
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
