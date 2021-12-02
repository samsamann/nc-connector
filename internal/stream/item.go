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
	SetPath(string)
	Attributes() Properties
	Data() io.Reader
}

type Properties map[string]interface{}

func (p Properties) Add() {

}

func (p Properties) Get() {

}

func NewFileSyncItem(path string, props Properties, content []byte) SyncItem {
	return &File{
		mode:    WRITE,
		path:    path,
		attrs:   props,
		content: bytes.NewReader(content),
	}
}

type File struct {
	mode    OperationMode
	path    string
	attrs   Properties
	content io.Reader
}

func (f *File) Mode() OperationMode {
	return f.mode
}
func (f *File) ChangeMode(m OperationMode) {
	f.mode = m
}
func (f *File) Path() string {
	return f.path
}
func (f *File) SetPath(path string) {
	f.path = path
}

func (f *File) Attributes() Properties {
	return f.attrs
}
func (f *File) SetAttributes(attrs Properties) {
	f.attrs = attrs
}
func (f *File) Data() io.Reader {
	return f.content
}
