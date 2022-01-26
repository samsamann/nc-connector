package sync

import (
	"time"

	"github.com/samsamann/nc-connector/internal/stream"
	"github.com/studio-b12/gowebdav"
)

type Manager interface {
	IsNewer(stream.SyncItem) bool
	Add(string, string, time.Time)
	Delete(stream.SyncItem)
	RemovableItems() []stream.SyncItem
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
	fileInfo, err := c.client.Stat(item.Path())
	if err != nil {
		return true
	}
	fi := fileInfo.(*gowebdav.File)
	entity := c.data.Get(item.Path())
	if entity != nil && (fi.ETag() == entity.eTag() &&
		fi.ModTime().Equal(entity.modifiedDate())) {
		entity.markAsDurable()
		return false
	}

	return true
}

func (c memManager) Add(path, etag string, modDate time.Time) {
	c.data.Add(path, convert(path, etag, modDate))
}

func (c memManager) Delete(item stream.SyncItem) {
	c.data.Delete(item.Path())
}

func (c memManager) RemovableItems() []stream.SyncItem {
	items := make([]stream.SyncItem, 0)
	for _, e := range c.data.removable() {
		item := stream.NewFileSyncItem(e.path(), make(stream.Properties), nil)
		item.ChangeMode(stream.DELETE)
		items = append(items, item)
	}
	return items
}
