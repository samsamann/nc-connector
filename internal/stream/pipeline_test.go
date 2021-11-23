package stream

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type stubConsumer struct {
	nextItem linker
}

func (s stubConsumer) next() linker {
	return s.nextItem
}

func (s stubConsumer) In() chan<- SyncItem {
	return nil
}

type stubProducer struct {
	stubConsumer
}

func (s stubProducer) out() <-chan SyncItem {
	return nil
}

func TestIsPipelineOK(t *testing.T) {
	consumer := new(stubConsumer)
	producer1 := stubProducer{stubConsumer: stubConsumer{nextItem: consumer}}
	producer2 := stubProducer{stubConsumer: stubConsumer{nextItem: producer1}}

	assert.True(t, isPipelineOK(producer2))
	assert.True(t, isPipelineOK(producer1))
	assert.False(t, isPipelineOK(consumer))
	assert.False(t, isPipelineOK(nil))
}
