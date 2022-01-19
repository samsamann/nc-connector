package consumer

import (
	"fmt"

	"github.com/samsamann/nc-connector/internal/config"
	"github.com/samsamann/nc-connector/internal/stream"
)

type InitConsumerFunc func(*config.GlobalConfig, map[string]interface{}) (stream.Consumer, error)

var consumerRegistry map[string]InitConsumerFunc

func init() {
	consumerRegistry = make(map[string]InitConsumerFunc)
	consumerRegistry[stubConsumerName] = initStubConsumer
	consumerRegistry[webdavConsumerName] = initWebdavConsumer
}

func CreateConsumer(opConfig config.StreamElem, c *config.GlobalConfig) (stream.Consumer, error) {
	name := opConfig.Name
	if f, ok := consumerRegistry[name]; ok {
		return f(c, opConfig.Config)
	}
	return nil, fmt.Errorf("no consumer found with name %q", name)
}
