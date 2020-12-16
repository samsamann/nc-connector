package webdav

import (
	"reflect"
)

var (
	// NotDefinedProp is the zero value representation of a WebDAV property
	NotDefinedProp = new(Property)

	ResourceTypeProp = NewProperty("resourcetype", davNamespace)
	// GetContentLengthProp represents the getlastmodified property in a WebDAV PropFind request.
	GetContentLengthProp = NewProperty("getcontentlength", davNamespace)
	// GetLastModifiedProp represents the getcontentlength property in a WebDAV PropFind request.
	GetLastModifiedProp = NewProperty("getlastmodified", davNamespace)
)

const davNamespace = "DAV:"

// Property represents a WebDAV property. It holds the name and namespace of the XML representation.
type Property struct {
	name  string
	xmlNS string
	pType reflect.Type
}

// NewProperty returns a WebDAV property
func NewProperty(name string, namespace string) *Property {
	return &Property{name: name, xmlNS: namespace}
}
