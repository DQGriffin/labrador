package types

type AwsLambda struct {
	Name        string `json:"name"`
	Runtime     string `json:"runtime"`
	Description string `json:"description"`
	FunctionArn string `json:"functionArn"`
	Handler     string `json:"handler"`
}
