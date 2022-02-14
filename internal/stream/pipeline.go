package stream

import (
	"sync"

	"github.com/sirupsen/logrus"
)

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

func (p *pipeline) Start(logger *logrus.Logger) {
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
	ctx := newContext(logger)
	ctx.Info("Pipeline start")
	execPipeline(ctx, producer, operators, consumer)
	ctx.Info("Pipeline end")
}

func execPipeline(ctx Context, p Producer, ops []Operator, c Consumer) {
	var outlet Outlet = p
	for _, operator := range ops {
		go transmit(ctx, operator, outlet)
		outlet = operator
	}
	go transmit(ctx, c, outlet)
	<-c.Wait()
}
