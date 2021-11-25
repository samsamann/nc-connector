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
	execPipeline(producer, operators, consumer)
}

func execPipeline(p Producer, ops []Operator, c Consumer) {
	var wg sync.WaitGroup
	wg.Add(len(ops) + 1)

	exec := func(i Inlet, o Outlet) {
		defer wg.Done()
		transmit(i, o)
	}

	var outlet Outlet = p
	for _, operator := range ops {
		go exec(operator, outlet)
		outlet = operator
	}
	go exec(c, outlet)
	wg.Wait()
}
