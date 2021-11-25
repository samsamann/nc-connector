package stream

import (
	"fmt"
	"testing"
	"time"
)

type stubConsumer struct{}

func (s stubConsumer) In() chan<- SyncItem {
	c := make(chan SyncItem)
	go func() {
		for {
			if i, ok := <-c; ok {
				fmt.Println(i.Data())
			} else {
				break
			}
		}
	}()
	return c
}

type item struct {
	data int
}

func (i item) Attributes() Properties {
	return make(Properties)
}

func (i item) Data() interface{} {
	return i.data
}

type stubProducer struct{}

func (s stubProducer) Out() <-chan SyncItem {
	c := make(chan SyncItem)
	go func() {
		defer close(c)
		for i := 0; i < 4; i++ {
			c <- item{data: i}
			time.Sleep(time.Second)
		}

	}()
	return c
}

func TestPipelineWithProducerAndConsumer(t *testing.T) {
	pip := NewStreamWithoutMiddleware(&stubProducer{}, &stubConsumer{})
	pip.Start()
}
