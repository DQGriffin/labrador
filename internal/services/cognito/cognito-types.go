package cognito

type CognitoConfig struct {
	Defaults *CognitoSettings  `json:"defaults,omitempty"`
	Pools    []CognitoSettings `json:"pools,omitempty"`
}

type CognitoSettings struct {
	ApplicationName      *string                      `json:"applicationName,omitempty"`
	Ref                  *string                      `json:"ref,omitempty"`
	DomainPrefix         *string                      `json:"domainPrefix,omitempty"`
	SignInIdentifiers    *[]string                    `json:"signInIdentifiers,omitempty"`
	SignUpAttributes     *[]string                    `json:"signUpAttributes,omitempty"`
	PasswordRequirements *CognitoPasswordRequirements `json:"passwordPolicy,omitempty"`
	AppClients           *[]CognitoAppClient          `json:"appClients,omitempty"`
	Tags                 map[string]string            `json:"tags,omitempty"`
}

type CognitoPasswordRequirements struct {
	MinLength        int  `json:"minLength"`
	RequireSymbols   bool `json:"requireSymbols"`
	RequireNumbers   bool `json:"requireNumbers"`
	RequireUppercase bool `json:"requireUppercase"`
	RequireLowercase bool `json:"requireLowercase"`
}

type CognitoAppClient struct {
	Name       string    `json:"name,omitempty"`
	ClientType string    `json:"type,omitempty"`
	ReturnUrls *[]string `json:"returnUrls,omitempty"`
	LogoutUrls *[]string `json:"logoutUrls,omitempty"`
}
