package sync

import (
	"github.com/samsamann/nc-connector/internal/stream"
	"github.com/studio-b12/gowebdav"
)

type Manager interface {
	IsNewer(stream.SyncItem) bool
	Save() error
}

type memManager struct {
	writer writer
	client *gowebdav.Client
	data   SearchableStorage
}

func NewInMemoryManager(loader Loader, client *gowebdav.Client) (Manager, error) {
	data, err := loader.Load()
	if err != nil {
		return nil, err
	}
	manager := new(memManager)
	manager.writer = loader
	manager.client = client
	manager.data = data

	return manager, nil
}

func (c memManager) Save() error {
	return c.writer.Unload(c.data)
}

func (c memManager) IsNewer(item stream.SyncItem) bool {
	entity := c.data.Get(item.Path())
	if entity != nil {
		fileInfo, err := c.client.Stat(item.Path())
		if err == nil {
			switch fi := fileInfo.(type) {
			case *gowebdav.File:
				if fi.ETag() == entity.eTag() ||
					fi.ModTime().Before(entity.modifiedDate()) {
					entity.markAsDurable()
					return false
				}
				c.data.Add(item.Path(), convert(fi.Name(), fi.ETag(), fi.ModTime()))
			}
		}
	}
	return true
}
