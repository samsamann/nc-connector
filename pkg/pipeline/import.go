package pipeline

// FileImporter is the interface that wraps the Import method.
type FileImporter interface {
	Import() <-chan File
}
