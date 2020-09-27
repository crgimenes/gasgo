package metadata

import (
	"github.com/google/uuid"
)

// Dataset define the base of the data type
type Dataset interface {
	Type() string
	String() string
}

// Loader define metadata loader interface
type Loader interface {
	LoadMetadata() Metadata
}

// Metadata define data structure
type Metadata struct {
	ID         uuid.UUID         `json:"id,omitempty"`
	Name       string            `json:"name"`
	Type       string            `json:"type"`
	Parameters map[string]string `json:"parameters,omitempty"`
	NestedMeta []Metadata        `json:"nested_meta,omitempty"`
}
