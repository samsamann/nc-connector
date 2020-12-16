package xml

import (
	"encoding/xml"
	"io"
	"reflect"
)

// Propfind represents the xml Propfind structure.
type Propfind struct {
	XMLName xml.Name           `xml:"DAV: propfind"`
	Prop    *propertyContainer `xml:"DAV: prop"`
}

// NewPropfind returns a new instance of Propfind
func NewPropfind() *Propfind {
	pf := new(Propfind)
	pf.Prop = &propertyContainer{}

	return pf
}

// Property adds a property to the WebDAV property list.
func (pf *Propfind) Property(name string, ns string) {
	if pf.Prop == nil {
		pf.Prop = &propertyContainer{}
	}
	pf.Prop.addProperty(xml.Name{Local: name, Space: ns}, nil)
}

// Marshal encodes the Propfind object as xml element.
func (pf Propfind) Marshal() (io.Reader, error) {
	return newReader(pf)
}

// Proppatch represents a propertyupdate XML element in a PROPPATCH HTTP request.
type Proppatch struct {
	XMLName    xml.Name           `xml:"DAV: propertyupdate"`
	SetList    *propertyContainer `xml:"DAV: set>prop"`
	RemoveList *propertyContainer `xml:"DAV: remove>prop"`
}

// NewProppatch returns a new instance of Proppatch
func NewProppatch() *Proppatch {
	pp := new(Proppatch)
	pp.SetList = &propertyContainer{}
	pp.RemoveList = &propertyContainer{}

	return pp
}

// Set adds a property with a namespace and the value to be updated to a PROPPATCH
// HTTP request.
func (pp *Proppatch) Set(name string, ns string, value interface{}) {
	if pp.SetList == nil {
		return
	}
	pp.SetList.addProperty(xml.Name{Local: name, Space: ns}, value)
}

// Remove adds a property to the remove list in a PROPPATCH HTTP request.
func (pp *Proppatch) Remove(name string, ns string) {
	if pp.RemoveList == nil {
		return
	}
	pp.RemoveList.addProperty(xml.Name{Local: name, Space: ns}, nil)
}

// Marshal encodes the Proppatch object as xml element.
func (pp Proppatch) Marshal() (io.Reader, error) {
	if pp.SetList.empty() {
		pp.SetList = nil
	}
	if pp.RemoveList.empty() {
		pp.RemoveList = nil
	}
	return newReader(pp)
}

type propertyContainer struct {
	propList map[xml.Name]interface{}
}

func (pc propertyContainer) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	var err error
	err = e.EncodeToken(start)
	if err != nil {
		return err
	}
	for name, val := range pc.propList {
		if err = e.Encode(
			struct {
				XMLName xml.Name
				Value   interface{} `xml:",innerxml"`
			}{
				XMLName: name,
				Value:   val,
			},
		); err != nil {
			return err
		}
	}
	return e.EncodeToken(start.End())
}

type innerXML struct {
	XMLName xml.Name
	Value   interface{} `xml:",innerxml"`
}

func (pc *propertyContainer) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for {
		token, err := d.Token()
		if err != nil {
			return err
		}
		switch tokenType := token.(type) {
		case xml.StartElement:
			if val, ok := pc.propList[tokenType.Name]; ok {
				val := reflect.New(reflect.TypeOf(val))
				err = d.DecodeElement(val.Interface(), &tokenType)
				if err != nil {
					d.Skip()
				}
				pc.propList[tokenType.Name] = val.Elem().Interface()
			}
		case xml.EndElement:
			if tokenType == start.End() {
				return err
			}
		}
	}
}

func (pc *propertyContainer) addProperty(name xml.Name, value interface{}) {
	if pc.propList == nil {
		pc.propList = make(map[xml.Name]interface{})
	}
	if value == nil {
		pc.propList[name] = nil
		return
	}
	switch val := value.(type) {
	case reflect.Value:
		pc.propList[name] = val.Interface()
	case reflect.Type:
		pc.propList[name] = reflect.New(val).Interface()
	default:
		pc.propList[name] = val
	}
}

func (pc *propertyContainer) empty() bool {
	return len(pc.propList) == 0
}
