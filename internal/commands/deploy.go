package commands

import (
	"fmt"

	"github.com/DQGriffin/labrador/internal/services/aws"
	"github.com/DQGriffin/labrador/pkg/types"
	lambdaTypes "github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

func HandleDeployCommand(config types.LabradorConfig, existingLambdas map[string]lambdaTypes.FunctionConfiguration, onlyCreate bool, onlyUpdate bool) {
	for _, stage := range config.Project.Stages {
		deployLambdaStage(&stage, existingLambdas, onlyCreate, onlyUpdate)
	}
}

func deployLambdaStage(stage *types.Stage, existingLambdas map[string]lambdaTypes.FunctionConfiguration, onlyCreate bool, onlyUpdate bool) {
	fmt.Printf("Deploying Stage: %s\n", stage.Name)
	fmt.Printf("Type: %s\n", stage.Type)

	for _, fnConfig := range stage.Functions {
		for _, fn := range fnConfig.Functions {
			if _, exists := existingLambdas[fn.Name]; exists {
				if !onlyCreate {
					fmt.Println("updating function", fn.Name)
					aws.UpdateLambda(fn)
				}
			} else {
				if !onlyUpdate {
					fmt.Println("creating function", fn.Name)
					aws.CreateLambda(fn)
				}
			}
		}
	}
}
