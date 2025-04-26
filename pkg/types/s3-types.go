package types

type S3Config struct {
	Defaults *S3Settings  `json:"defaults"`
	Buckets  []S3Settings `json:"buckets"`
}

type S3Settings struct {
	Name              *string                `json:"name,omitempty"`
	Region            *string                `json:"region,omitempty"`
	Versioning        *bool                  `json:"versioning,omitempty"`
	OnDelete          *string                `json:"onDelete,omitempty"`
	BlockPublicAccess *bool                  `json:"blockPublicAccess,omitempty"`
	StaticHosting     *StaticHostingSettings `json:"staticHosting,omitempty"`
	Tags              map[string]string      `json:"tags,omitempty"`
}

type StaticHostingSettings struct {
	Enabled       bool    `json:"enabled"`
	IndexDocument *string `json:"indexDocument,omitempty"`
	ErrorDocument *string `json:"errorDocument,omitempty"`
}
