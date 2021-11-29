package operator

import (
	"fmt"

	"github.com/samsamann/nc-connector/internal/stream"
)

type InitOperatorFunc func() (stream.Operator, error)

var operatorRegistry map[string]InitOperatorFunc

func init() {
	operatorRegistry = make(map[string]InitOperatorFunc)
}

func CreateOperator(name string, config map[string]interface{}) (stream.Operator, error) {
	if f, ok := operatorRegistry[name]; ok {
		return f()
	}
	return nil, fmt.Errorf("no operator found with name %q", name)
}
