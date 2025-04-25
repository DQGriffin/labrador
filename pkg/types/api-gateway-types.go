package types

type ApiGatewayConfig struct {
	Defaults *ApiGatewaySettings  `json:"defaults"`
	Gateways []ApiGatewaySettings `json:"gateways"`
}

type ApiGatewaySettings struct {
	Name         *string                 `json:"name"`
	OnDelete     *string                 `json:"onDelete"`
	Description  *string                 `json:"description"`
	Region       *string                 `json:"region"`
	Protocol     *string                 `json:"protocol"`
	Stages       *[]ApiGatewayStage      `json:"stages,omitempty"`
	Integrations []ApiGatewayIntegration `json:"integrations"`
	Routes       []ApiGatewayRoute       `json:"routes"`
	Tags         map[string]string       `json:"tags,omitempty"`
}

type ApiGatewayIntegration struct {
	Type              string         `json:"type"`
	PayloadVersion    string         `json:"payloadVersion"`
	IntegrationMethod string         `json:"integrationMethod"`
	Ref               string         `json:"ref"`
	Target            ResourceTarget `json:"target"`
}

type ApiGatewayRoute struct {
	Method string         `json:"method"`
	Route  string         `json:"route"`
	Target ResourceTarget `json:"target"`
}

type ApiGatewayStage struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	AutoDeploy  bool              `json:"automaticDeployment"`
	Tags        map[string]string `json:"tags,omitempty"`
}
