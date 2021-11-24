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

func TestPipeline(t *testing.T) {
	pip := NewStreamWithoutMiddleware(&stubProducer{}, &stubConsumer{})
	pip.Start()
}

/*func TestIsPipelineOK(t *testing.T) {
	consumer := new(stubConsumer)
	producer1 := stubProducer{stubConsumer: stubConsumer{nextItem: consumer}}
	producer2 := stubProducer{stubConsumer: stubConsumer{nextItem: producer1}}

	assert.True(t, isPipelineOK(producer2))
	assert.True(t, isPipelineOK(producer1))
	assert.False(t, isPipelineOK(consumer))
	assert.False(t, isPipelineOK(nil))
}*/
