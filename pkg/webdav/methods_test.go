package webdav

import (
	"net/url"
	"testing"

	"github.com/samsamann/nc-connector/pkg/webdav/xml"
)

func TestBuildMultiStatusResponse(t *testing.T) {
	multistatusMock := new(xml.MultiStatus)

	multistatusMock.Responses = append(multistatusMock.Responses, xml.Response{Href: url.URL{Path: "test"}})
	res := buildMultiStatusResponse(multistatusMock)
	temp(res)
}
