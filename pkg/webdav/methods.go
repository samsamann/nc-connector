package webdav

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/samsamann/nc-connector/pkg/webdav/xml"
)

const defaultClientTimeout = 3 * time.Second

var defaultClient = &http.Client{
	Timeout: defaultClientTimeout,
}

func Propfind(ctx context.Context, builder *PropfindBuilder) error {
	requestCtx, _ := context.WithCancel(ctx)
	propfindReq := builder.build(requestCtx)
	// TODO: Make url configurable
	baseURL, _ := url.Parse("")
	req, err := propfindReq.build(
		*baseURL,
		HTTPHeader{Name: "Content-Type", Value: "application/xml; charset=\"utf-8\""},
		HTTPHeader{Name: "Accept", Value: "application/xml"},
	)
	if err != nil {
		return err
	}
	var multiStatus *xml.MultiStatus
	parseBodyFunc := func(responseBody []byte) {
		multiStatus, err = xml.ParseMultiStatusResponse(responseBody)
	}
	_, err = doRequest(req, parseBodyFunc)
	_ = multiStatus
	return err
}

func Proppatch(ctx context.Context, builder *ProppatchBuilder) error {
	requestCtx, _ := context.WithCancel(ctx)
	proppatchReq := builder.build(requestCtx)
	// TODO: Make url configurable
	baseURL, _ := url.Parse("")
	req, err := proppatchReq.build(
		*baseURL,
		HTTPHeader{Name: "Content-Type", Value: "application/xml; charset=\"utf-8\""},
		HTTPHeader{Name: "Accept", Value: "application/xml"},
	)
	if err != nil {
		return err
	}
	var multiStatus *xml.MultiStatus
	parseBodyFunc := func(responseBody []byte) {
		multiStatus, err = xml.ParseMultiStatusResponse(responseBody)
	}
	_, err = doRequest(req, parseBodyFunc)
	_ = multiStatus
	return err
}

func temp(t interface{}) {

}

func doRequest(req *http.Request, parseBodyFunc func([]byte)) (*http.Response, error) {
	res, err := defaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
	case http.StatusMultiStatus:
		responseBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		parseBodyFunc(responseBody)
		return res, nil
	default:
		return nil, fmt.Errorf("respone with status code %q can not be processed", res.Status)
	}
	return res, nil
}

func buildMultiStatusResponse(xmlMultiStatus *xml.MultiStatus) *MultistatusResponse {
	multiStatus := new(MultistatusResponse)
	for _, xmlRes := range xmlMultiStatus.Responses {
		response := new(Response)
		href := xmlRes.Href
		response.href = &href
		response.propertyStatus = make(map[*Property]*PropertyStatus)
		for _, prop := range xmlRes.PropStats {
			propStatus := NewStatus(prop.Status())
			for tuple := range prop.Properties() {
				prop := NewProperty(tuple.Get(0).(string), tuple.Get(1).(string))
				response.propertyStatus[prop] = newPropertyStatus(propStatus, tuple.Get(2))
			}
		}

		multiStatus.responses = append(multiStatus.responses, response)
	}

	return multiStatus
}
