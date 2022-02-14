package stream

import (
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

type stubConsumer struct {
	waitChan chan interface{}
}

func (s stubConsumer) In(ctx Context) chan<- SyncItem {
	c := make(chan SyncItem)
	go func() {
		for {
			if i, ok := <-c; ok {
				fmt.Println(i.Data())
			} else {
				break
			}
		}
		s.waitChan <- nil
	}()
	return c
}

func (s stubConsumer) Wait() <-chan interface{} {
	return s.waitChan
}

type item struct {
	data string
}

func (i item) Mode() OperationMode {
	return NONE
}

func (i item) ChangeMode(OperationMode) {}

func (i item) Path() string {
	return ""
}

func (i item) SetPath(path string) {}

func (i item) Attributes() Properties {
	return make(Properties)
}

func (i *item) Data() io.Reader {
	return strings.NewReader(i.data)
}

func (i *item) SetData(reader io.ReadCloser) {
	data, _ := ioutil.ReadAll(reader)
	i.data = string(data)
}

type stubProducer struct{}

func (s stubProducer) Out(ctx Context) <-chan SyncItem {
	c := make(chan SyncItem)
	go func() {
		defer close(c)
		for i := 0; i < 4; i++ {
			c <- &item{data: strconv.Itoa(i)}
			time.Sleep(time.Second)
		}

	}()
	return c
}

func TestPipelineWithProducerAndConsumer(t *testing.T) {
	pip := NewStreamWithoutMiddleware(&stubProducer{}, &stubConsumer{waitChan: make(chan interface{})})
	pip.Start(logrus.New())
}
