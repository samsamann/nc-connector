package stream

type linker interface {
	prev() linker
	element() interface{}
}

type Inlet interface {
	In() chan<- SyncItem
}

type Outlet interface {
	Out() <-chan SyncItem
}

type Producer interface {
	Outlet
}

type Operator interface {
	Inlet
	Outlet
}

type Consumer interface {
	Inlet
	Wait() <-chan interface{}
}

type Flow interface {
	linker
	Via(Operator) Flow
	To(Consumer) Pipeline
}

type Pipeline interface {
	linker
	Start()
}
