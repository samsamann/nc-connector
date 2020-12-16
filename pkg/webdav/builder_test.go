package webdav

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPropfindBuilder(t *testing.T) {
	sot := new(PropfindBuilder)
	retVal := sot.Property(GetContentLengthProp).Property(GetLastModifiedProp)

	assert.Same(t, sot, retVal)
	assert.Len(t, sot.properties, 2)
	assert.Contains(t, sot.properties, GetContentLengthProp)
	assert.Contains(t, sot.properties, GetLastModifiedProp)
}
