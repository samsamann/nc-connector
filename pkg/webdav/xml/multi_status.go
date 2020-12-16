package xml

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"

	"github.com/samsamann/nc-connector/pkg/util"
)

// MultiStatus represents the Multi-Status response body, which conveys information about
// multiple resources in situations where multiple status codes might be appropriate.
type MultiStatus struct {
	Responses []Response `xml:"response"`
}

// Response holds the information about a resource and its properties which was effected
// by the previous method.
type Response struct {
	Href      url.URL        `xml:"href"`
	Status    responseStatus `xml:"status"`
	PropStats []propStat     `xml:"propstat"`
}

// UnmarshalXML decodes a single response xml element.
func (r *Response) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for {
		t, err := d.Token()
		if err != nil {
			return err
		}
		switch tt := t.(type) {
		case xml.StartElement:
			switch tt.Name.Local {
			case "href":
				var href string
				if err := d.DecodeElement(&href, &tt); err != nil {
					return err
				}
				parsedURL, err := url.Parse(href)
				if err != nil {
					return err
				}
				r.Href = *parsedURL

			case "status":
				if err := d.DecodeElement(&r.Status, &tt); err != nil {
					return err
				}

			case "propstat":
				el := &propStat{}
				err = d.DecodeElement(el, &tt)
				if err != nil {
					return err
				}
				r.PropStats = append(r.PropStats, *el)
			}
			break
		case xml.EndElement:
			if tt == start.End() {
				return nil
			}
		}
	}
}

type propStat struct {
	ResStatus responseStatus     `xml:"status"`
	Prop      *propertyContainer `xml:"prop"`
}

// Properties returns all properties as a tuple
func (ps propStat) Properties() <-chan *util.Tuple {
	chnl := make(chan *util.Tuple)
	go func(pc propertyContainer) {
		for prop, val := range pc.propList {
			tuple := util.NewTuple(3)
			tuple.Set(0, prop.Local)
			tuple.Set(1, prop.Space)
			tuple.Set(2, val)
			chnl <- tuple
		}
		close(chnl)
	}(*ps.Prop)
	return chnl
}

// Status returns the HTTP status as integer and string.
func (ps propStat) Status() (int, string) {
	return ps.ResStatus.StatusCode, ps.ResStatus.RawStatus
}

const statusRegex = `^HTTP/(?:1\.0|1\.1|2)\s((?:1|2|3|4|5)\d{2})`

const statusCodeErrMsg = "Can not parse response status: %w"

type responseStatus struct {
	StatusCode int
	RawStatus  string
}

func (rs *responseStatus) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	if err := d.DecodeElement(&rs.RawStatus, &start); err != nil {
		return err
	}
	if rs.RawStatus == "" {
		return nil
	}

	validateHTTPStatus := regexp.MustCompile(statusRegex)
	match := validateHTTPStatus.FindStringSubmatch(rs.RawStatus)
	if match == nil || len(match) != 2 {
		return fmt.Errorf(statusCodeErrMsg, errors.New("Status code did not match"))
	}
	rs.StatusCode, _ = strconv.Atoi(match[1])

	return nil
}

// ParseMultiStatusResponse reads the response body and generate a MultiStatus object.
func ParseMultiStatusResponse(bytes []byte) (*MultiStatus, error) {
	obj := new(MultiStatus)
	return obj, xml.Unmarshal(bytes, obj)
}
