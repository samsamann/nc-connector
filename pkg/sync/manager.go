package sync

import (
	"sync"
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
	mu     sync.Mutex
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

func (m *memManager) Save() error {
	return m.writer.Unload(m.data)
}

func (m *memManager) IsNewer(item stream.SyncItem) bool {
	fileInfo, err := m.client.Stat(item.Path())
	if err != nil {
		return true
	}
	fi := fileInfo.(*gowebdav.File)
	entity := m.data.Get(item.Path())
	if entity != nil && (fi.ETag() == entity.eTag() &&
		fi.ModTime().Equal(entity.modifiedDate())) {
		m.mu.Lock()
		entity.markAsDurable()
		m.mu.Unlock()
		return false
	}

	return true
}

func (m *memManager) Add(path, etag string, modDate time.Time) {
	m.mu.Lock()
	m.data.Add(path, convert(path, etag, modDate))
	m.mu.Unlock()
}

func (m *memManager) Delete(item stream.SyncItem) {
	m.mu.Lock()
	m.data.Delete(item.Path())
	m.mu.Unlock()
}

func (m *memManager) RemovableItems() []stream.SyncItem {
	defer m.mu.Unlock()
	m.mu.Lock()
	items := make([]stream.SyncItem, 0)
	for _, e := range m.data.removable() {
		item := stream.NewFileSyncItem(e.path(), make(stream.Properties), nil)
		item.ChangeMode(stream.DELETE)
		items = append(items, item)
	}
	return items
}
