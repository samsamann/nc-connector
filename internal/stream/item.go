package stream

import (
	"bytes"
	"io"
)

// SyncItem reflects an object that can be synced.
type SyncItem interface {
	Path() string
	Attributes() Properties
	Data() io.Reader
}

type Properties map[string]interface{}

func (p Properties) Add() {

}

func (p Properties) Get() {

}

func NewFileSyncItem(path string, props Properties, content []byte) SyncItem {
	return &file{
		path:    path,
		attrs:   props,
		content: bytes.NewReader(content),
	}
}

type file struct {
	path    string
	attrs   Properties
	content io.Reader
}

func (f file) Path() string {
	return f.path
}
func (f file) Attributes() Properties {
	return f.attrs
}
func (f file) Data() io.Reader {
	return f.content
}
