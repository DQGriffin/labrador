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
	IamRoles     []IamRoleConfig    `json:"-"`
}

type Hooks struct {
	WorkingDir     string   `json:"workingDir,omitempty"`
	SuppressStdout bool     `json:"suppressStdout,omitempty"`
	SuppressStderr bool     `json:"suppressStderr,omitempty"`
	StopOnError    bool     `json:"stopOnError,omitempty"`
	PreDeploy      []string `json:"preDeploy,omitempty"`
	PostDeploy     []string `json:"postDeploy,omitempty"`
	PreDestroy     []string `json:"preDestroy,omitempty"`
	PostDestroy    []string `json:"postDestroy,omitempty"`
}

func (e Project) String() string {
	b, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return fmt.Sprintf("TelephonyEvent<error: %v>", err)
	}
	return string(b)
}
