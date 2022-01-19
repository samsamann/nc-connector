package consumer

import (
	"net/url"

	"github.com/samsamann/nc-connector/internal/config"
	"github.com/samsamann/nc-connector/internal/stream"
	"github.com/samsamann/nc-connector/internal/stream/util"
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

func (w webdavConsumer) In() chan<- stream.SyncItem {
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
			for file := range channel {
				if file.Mode() == stream.WRITE {
					c.WriteStream(file.Path(), file.Data(), 0)
				} else if file.Mode() == stream.DELETE {
					c.Remove(file.Path())
				}
			}
			c = nil
		}()
	}
	return channel
}

func (w webdavConsumer) Wait() <-chan interface{} {
	return w.waitChan
}
type webdavConfig struct {
	connURL   *url.URL
	username  string
	password  string
	cachePath string
}

func processConfig(config config.NCClientConfig, opConfig map[string]interface{}) (webdavConfig, error) {
	configMap := util.NewConfigMap(opConfig)
	cachePath := configMap.Get(webdavCachePath).String()
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
	}, nil
}
