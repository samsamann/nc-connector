package operator

import "github.com/samsamann/nc-connector/internal/stream"

type filter struct {
}

func NewFilter() stream.Operator {
	return new(filter)
}

func (f filter) In() chan<- stream.SyncItem {
	return nil
}

func (f filter) Out() <-chan stream.SyncItem {
	return nil
}

func (f filter) next() stream.Inlet {
	return nil
}
