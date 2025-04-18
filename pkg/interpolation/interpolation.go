package interpolation

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/DQGriffin/labrador/pkg/types"
)

func InterpolateProjectVariables(project *types.Project) {
	project.Name = ResolveVariable(project.Name, project.Variables)
	project.Environment = ResolveVariable(project.Environment, project.Variables)

	for i := range project.Stages {
		InterpolateStage(&project.Stages[i], project.Variables)
	}

}

func ResolveVariable(value string, vars map[string]string) string {
	for k, v := range vars {
		value = strings.ReplaceAll(value, "{{"+k+"}}", v)
	}

	// Replace $VAR from the OS environment
	value = os.ExpandEnv(value)

	return value
}

func InterpolateStage(stage *types.Stage, vars map[string]string) {
	stage.Name = ResolveVariable(stage.Name, vars)
	stage.ConfigFile = ResolveVariable(stage.ConfigFile, vars)
	stage.OnConflict = ResolveVariable(stage.OnConflict, vars)
	stage.OnError = ResolveVariable(stage.OnError, vars)
	stage.Hooks.Pre = ResolveVariable(stage.Hooks.Pre, vars)
	stage.Hooks.Post = ResolveVariable(stage.Hooks.Post, vars)
}

func Interpolate(target any, vars map[string]string) error {
	v := reflect.ValueOf(target)
	if v.Kind() != reflect.Pointer || v.IsNil() {
		return fmt.Errorf("target must be a non-nil pointer")
	}
	return interpolateValue(v.Elem(), vars)
}

func interpolateValue(v reflect.Value, vars map[string]string) error {
	switch v.Kind() {
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)

			if !field.CanSet() {
				continue // skip unexported fields
			}
			if err := interpolateValue(field, vars); err != nil {
				return err
			}
		}

	case reflect.Ptr:
		if !v.IsNil() {
			return interpolateValue(v.Elem(), vars)
		}

	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			if err := interpolateValue(v.Index(i), vars); err != nil {
				return err
			}
		}

	case reflect.Map:
		if v.Type().Key().Kind() == reflect.String && v.Type().Elem().Kind() == reflect.String {
			for _, key := range v.MapKeys() {
				val := v.MapIndex(key)
				if val.Kind() == reflect.String {
					resolved := ResolveVariable(val.String(), vars)
					v.SetMapIndex(key, reflect.ValueOf(resolved))
				}
			}
		}

	case reflect.String:
		v.SetString(ResolveVariable(v.String(), vars))
	}

	return nil
}
