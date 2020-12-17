package pipeline

// FileData represents the transfer object that was loaded from an inmport source.
type FileData interface {
}

// File represents a file.
type File struct {
	Name string
}

// Folder represents a folder.
type Folder struct {
	Name string
}
