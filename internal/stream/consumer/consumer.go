package consumer

import (
	"fmt"

	"github.com/samsamann/nc-connector/internal/stream"
)

type InitConsumerFunc func(map[string]interface{}) (stream.Consumer, error)

var consumerRegistry map[string]InitConsumerFunc

func init() {
	consumerRegistry = make(map[string]InitConsumerFunc)
	consumerRegistry[stubConsumerName] = initStubConsumer
}

func CreateConsumer(name string, config map[string]interface{}) (stream.Consumer, error) {
	if f, ok := consumerRegistry[name]; ok {
		return f(config)
	}
	return nil, fmt.Errorf("no consumer found with name %q", name)
}
