package webdav

import (
	"net/url"
)

type ResponseScanner interface {
	Status(Property) bool
	Scan(Property, interface{}) error
}

type MultistatusResponse struct {
	responses []*Response
}

/* func (pr PropfindResponse) Scan(p Property, valReciv interface{}) {

}*/

type Response struct {
	href           *url.URL
	propertyStatus map[*Property]*PropertyStatus
}

func (r Response) Scan(p *Property, valReciv interface{}) {

}

type PropertyStatus struct {
	status *Status
	value  interface{}
}

func newPropertyStatus(status *Status, value interface{}) *PropertyStatus {
	ps := new(PropertyStatus)
	ps.status = status
	ps.value = value

	return ps
}

func (ps PropertyStatus) Status() *Status {
	return ps.status
}

type Status struct {
	statusCode int
	rawStatus  string
}

func NewStatus(statusCode int, rawStatus string) *Status {
	s := new(Status)
	s.statusCode = statusCode
	s.rawStatus = rawStatus
	return s
}

func (s Status) StatusCode() int {
	return s.statusCode
}

func (s Status) RawStatus() string {
	return s.rawStatus
}
