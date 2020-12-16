package webdav

import (
	"context"
	"net/url"
	"strconv"

	"github.com/samsamann/nc-connector/pkg/webdav/xml"
)

type Builder interface {
	build(context.Context) *request
}

type baseBuilder struct {
	url *url.URL
}

func (bb *baseBuilder) setURL(rawURL string) error {
	var err error
	bb.url, err = url.Parse(rawURL)
	return err
}

type PropfindBuilder struct {
	*baseBuilder
	properties []*Property
	depth      int
}

func NewPropfindBuilder(rawUrl string) (*PropfindBuilder, error) {
	builder := new(PropfindBuilder)
	builder.baseBuilder = new(baseBuilder)
	if err := builder.setURL(rawUrl); err != nil {
		return nil, err
	}
	return builder, nil
}

func (pb *PropfindBuilder) Depth(depth int) *PropfindBuilder {
	pb.depth = depth
	return pb
}

func (pb *PropfindBuilder) Property(property *Property) *PropfindBuilder {
	if property != nil {
		pb.properties = append(pb.properties, property)
	}
	return pb
}

func (pb PropfindBuilder) build(ctx context.Context) *request {
	headers := make(map[string]string, 1)
	headers[httpHeaderDepth] = strconv.Itoa(pb.depth)

	propfind := &xml.Propfind{}
	for _, prop := range pb.properties {
		propfind.Property(prop.name, prop.xmlNS)
	}

	if ctx == nil {
		ctx = context.Background()
	}
	return &request{ctx: ctx, method: httpMethodPropfind, url: pb.url, headers: headers, body: propfind}
}

type ProppatchBuilder struct {
	*baseBuilder
	setList    map[*Property]interface{}
	removeList []*Property
}

func NewProppatchBuilder(rawUrl string) (*ProppatchBuilder, error) {
	builder := new(ProppatchBuilder)
	builder.baseBuilder = new(baseBuilder)
	builder.setList = make(map[*Property]interface{}, 0)
	if err := builder.setURL(rawUrl); err != nil {
		return nil, err
	}
	return builder, nil
}

func (pb *ProppatchBuilder) Set(property *Property, value interface{}) *ProppatchBuilder {
	if property != nil && value != nil {
		pb.setList[property] = value
	}
	return pb
}

func (pb *ProppatchBuilder) Remove(property *Property) *ProppatchBuilder {
	if property != nil {
		pb.removeList = append(pb.removeList, property)
	}
	return pb
}

func (pb ProppatchBuilder) build(ctx context.Context) *request {
	proppatch := xml.NewProppatch()
	for prop, val := range pb.setList {
		proppatch.Set(prop.name, prop.xmlNS, val)
	}
	for _, prop := range pb.removeList {
		proppatch.Remove(prop.name, prop.xmlNS)
	}

	if ctx == nil {
		ctx = context.Background()
	}
	return &request{ctx: ctx, method: httpMethodProppatch, url: pb.url, body: proppatch}
}
