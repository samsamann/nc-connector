package producer

import (
	"fmt"

	"github.com/samsamann/nc-connector/internal/stream"
)

type InitProducerFunc func() (stream.Producer, error)

var producerRegistry map[string]InitProducerFunc

func init() {
	producerRegistry = make(map[string]InitProducerFunc)
	producerRegistry[mssqlProducer] = initMssqlProducer
}

func CreateProducer(name string) (stream.Producer, error) {
	if f, ok := producerRegistry[name]; ok {
		return f()
	}
	return nil, fmt.Errorf("no producer found with name %q", name)
}
