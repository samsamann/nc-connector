package xml

import (
	"encoding/xml"
	"errors"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalResponse(t *testing.T) {
	tests := []struct {
		xmlStr string
		href   url.URL
	}{
		{"<p><href>http://localhost/test</href></p>", url.URL{Path: "/test"}},
	}

	for _, test := range tests {
		response := new(Response)
		xml.Unmarshal([]byte(test.xmlStr), response)
		assert.Equal(t, test.href.Path, response.Href.Path)
	}
}

func TestResponseStatus(t *testing.T) {
	tests := []struct {
		xmlStr     string
		statusCode int
		rawStatus  string
		err        error
	}{
		{"", 0, "", errors.New("EOF")},
		{"<empty></empty>", 0, "", nil},
		{"<s>HTTP/1.1 200 OK</s>", 200, "HTTP/1.1 200 OK", nil},
		{"<status>HTTP/2 404 FOOBAR</status>", 404, "HTTP/2 404 FOOBAR", nil},
		{"<s>HTTP/123 123 STATUS NN</s>", 0, "", errors.New("Can not parse response status: Status code did not match")},
		{"<s>HTTP/1.2 500 Internal Server Error<s", 0, "", errors.New("XML syntax error on line 1: unexpected EOF")},
	}

	for _, test := range tests {
		sut := new(responseStatus)
		err := xml.Unmarshal([]byte(test.xmlStr), &sut)
		if test.err == nil {
			assert.NoError(t, err)
			assert.Equal(t, test.statusCode, sut.StatusCode)
			assert.Equal(t, test.rawStatus, sut.RawStatus)
		} else {
			assert.EqualError(t, err, test.err.Error())
		}
	}
}

func TestPropStatProperties(t *testing.T) {
	sut := new(propStat)
	sut.Prop = new(propertyContainer)
	sut.Prop.addProperty(xml.Name{Local: "prop", Space: "space"}, nil)
	// sut.Prop.addProperty(xml.Name{Local: "prop2", Space: "space"}, nil)

	for tuple := range sut.Properties() {
		assert.Equal(t, "prop", tuple.Get(0))
		assert.Equal(t, "space", tuple.Get(1))
		assert.Nil(t, tuple.Get(2))
	}
}
