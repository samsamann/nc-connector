package groupfolder

import (
	"crypto"

	pip "github.com/samsamann/nc-connector/pkg/pipeline"
)

type groupfolderExporter struct{}

// NewGroupfolderExporter retutns a new instance of groupfolderexporter.
func NewGroupfolderExporter() pip.FileExporter {
	return new(groupfolderExporter)
}

func (g groupfolderExporter) Export(ctx pip.ExportContext, channel <-chan pip.FileData) error {
	var ok = true
	for ok {
		select {
		case <-ctx.Done():
			ok = false
		case fd, ok := <-channel:
			if !ok {
				break
			}
			switch fdType := fd.(type) {
			case *pip.File:
				exportFile(fdType)
			case *pip.Folder:
				exportFolder(fdType)
			default:
				// TODO: log warning
			}
		}
	}
	return nil
}

func exportFile(file *pip.File) {}

func exportFolder(folder *pip.Folder) {}

func hash(b []byte) []byte {
	hash := crypto.SHA256
	return hash.New().Sum(b)
}
