package types

// I didn't put much though into the naming here

type LabradorConfig struct {
	Project      Project
	FunctionData []LambdaData
}

type ResourceTarget struct {
	Ref      *string            `json:"ref,omitempty"`
	External *ExternalReference `json:"external,omitempty"`
}

type ExternalReference struct {
	Arn     *string                 `json:"arn,omitempty"`
	Dynamic *DynamicResourceRefData `json:"dynamic"`
}

type DynamicResourceRefData struct {
	Name   string `json:"name"`
	Region string `json:"region"`
	Type   string `json:"type"`
}
