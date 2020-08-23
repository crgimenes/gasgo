package metadata

import "github.com/google/uuid"

// Dataset define the base of the data type
type Dataset interface {
	Type() string
	String() string
}

// Metadata define data structure
type Metadata struct {
	ID         uuid.UUID
	Name       string
	Type       string
	Unique     bool
	Parameters map[string]string
	NestedMeta []Metadata
}
