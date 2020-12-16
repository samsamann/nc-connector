package webdav

import (
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/samsamann/nc-connector/pkg/webdav/xml"
)

const (
	httpMethodPropfind  = "PROPFIND"
	httpMethodProppatch = "PROPPATCH"
)

const (
	httpHeaderDepth = "Depth"
)

// RequestBuilder builds a HTTP request
type RequestBuilder interface {
	build(url string, extraHeaders ...HTTPHeader) (*http.Request, error)
}

// HTTPHeader represents a single http header entry
type HTTPHeader struct {
	Name  string
	Value string
}

type request struct {
	ctx     context.Context
	method  string
	url     *url.URL
	headers map[string]string
	body    xml.Marshaler
}

func (r request) build(url url.URL, extraHeaders ...HTTPHeader) (*http.Request, error) {
	r.url = url.ResolveReference(r.url)
	var reader io.Reader
	if r.body != nil {
		var err error
		reader, err = r.body.Marshal()
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequestWithContext(r.ctx, r.method, r.url.String(), reader)
	if err != nil {
		return nil, err
	}
	for header, val := range r.headers {
		req.Header.Set(header, val)
	}
	for _, header := range extraHeaders {
		req.Header.Add(header.Name, header.Value)
	}

	return req, nil
}
