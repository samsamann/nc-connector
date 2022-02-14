package consumer

import (
	"github.com/samsamann/nc-connector/internal/config"
	"github.com/samsamann/nc-connector/internal/stream"
)

const (
	stubConsumerName = "stub"
)

func initStubConsumer(global *config.GlobalConfig, opConfig map[string]interface{}) (stream.Consumer, error) {
	return &stubConsumer{waitChan: make(chan interface{})}, nil
}

type stubConsumer struct {
	waitChan chan interface{}
}

func (f stubConsumer) In(ctx stream.Context) chan<- stream.SyncItem {
	channel := make(chan stream.SyncItem)
	go func() {
		for range channel {
		}
		f.waitChan <- nil
	}()
	return channel
}

func (s stubConsumer) Wait() <-chan interface{} {
	return s.waitChan
}
