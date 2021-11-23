package stream

type linker interface {
	next() linker
}

type Inlet interface {
	// linker
	In() chan<- SyncItem
}

type Outlet interface {
	//linker
	Out() <-chan SyncItem
}

type Producer interface {
	Outlet
}

type Operator interface {
	Inlet
	Outlet
	// Via(Operator) Operator
	// To(Consumer)
}

type Consumer interface {
	Inlet
}

type Flow interface {
	Via(Operator) Flow
	To(Consumer)
}
type operatorLink struct {
	opt Operator
}

/*func newFlowElement(pipeline *Pipeline) Flow {
	return new(flowImpl)
}

func (f flowImpl) Via(Operator) Flow {
	return newFlowElement(nil)
}

func (f flowImpl) To(Consumer) {

}*/
