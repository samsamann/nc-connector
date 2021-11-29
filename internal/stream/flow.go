package stream

func NewStream(producer Producer) Flow {
	return newFlowElement(producer, nil)
}
func NewStreamWithoutMiddleware(producer Producer, consumer Consumer) Pipeline {
	return NewStream(producer).To(consumer)
}

type flow struct {
	prevFlow Flow
	elem     interface{}
}

func newFlowElement(elem interface{}, prev Flow) Flow {
	return &flow{
		prevFlow: prev,
		elem:     elem,
	}
}

func (f flow) To(consumer Consumer) Pipeline {
	return newPipeline(newDestElement(consumer, &f))
}

func (f flow) Via(operator Operator) Flow {
	return newFlowElement(operator, &f)
}

func (f flow) prev() linker {
	return f.prevFlow
}

func (f flow) element() interface{} {
	return f.elem
}
