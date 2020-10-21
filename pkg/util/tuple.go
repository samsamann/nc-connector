package util

// Tuple represents an n-tuple
//
// A tuple is a finite ordered list (sequence) of elements.
type Tuple struct {
	data []interface{}
}

// NewTuple creates an empty Tuple of length n
func NewTuple(length int) *Tuple {
	t := new(Tuple)
	t.data = make([]interface{}, length)
	return t
}

// Get returns the item at index i
func (t Tuple) Get(i int) interface{} {
	if len(t.data) > i {
		item := t.data[i]
		return item
	}
	return nil
}

// Set adds or replaces the item at index i
func (t *Tuple) Set(i int, item interface{}) {
	if len(t.data) > i {
		t.data[i] = item
	}
}
