package stream

import (
	"fmt"
	"io"
	"strconv"
	"strings"
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
	data string
}

func (i item) Path() string {
	return ""
}

func (i item) Attributes() Properties {
	return make(Properties)
}

func (i item) Data() io.Reader {
	return strings.NewReader(i.data)
}

type stubProducer struct{}

func (s stubProducer) Out() <-chan SyncItem {
	c := make(chan SyncItem)
	go func() {
		defer close(c)
		for i := 0; i < 4; i++ {
			c <- item{data: strconv.Itoa(i)}
			time.Sleep(time.Second)
		}

	}()
	return c
}

func TestPipelineWithProducerAndConsumer(t *testing.T) {
	pip := NewStreamWithoutMiddleware(&stubProducer{}, &stubConsumer{})
	pip.Start()
}
