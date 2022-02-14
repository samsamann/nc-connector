package operator

import (
	"path"
	"strings"
	"text/template"

	"github.com/samsamann/nc-connector/internal/config"
	"github.com/samsamann/nc-connector/internal/stream"
	"github.com/samsamann/nc-connector/internal/stream/util"
)

const pathManipulatorName = "path"

const (
	pathAppendOpCName = "append"
)

func initPathManipulator(global *config.GlobalConfig, config map[string]interface{}) (stream.Operator, error) {
	cMap := util.NewConfigMap(config)
	appendTemp := cMap.Get(pathAppendOpCName).Required().String()
	if err := cMap.Error(); err != nil {
		return nil, err
	}
	t, err := template.New("path").Parse(appendTemp)
	if err != nil {
		return nil, err
	}
	return newPathManipulator(t), nil
}

type pathManipulator struct {
	channel  chan stream.SyncItem
	template *template.Template
}

func newPathManipulator(t *template.Template) *pathManipulator {
	return &pathManipulator{
		channel:  make(chan stream.SyncItem),
		template: t,
	}
}

func (ms pathManipulator) In(ctx stream.Context) chan<- stream.SyncItem {
	return ms.channel
}

func (ms pathManipulator) Out(ctx stream.Context) <-chan stream.SyncItem {
	channel := make(chan stream.SyncItem)
	go func() {
		defer close(channel)
		for file := range ms.channel {
			manipulatePath(file, ms.template)
			channel <- file
		}
	}()
	return channel
}

func manipulatePath(file stream.SyncItem, t *template.Template) {
	var builder strings.Builder
	if err := t.Execute(&builder, file.Attributes()); err != nil {
		return
	}
	dir, fileName := path.Split(file.Path())
	file.SetPath(
		path.Join(dir, "/", builder.String(), "/", fileName),
	)
}
