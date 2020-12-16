package webdav

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPropfind(t *testing.T) {
	builder, err := NewPropfindBuilder("addressbooks/users/foo/")
	assert.NoError(t, err)
	builder.Depth(0)
	builder.Property(ResourceTypeProp)

	err = Propfind(context.Background(), builder)
	assert.NoError(t, err)
}

func TestProppatch(t *testing.T) {
	builder, err := NewProppatchBuilder("addressbooks/users/foo/")
	assert.NoError(t, err)
	builder.Remove(NewProperty("custom-prop", "http://property.local/ns"))

	err = Proppatch(context.Background(), builder)
	assert.NoError(t, err)
}
