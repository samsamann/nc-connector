package xml

import (
	"encoding/xml"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPropContainerMarshal(t *testing.T) {
	sut := new(propertyContainer)
	sut.addProperty(xml.Name{Local: "prop", Space: "ns"}, "test")

	bytes, err := xml.Marshal(sut)
	assert.NoError(t, err)
	assert.Equal(t, "<propertyContainer><prop xmlns=\"ns\">test</prop></propertyContainer>", string(bytes))

	sut = new(propertyContainer)
	sut.addProperty(xml.Name{Local: "p", Space: "DAV:"}, nil)

	bytes, err = xml.Marshal(sut)
	assert.NoError(t, err)
	assert.Equal(t, "<propertyContainer><p xmlns=\"DAV:\"></p></propertyContainer>", string(bytes))

	sut = new(propertyContainer)
	sut.addProperty(xml.Name{Local: "prop", Space: "ns"}, reflect.ValueOf(string("Test")))

	bytes, err = xml.Marshal(sut)
	assert.NoError(t, err)
	assert.Equal(t, "<propertyContainer><prop xmlns=\"ns\">Test</prop></propertyContainer>", string(bytes))

	sut = new(propertyContainer)
	sut.addProperty(
		xml.Name{Local: "prop", Space: "ns"},
		struct {
			XMLName xml.Name
			Attr    string `xml:"a,attr"`
			Val     string `xml:"val"`
		}{XMLName: xml.Name{Local: "struct"}, Attr: "attr", Val: "test"},
	)

	bytes, err = xml.Marshal(sut)
	assert.NoError(t, err)
	assert.Equal(t, "<propertyContainer><prop xmlns=\"ns\"><struct a=\"attr\"><val>test</val></struct></prop></propertyContainer>", string(bytes))
}

func TestPropContainerUnmarshal(t *testing.T) {
	sut := new(propertyContainer)
	err := xml.Unmarshal([]byte("<propertyContainer></propertyContainer>"), sut)
	assert.NoError(t, err)
	assert.Empty(t, sut.propList)

	sut = new(propertyContainer)
	err = xml.Unmarshal([]byte("<propertyContainer><prop xmlns=\"ns\">test</prop></propertyContainer>"), sut)
	assert.NoError(t, err)
	assert.Empty(t, sut.propList)

	sut = new(propertyContainer)
	sut.addProperty(xml.Name{Local: "p", Space: "ns"}, "")
	err = xml.Unmarshal([]byte("<propertyContainer><p xmlns=\"ns\">test</p></propertyContainer>"), sut)
	assert.NoError(t, err)
	assert.IsType(t, string(""), sut.propList[xml.Name{Local: "p", Space: "ns"}])
	assert.Equal(t, "test", sut.propList[xml.Name{Local: "p", Space: "ns"}])

	sut = new(propertyContainer)
	sut.addProperty(xml.Name{Local: "p", Space: "ns"}, int(0))
	err = xml.Unmarshal([]byte("<propertyContainer><p xmlns=\"ns\">1234</p></propertyContainer>"), sut)
	assert.NoError(t, err)
	assert.IsType(t, int(0), sut.propList[xml.Name{Local: "p", Space: "ns"}])
	assert.Equal(t, 1234, sut.propList[xml.Name{Local: "p", Space: "ns"}])

	sut = new(propertyContainer)
	sut.addProperty(xml.Name{Local: "p1", Space: "ns"}, int(0))
	sut.addProperty(xml.Name{Local: "p2", Space: "ns"}, int(0))
	err = xml.Unmarshal([]byte("<propertyContainer><p1 xmlns=\"ns\">4321</p1><p2 xmlns=\"ns\">1234</p2></propertyContainer>"), sut)
	assert.NoError(t, err)
	assert.Len(t, sut.propList, 2)
	assert.Nil(t, sut.propList[xml.Name{Local: "p", Space: "ns"}])
	assert.Equal(t, 4321, sut.propList[xml.Name{Local: "p1", Space: "ns"}])
}

func TestProppatchMarshal(t *testing.T) {
	expectedOutput := xmlHeader + "<propertyupdate xmlns=\"DAV:\"><set><prop xmlns=\"DAV:\"><getcontentlength xmlns=\"DAV:\">test</getcontentlength></prop></set></propertyupdate>"
	proppatch := NewProppatch()
	proppatch.Set("getcontentlength", "DAV:", "test")
	encodeProppatch(t, proppatch, expectedOutput)

	expectedOutput = xmlHeader + "<propertyupdate xmlns=\"DAV:\"></propertyupdate>"
	proppatch = NewProppatch()
	encodeProppatch(t, proppatch, expectedOutput)

	expectedOutput = xmlHeader + "<propertyupdate xmlns=\"DAV:\"><remove><prop xmlns=\"DAV:\"><test xmlns=\"http://test.net/ns\"></test></prop></remove></propertyupdate>"
	proppatch = NewProppatch()
	proppatch.Remove("test", "http://test.net/ns")
	encodeProppatch(t, proppatch, expectedOutput)
}

func encodeProppatch(t *testing.T, propFind *Proppatch, expectedOutput string) {
	reader, err := propFind.Marshal()
	if assert.NoError(t, err) {
		output, err := ioutil.ReadAll(reader)
		assert.NoError(t, err)
		assert.Equal(t, expectedOutput, string(output))
	}
}
