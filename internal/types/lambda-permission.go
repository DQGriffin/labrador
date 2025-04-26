package types

type LambdaPermission struct {
	Action       string
	FunctionName string
	Principal    string
	StatementId  string
	SourceArn    string
}
