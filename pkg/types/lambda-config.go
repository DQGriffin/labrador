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
	MemorySize  *uint16           `json:"memory,omitempty"`
	Timeout     *uint16           `json:"timeout,omitempty"`
	Description *string           `json:"description,omitempty"`
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
	MemorySize  *uint16           `json:"memory,omitempty"`
	Timeout     *uint16           `json:"timeout,omitempty"`
	Description *string           `json:"description,omitempty"`
	OnDelete    *string           `json:"onDelete,omitempty"`
	Tags        map[string]string `json:"tags,omitempty"`
	Environment map[string]string `json:"environment,omitempty"`
}
