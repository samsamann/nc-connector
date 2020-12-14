package pipeline

import "errors"

// CreateFileImporterFunc returns a new instance of a specific FileImporter.
type CreateFileImporterFunc func(map[string]interface{}) FileImporter

var (
	register map[string]CreateFileImporterFunc

	errNotFound = errors.New("No element found")
)

// RegisterFileImporter registers a FileImporter.
func RegisterFileImporter(name string, nFunc CreateFileImporterFunc) {
	if register == nil {
		register = make(map[string]CreateFileImporterFunc)
	}
	register[name] = nFunc
}

// GetFileImporter returns a new FileImporter instance with the specified name.
func GetFileImporter(name string, config map[string]interface{}) (FileImporter, error) {
	if register == nil {
		return nil, errNotFound
	}
	if nFunc, ok := register[name]; ok {
		return nFunc(config), nil
	}
	return nil, errNotFound
}
