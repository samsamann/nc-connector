package stream

import (
	"bytes"
	"io"
)

type OperationMode uint8

const (
	NONE OperationMode = iota
	WRITE
	DELETE
)

// SyncItem reflects an object that can be synced.
type SyncItem interface {
	Mode() OperationMode
	ChangeMode(OperationMode)
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
		mode:    WRITE,
		path:    path,
		attrs:   props,
		content: bytes.NewReader(content),
	}
}

type file struct {
	mode    OperationMode
	path    string
	attrs   Properties
	content io.Reader
}

func (f *file) Mode() OperationMode {
	return f.mode
}
func (f *file) ChangeMode(m OperationMode) {
	f.mode = m
}
func (f *file) Path() string {
	return f.path
}
func (f *file) Attributes() Properties {
	return f.attrs
}
func (f *file) Data() io.Reader {
	return f.content
}
