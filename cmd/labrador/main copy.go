package main

// import (
// 	"fmt"
// 	"log"
// 	"os"

// 	"github.com/DQGriffin/labrador/internal/aws/lambda"
// 	"github.com/DQGriffin/labrador/internal/validation"
// 	"github.com/DQGriffin/labrador/pkg/interpolation"
// 	"github.com/DQGriffin/labrador/pkg/utils"
// 	"github.com/joho/godotenv"
// )

// func main() {
// 	fmt.Println("Hello, Labrador!")

// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}

// 	project, err := utils.ReadProjectData("/Users/griffin/dev/labrador/templates/project.json")
// 	if err != nil {
// 		fmt.Println("Failed to read project config")
// 		os.Exit(1)
// 	}

// 	errs := validation.ValidateProject(project)
// 	if len(errs) > 0 {
// 		fmt.Println("Errors validating project config")
// 		for _, err := range errs {
// 			fmt.Println(err)
// 		}
// 		os.Exit(1)
// 	}

// 	fmt.Println("Successfully read project config")
// 	interpolation.InterpolateProjectVariables(&project)
// 	project.Variables["project_name"] = project.Name
// 	project.Variables["env"] = project.Environment

// 	functionData, readErr := utils.ReadFunctionConfigs(project.Stages)

// 	if readErr != nil {
// 		fmt.Println(readErr)
// 		os.Exit(1)
// 	}

// 	for i := range functionData {
// 		interpolation.Interpolate(&functionData[i], project.Variables)

// 		for functionIndex := range functionData[i].Functions {
// 			project.Variables["name"] = functionData[i].Functions[functionIndex].Name
// 			interpolation.Interpolate(&functionData[i].Functions[functionIndex], project.Variables)
// 		}
// 	}

// 	utils.ApplyDefaultsToFunctions(&functionData[0])

// 	lambda.CreateLambda(functionData[0].Functions[0])
// }
