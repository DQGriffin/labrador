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
	Name         string             `json:"name"`
	Type         string             `json:"type"`
	Enabled      bool               `json:"enabled,omitempty"`
	OnConflict   string             `json:"onConflict"`
	OnError      string             `json:"onError"`
	ConfigFile   string             `json:"config"`
	DependsOn    []string           `json:"dependsOn,omitempty"`
	Environments []string           `json:"environments"`
	Hooks        *Hooks             `json:"hooks,omitempty"`
	Functions    []LambdaData       `json:"-"`
	Buckets      []S3Config         `json:"-"`
	Gateways     []ApiGatewayConfig `json:"-"`
}

type Hooks struct {
	Pre  string `json:"pre,omitempty"`
	Post string `json:"post,omitempty"`
}

func (e Project) String() string {
	b, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return fmt.Sprintf("TelephonyEvent<error: %v>", err)
	}
	return string(b)
}
