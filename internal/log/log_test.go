package log

import (
	"bytes"
	"errors"
	"io"
	"regexp"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestLogError(t *testing.T) {
	tests := []struct {
		level       logrus.Level
		format      string
		params      []interface{}
		regexLogMsg string
	}{
		{
			level:       logrus.ErrorLevel,
			format:      "logfmt",
			params:      []interface{}{"error msg"},
			regexLogMsg: "time=\".*\" level=error msg=\"error msg\"",
		},
		{
			level:       logrus.FatalLevel,
			format:      "logfmt",
			params:      []interface{}{"error msg"},
			regexLogMsg: "",
		},
		{
			level:       logrus.InfoLevel,
			format:      "logfmt",
			params:      []interface{}{"test", "err", errors.New("Test")},
			regexLogMsg: "time=\".*\" level=error msg=test err=Test",
		},
		{
			level:       logrus.InfoLevel,
			format:      "json",
			params:      []interface{}{"json test", "key", 123},
			regexLogMsg: "{\"key\":123,\"level\":\"error\",\"msg\":\"json test\",\"time\":\".*\"}",
		},
	}
	for _, test := range tests {
		var buf bytes.Buffer
		sut, _ := NewLogger(io.Writer(&buf), &Config{Level: &Level{l: test.level}, Format: &Format{s: test.format}})
		sut.Error(test.params...)
		assert.Regexp(t, regexp.MustCompile(test.regexLogMsg), buf.String())
	}
}
