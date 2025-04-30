package types

type ApiGatewayConfig struct {
	Defaults *ApiGatewaySettings  `json:"defaults"`
	Gateways []ApiGatewaySettings `json:"gateways"`
}

type ApiGatewaySettings struct {
	Name         *string                 `json:"name,omitempty"`
	OnDelete     *string                 `json:"onDelete,omitempty"`
	Description  *string                 `json:"description,omitempty"`
	Region       *string                 `json:"region,omitempty"`
	Protocol     *string                 `json:"protocol,omitempty"`
	Stages       *[]ApiGatewayStage      `json:"stages,omitempty"`
	Integrations []ApiGatewayIntegration `json:"integrations,omitempty"`
	Routes       []ApiGatewayRoute       `json:"routes,omitempty"`
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
