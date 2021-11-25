package stream

import "errors"

// type HandlerFunc func(SyncItem) interface{}

func getProducer(link linker) (Producer, error) {
	if link.prev() == nil {
		if p, ok := link.element().(Producer); ok {
			return p, nil
		}
		return nil, errors.New("no producer found")
	}
	return getProducer(link.prev())
}

func getAllOperators(link linker) []Operator {
	if link == nil {
		return []Operator{}
	}

	outlets := getAllOperators(link.prev())
	if o, ok := link.element().(Operator); ok {
		return append(outlets, o)
	}
	return outlets
}

func transmit(i Inlet, o Outlet) {
	defer close(i.In())
	for ele := range o.Out() {
		i.In() <- ele
	}
}

/*func DoStream(i Inlet, o Outlet) {
	go transmit(i, o)
}*/