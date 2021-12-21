package operator

import (
	"net/url"

	"github.com/samsamann/nc-connector/internal/config"
	"github.com/samsamann/nc-connector/internal/stream"
	"github.com/samsamann/nc-connector/internal/stream/util"
	"github.com/samsamann/nc-connector/pkg/sync"
	"github.com/studio-b12/gowebdav"
)

const ItemNotChangedFilterName = "not-changed"

const (
	cacheFilePathCName = "cachePath"
)

func initNotChangedFilter(global *config.GlobalConfig, opConfig map[string]interface{}) (stream.Operator, error) {
	return newNotChangedFilter(
			global,
			util.NewConfigMap(opConfig)),
		nil
}

type notChangedFilter struct {
	channel chan stream.SyncItem
	manager sync.Manager
}

func newNotChangedFilter(globalConfig *config.GlobalConfig, cMap *util.ConfigMap) stream.Operator {
	cachePath := cMap.Get(cacheFilePathCName).Required().String()
	if cMap.Error() != nil {
		return nil
	}

	client := gowebdav.NewClient(
		(&url.URL{Scheme: "https", Host: globalConfig.NCClient.Host, Path: globalConfig.NCClient.BasePath}).String(),
		globalConfig.NCClient.Username,
		globalConfig.NCClient.Password,
	)

	m, _ := sync.NewInMemoryManager(sync.NewJsonFileLoader(cachePath), client)
	return &notChangedFilter{
		channel: make(chan stream.SyncItem),
		manager: m,
	}
}

func (f *notChangedFilter) In() chan<- stream.SyncItem {
	return f.channel
}

func (f *notChangedFilter) Out() <-chan stream.SyncItem {
	channel := make(chan stream.SyncItem)
	go func() {
		defer func() {
			close(channel)
			f.manager.Save()
		}()
		//f.cache.refreshCache()
		for item := range f.channel {
			if f.manager.IsNewer(item) {
				channel <- item
			}
		}

	}()

	return channel
}
