package database

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

// Database interface defines the functions required for the database manipulation.
type Database interface {
	Save(data []Dataset, m Metadata, parent uuid.UUID) (uuid.UUID, []Dataset, error) // insert or update
	List(where string) ([][]Dataset, error)                                          // return multiple data
	Count(where string) (int, error)                                                 // return the number of datasets
	Get(id uuid.UUID) ([]Dataset, error)                                             // return onli one item
	Delete(id uuid.UUID) error                                                       // delete dataset
	String() string                                                                  // return database data in strigable form
	Close() error                                                                    // close database
}
