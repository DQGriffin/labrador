package types

type UniversalResourceDefinition struct {
	StageName    string `json:"stageName"`
	Name         string `json:"name"`
	Arn          string `json:"arn"`
	ResourceType string `json:"resourceType"`
	Region       string `json:"region"`
}
