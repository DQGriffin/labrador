package types

import (
	"encoding/json"
	"fmt"
)

type Project struct {
	Name        string            `json:"name"`
	Environment string            `json:"environment"`
	Stages      []Stage           `json:"stages"`
	Variables   map[string]string `json:"variables,omitempty" ,interpolate:"false"`
}

type Stage struct {
	Name         string   `json:"name"`
	Type         string   `json:"type"`
	Enabled      bool     `json:"enabled"`
	OnConflict   string   `json:"onConflict"`
	OnError      string   `json:"onError"`
	ConfigFile   string   `json:"config"`
	DependsOn    []string `json:"dependsOn"`
	Environments []string `json:"environments"`
	Hooks        Hooks    `json:"hooks"`
	Functions    []LambdaData
	Buckets      []S3Config
}

type Hooks struct {
	Pre  string `json:"pre"`
	Post string `json:"post"`
}

func (e Project) String() string {
	b, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return fmt.Sprintf("TelephonyEvent<error: %v>", err)
	}
	return string(b)
}
