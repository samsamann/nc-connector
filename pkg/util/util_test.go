package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTuple(t *testing.T) {
	sut := NewTuple(2)
	sut.Set(0, "test")
	assert.Equal(t, "test", sut.Get(0))
	sut.Set(0, 56)
	assert.Equal(t, 56, sut.Get(0))
	assert.Nil(t, sut.Get(1))
	assert.Nil(t, sut.Get(999))

	sut.Set(1010, "not set")
	assert.Nil(t, sut.Get(1010))
}

func TestGetCaller(t *testing.T) {
	assert.Equal(t, "github.com/samsamann/nc-connector/pkg/util.TestGetCaller", GetFuncName())
}
