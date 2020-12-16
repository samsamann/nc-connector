package xml

import (
	"bytes"
	"encoding/xml"
	"io"
)

const xmlHeader string = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>"

// Marshaler is the interface implemented by objects that can marshal
// themselves into valid XML elements.
type Marshaler interface {
	// Marshal encodes the receiver as zero or more XML elements.
	Marshal() (io.Reader, error)
}

func newReader(entity interface{}) (io.Reader, error) {
	body, err := xml.Marshal(entity)
	if err != nil {
		return nil, err
	}
	header := []byte(xmlHeader)
	return bytes.NewReader(append(header, body...)), nil
}
