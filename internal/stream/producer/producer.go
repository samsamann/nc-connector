package producer

import (
	"fmt"

	"github.com/samsamann/nc-connector/internal/stream"
)

type InitProducerFunc func(map[string]interface{}) (stream.Producer, error)

var producerRegistry map[string]InitProducerFunc

func init() {
	producerRegistry = make(map[string]InitProducerFunc)
	producerRegistry[fsProducerName] = initFsProducer
	producerRegistry[mssqlProducerName] = initMssqlProducer
}

func CreateProducer(name string, config map[string]interface{}) (stream.Producer, error) {
	if f, ok := producerRegistry[name]; ok {
		return f(config)
	}
	return nil, fmt.Errorf("no producer found with name %q", name)
}
