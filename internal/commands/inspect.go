package commands

import (
	"fmt"
	"sort"

	"github.com/DQGriffin/labrador/internal/helpers"
	"github.com/DQGriffin/labrador/pkg/types"
)

func HandleInspectCommand(config *types.LabradorConfig, format string, verbose bool) {
	switch format {
	case "plain":
		printPlainText(config, verbose)
	default:
		printPlainText(config, verbose)
	}
}

func printPlainText(config *types.LabradorConfig, verbose bool) {
	fmt.Println("============================")
	plainPrintProject(&config.Project, verbose)
	plainPrintStages(&config.Project.Stages, verbose)

	if !verbose {
		fmt.Println("\nInspect more with --verbose")
	}
}

func plainPrintProject(project *types.Project, verbose bool) {
	fmt.Printf("Project: %s\n", project.Name)
	fmt.Printf("Environment: %s\n", project.Environment)
}

func plainPrintStages(stages *[]types.Stage, verbose bool) {
	fmt.Println("\nStages:")
	for _, stage := range *stages {
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

	}
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
