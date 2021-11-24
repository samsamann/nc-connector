package stream

import "sync"

type dest struct {
	prevFlow Flow
	elem     Inlet
}

func newDestElement(elem Inlet, prev Flow) *dest {
	return &dest{
		prevFlow: prev,
		elem:     elem,
	}
}

func (d dest) prev() linker {
	return d.prevFlow
}

func (d dest) element() interface{} {
	return d.elem
}

type pipeline struct {
	*dest
	m sync.Mutex
}

func newPipeline(dest *dest) Pipeline {
	return &pipeline{
		dest: dest,
	}
}

func (p *pipeline) Start() {
	p.m.Lock()
	defer p.m.Unlock()

	producer, err := getProducer(p.dest)
	if err != nil {
		return
	}
	operators := getAllOperators(p.dest)
	var consumer Consumer
	if c, ok := p.dest.element().(Consumer); ok {
		consumer = c
	} else {
		return
	}

	var wg sync.WaitGroup
	wg.Add(len(operators) + 1)
	var outlet Outlet = producer
	for _, o := range operators {
		go func() {
			defer wg.Done()
			transmit(o, outlet)
		}()
		outlet = o
	}
	go func() {
		defer wg.Done()
		transmit(consumer, outlet)
	}()
	wg.Wait()
}

/*func isPipelineOK(item linker) bool {
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
}*/
