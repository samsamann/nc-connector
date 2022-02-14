package consumer

import (
	"context"
	"fmt"
	"net/url"
	"runtime"
	wait "sync"

	"github.com/samsamann/nc-connector/internal/config"
	"github.com/samsamann/nc-connector/internal/stream"
	"github.com/samsamann/nc-connector/internal/stream/util"
	"github.com/samsamann/nc-connector/pkg/sync"
	"github.com/studio-b12/gowebdav"
)

const (
	webdavConsumerName  = "webdav"
	webdavDefaultScheme = "https"
)

const (
	webdavSchemeCName   = "scheme"
	webdavHostCName     = "host"
	webdavPathCName     = "basePath"
	webdavUsernameCName = "username"
	webdavPasswordCName = "password"
	webdavWorkers       = "workers"
	webdavCachePath     = "cachePath"
)

func initWebdavConsumer(global *config.GlobalConfig, opConfig map[string]interface{}) (stream.Consumer, error) {
	c, err := processConfig(global.NCClient, opConfig)
	if err != nil {
		return nil, err
	}
	return &webdavConsumer{config: c, waitChan: make(chan interface{})}, nil
}

type webdavConsumer struct {
	config   webdavConfig
	waitChan chan interface{}
}

func (w webdavConsumer) In(ctx stream.Context) chan<- stream.SyncItem {
	channel := make(chan stream.SyncItem)

	c := gowebdav.NewClient(
		w.config.connURL.String(),
		w.config.username,
		w.config.password,
	)
	if err := c.Connect(); err != nil {
		go func() {
			defer func() {
				w.waitChan <- nil
			}()
			for range channel {
			}
		}()
	} else {
		go func() {
			defer func() {
				w.waitChan <- nil
			}()
			m, _ := sync.NewInMemoryManager(sync.NewJsonFileLoader(w.config.cachePath), c)
			ctx := context.TODO()
			var wg wait.WaitGroup
			for i := 0; i < w.config.workers; i++ {
				wg.Add(1)
				go worker(ctx, &wg, channel, c, m)
			}
			wg.Wait()

			removeChan := make(chan stream.SyncItem)
			go func() {
				for _, f := range m.RemovableItems() {
					removeChan <- f
				}
				close(removeChan)
			}()
			for i := 0; i < w.config.workers; i++ {
				wg.Add(1)
				go worker(ctx, &wg, removeChan, c, m)
			}
			wg.Wait()
			m.Save()
			c = nil
		}()
	}
	return channel
}

func (w webdavConsumer) Wait() <-chan interface{} {
	return w.waitChan
}

func worker(ctx context.Context, wg *wait.WaitGroup, data <-chan stream.SyncItem, c *gowebdav.Client, m sync.Manager) {
	defer wg.Done()
	for {
		select {
		case file, ok := <-data:
			if !ok {
				return
			}
			if file.Mode() == stream.WRITE {
				execWriteOperation(file, c, m)

			} else if file.Mode() == stream.DELETE {
				execDeleteOperation(file, c, m)
			}
		case <-ctx.Done():
			fmt.Printf("cancelled worker. Error detail: %v\n", ctx.Err())
			return
		}
	}
}

func execWriteOperation(file stream.SyncItem, client *gowebdav.Client, manager sync.Manager) {
	if manager.IsNewer(file) {
		err := client.WriteStream(file.Path(), file.Data(), 0)
		if err == nil {
			fileInfo, err := client.Stat(file.Path())
			if err != nil {
				return
			}
			fi := fileInfo.(*gowebdav.File)
			manager.Add(file.Path(), fi.ETag(), fi.ModTime())
		}
	}
}

func execDeleteOperation(file stream.SyncItem, client *gowebdav.Client, manager sync.Manager) {
	err := client.Remove(file.Path())
	if err == nil {
		manager.Delete(file)
	}
}

type webdavConfig struct {
	connURL   *url.URL
	username  string
	password  string
	cachePath string
	workers   int
}

func processConfig(config config.NCClientConfig, opConfig map[string]interface{}) (webdavConfig, error) {
	configMap := util.NewConfigMap(opConfig)
	cachePath := configMap.Get(webdavCachePath).String()
	workers := configMap.Get(webdavWorkers).IntWithDefault(runtime.NumCPU() * 2)
	/*host := config.Get(webdavHostCName).Required().String()
	path := config.Get(webdavPathCName).Required().String()
	username := config.Get(webdavUsernameCName).String()
	password := config.Get(webdavPasswordCName).String()
	if err := config.Error(); err != nil {
		return webdavConfig{}, err
	}*/

	u := &url.URL{
		Scheme: webdavDefaultScheme,
		Host:   config.Host,
		Path:   config.BasePath,
	}

	return webdavConfig{
		connURL:   u,
		username:  config.Username,
		password:  config.Password,
		cachePath: cachePath,
		workers:   workers,
	}, nil
}
