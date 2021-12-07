package consumer

import (
	"net/url"

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
)

func initWebdavConsumer(config map[string]interface{}) (stream.Consumer, error) {
	c, err := processConfig(util.NewConfigMap(config))
	if err != nil {
		return nil, err
	}
	return &webdavConsumer{config: c}, nil
}

type webdavConsumer struct {
	config webdavConfig
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
			for range channel {
			}
		}()
	} else {
		go func() {
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

type webdavConfig struct {
	connURL  *url.URL
	username string
	password string
}

func processConfig(config *util.ConfigMap) (webdavConfig, error) {
	scheme := config.Get(webdavSchemeCName).StringWithDefault(webdavDefaultScheme)
	host := config.Get(webdavHostCName).Required().String()
	path := config.Get(webdavPathCName).Required().String()
	username := config.Get(webdavUsernameCName).String()
	password := config.Get(webdavPasswordCName).String()
	if err := config.Error(); err != nil {
		return webdavConfig{}, err
	}

	u := &url.URL{
		Scheme: scheme,
		Host:   host,
		Path:   path,
	}

	return webdavConfig{
		connURL:  u,
		username: username,
		password: password,
	}, nil
}
