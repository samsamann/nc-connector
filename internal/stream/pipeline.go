package stream

import "sync"

type Pipeline struct {
	m        sync.Mutex
	started  bool
	producer Producer
}

func NewPipeline(producer Producer) Flow {
	return &Pipeline{
		started:  false,
		producer: producer,
	}
}

func (p *Pipeline) Start() {
	p.m.Lock()
	if p.started {
		p.m.Unlock()
		return
	}
	p.started = true
	p.m.Unlock()
}

func (p *Pipeline) Via(Operator) Flow {
	return nil //newFlowElement(nil)
}

func (p *Pipeline) To(Consumer) {

}

func (p *Pipeline) Source(producer Producer) Flow {
	return nil
}

func isPipelineOK(item linker) bool {
	if item == nil {
		return false
	}
	nextItem := item.next()
	for nextItem != nil {
		if _, ok := nextItem.(Consumer); ok && nextItem.next() == nil {
			return true
		}
		nextItem = nextItem.next()
	}
	return false
}
