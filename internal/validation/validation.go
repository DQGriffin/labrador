package validation

import (
	"fmt"

	"github.com/DQGriffin/labrador/pkg/types"
)

func ValidateProject(project types.Project) []error {
	var errs []error
	for _, stage := range project.Stages {
		conflicResolutionError := validateConflictResolution(stage.OnConflict)
		if conflicResolutionError != nil {
			errs = append(errs, fmt.Errorf("stage %q: %w", stage.Name, conflicResolutionError))
		}

		errorResolutionError := validateErrorResolution(stage.OnError)
		if errorResolutionError != nil {
			errs = append(errs, fmt.Errorf("stage %q: %w", stage.Name, errorResolutionError))
		}

		stageTypeError := validateStageType(stage.Type)
		if stageTypeError != nil {
			errs = append(errs, fmt.Errorf("stage %q: %w", stage.Name, stageTypeError))
		}
	}

	return errs
}

func validateConflictResolution(value string) error {
	switch value {
	case "stop":
		return nil
	case "update":
		return nil
	default:
		return fmt.Errorf("onConflict must be one of: stop, update, skip")
	}
}

func validateErrorResolution(value string) error {
	switch value {
	case "stop":
		return nil
	case "skip":
		return nil
	case "rollback":
		return nil
	default:
		return fmt.Errorf("onError must be one of: stop, skip, rollback")
	}
}

// TODO: Refactor this
func validateStageType(value string) error {
	switch value {
	case "lambda":
		return nil
	case "s3":
		return nil
	case "api":
		return nil
	case "iam-role":
		return nil
	default:
		return fmt.Errorf("type must be one of: lambda, s3, api, iam-role")
	}
}

func ValidateFunctions(functionData types.LambdaData) []error {
	var errs []error

	// Not yet implemented

	return errs
}
