package producer

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/samsamann/nc-connector/internal/stream"
	"github.com/samsamann/nc-connector/internal/stream/util"
)

const fsProducerName = "fs"

const (
	fsPathCName       = "path"
	fsRemoveFileCName = "removeFile"
)

func initFsProducer(config map[string]interface{}) (stream.Producer, error) {
	c := util.NewConfigMap(config)
	path := c.Get(fsPathCName).Required().String()
	removeFile := c.Get(fsRemoveFileCName).Bool()
	if err := c.Error(); err != nil {
		return nil, err
	}
	return newFsProducer(path, removeFile), nil
}

type fsProducer struct {
	path       string
	removeFile bool
}

func newFsProducer(path string, removeFile bool) stream.Producer {
	p := new(fsProducer)
	p.path = path
	p.removeFile = removeFile
	return p
}

func (fs fsProducer) Out(ctx stream.Context) <-chan stream.SyncItem {
	channel := make(chan stream.SyncItem)

	go func() {
		defer func() {
			close(channel)
		}()
		fileInfos, err := ioutil.ReadDir(fs.path)
		if err != nil {
			// ctx.
			return
		}
		for _, fi := range fileInfos {
			if fi.IsDir() {
				continue
			}
			filepath := path.Join(fs.path, fi.Name())
			content, err := ioutil.ReadFile(filepath)
			if err != nil {
				continue
			}
			channel <- stream.NewFileSyncItem(fi.Name(), make(stream.Properties), content)
			if fs.removeFile {
				os.Remove(filepath)
			}
		}

	}()

	return channel
}
