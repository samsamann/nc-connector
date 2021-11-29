package consumer

import (
	"github.com/samsamann/nc-connector/internal/stream"
)

const (
	stubConsumerName = "stub"
)

func initStubConsumer() (stream.Consumer, error) {
	return &stubConsumer{}, nil
}

type stubConsumer struct{}

func (f stubConsumer) In() chan<- stream.SyncItem {
	channel := make(chan stream.SyncItem)
	go func() {
		for range channel {
		}
	}()
	return channel
}
