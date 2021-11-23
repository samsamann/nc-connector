package stream

// SyncItem reflects an object that can be synced.
type SyncItem interface {
	Attributes() Properties
	Data() interface{}
}

type Properties map[string]interface{}

func (p Properties) Add() {

}

func (p Properties) Get() {

}
