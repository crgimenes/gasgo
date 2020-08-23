package database

import (
	"github.com/gasgo/gasgo/metadata"
	"github.com/google/uuid"
)

// Database interface defines the functions required for the database manipulation.
type Database interface {
	Save(data []metadata.Dataset,
		m metadata.Metadata,
		parent uuid.UUID) (uuid.UUID, []metadata.Dataset, error) // insert or update
	List(where string) ([][]metadata.Dataset, error) // return multiple data
	Count(where string) (int, error)                 // return the number of metadata.Datasets
	Get(id uuid.UUID) ([]metadata.Dataset, error)    // return onli one item
	Delete(id uuid.UUID) error                       // delete metadata.Dataset
	String() string                                  // return database data in strigable form
	Close() error                                    // close database
}
