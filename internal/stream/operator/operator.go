package operator

import (
	"fmt"

	"github.com/samsamann/nc-connector/internal/config"
	"github.com/samsamann/nc-connector/internal/stream"
)

type InitOperatorFunc func(*config.GlobalConfig, map[string]interface{}) (stream.Operator, error)

var operatorRegistry map[string]InitOperatorFunc

func init() {
	operatorRegistry = make(map[string]InitOperatorFunc)
	operatorRegistry[extractPageOperatorrName] = initPageSplitOperator
	operatorRegistry[pathManipulatorName] = initPathManipulator
	operatorRegistry[renameOperatorName] = initRenameOperator
	operatorRegistry[splitOperatorrName] = initSplitOperator
	operatorRegistry[apiOperatorName] = initToPDFAPIOperator
}

func CreateOperator(opConfig config.StreamElem, c *config.GlobalConfig) (stream.Operator, error) {
	name := opConfig.Name
	if f, ok := operatorRegistry[name]; ok {
		return f(c, opConfig.Config)
	}
	return nil, fmt.Errorf("no operator found with name %q", name)
}
