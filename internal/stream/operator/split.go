package operator

import (
	"strings"

	"github.com/samsamann/nc-connector/internal/config"
	"github.com/samsamann/nc-connector/internal/stream"
	"github.com/samsamann/nc-connector/internal/stream/util"
)

const splitOperatorrName = "split"

const (
	splitAttrNameCName    = "attr"
	splitSeperatorCName   = "sep"
	splitNewAttrNameCName = "newAttr"
)

func initSplitOperator(global *config.GlobalConfig, config map[string]interface{}) (stream.Operator, error) {
	cMap := util.NewConfigMap(config)
	attrName := cMap.Get(splitAttrNameCName).Required().String()
	sep := cMap.Get(splitSeperatorCName).Required().String()
	newAttrName := cMap.Get(splitNewAttrNameCName).String()
	if err := cMap.Error(); err != nil {
		return nil, err
	}
	return newSplitOperator(attrName, sep, newAttrName), nil
}

type splitOperator struct {
	channel     chan stream.SyncItem
	attrName    string
	newAttrName string
	sep         string
}

func newSplitOperator(attrName, sep, newAttrName string) *splitOperator {
	return &splitOperator{
		channel:     make(chan stream.SyncItem),
		attrName:    attrName,
		newAttrName: newAttrName,
		sep:         sep,
	}
}

func (ms splitOperator) In(ctx stream.Context) chan<- stream.SyncItem {
	return ms.channel
}

func (ms splitOperator) Out(ctx stream.Context) <-chan stream.SyncItem {
	channel := make(chan stream.SyncItem)
	go func() {
		defer close(channel)
		for file := range ms.channel {
			attrs := file.Attributes()
			for _, v := range split(attrs, ms.attrName, ms.sep) {
				copyAttrs := make(stream.Properties, len(attrs))
				for k, v := range attrs {
					copyAttrs[k] = v
				}
				if len(ms.newAttrName) > 0 {
					copyAttrs[ms.newAttrName] = v
				} else {
					copyAttrs[ms.attrName] = v
				}
				switch t := file.(type) {
				case *stream.File:
					concreteFile := *t
					concreteFile.SetAttributes(copyAttrs)
					channel <- &concreteFile
				}
			}
		}
	}()
	return channel
}

func split(attrs stream.Properties, name, sep string) []string {
	attr, ok := attrs[name]
	if ok {
		if str, ok := attr.(string); ok {
			return strings.Split(str, sep)
		}
	}
	return []string{""}
}
