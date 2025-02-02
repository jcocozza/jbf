package dal

import "github.com/jcocozza/jbf/internal/metadata"

type Repository interface {
	CreateTag(metadataID int, name string) error
	ReadTagExists(tagName string) bool
	ReadTags(metadataID int) ([]string, error)
	DeleteTag(tagName string) error

	CreateMetadata(m metadata.Metadata) (int, error)
	ReadMetadataExists(filepath string) bool
	ReadMetadata(filepath string) (metadata.Metadata, error)
	ReadMetadataFiles() ([]string, error)
	ReadAllMetadata() ([]metadata.Metadata, error)
	UpdateMetadata(m metadata.Metadata) error
	DeleteMetadata(filepath string) error
}
