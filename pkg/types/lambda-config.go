package types

type LambdaData struct {
	Defaults  *LambdaDefaults `json:"defaults,omitempty"`
	Functions []LambdaConfig  `json:"functions"`
}

type LambdaDefaults struct {
	Region      *string           `json:"region,omitempty"`
	RoleArn     *string           `json:"roleArn,omitempty"`
	Handler     *string           `json:"handler,omitempty"`
	Runtime     *string           `json:"runtime,omitempty"`
	Code        *string           `json:"code,omitempty"`
	MemorySize  *uint16           `json:"memory"`
	Timeout     *uint16           `json:"timeout"`
	Description *string           `json:"description"`
	Tags        map[string]string `json:"tags,omitempty"`
	Environment map[string]string `json:"environment,omitempty"`
}

type LambdaConfig struct {
	Name        string            `json:"name"`
	Region      *string           `json:"region,omitempty"`
	RoleArn     *string           `json:"roleArn,omitempty"`
	Handler     *string           `json:"handler,omitempty"`
	Runtime     *string           `json:"runtime,omitempty"`
	Code        *string           `json:"code,omitempty"`
	MemorySize  *uint16           `json:"memory"`
	Timeout     *uint16           `json:"timeout"`
	Description *string           `json:"description"`
	OnDelete    *string           `json:"onDelete"`
	Tags        map[string]string `json:"tags,omitempty"`
	Environment map[string]string `json:"environment,omitempty"`
}
