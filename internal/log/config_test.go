package log

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestSetLogLevel(t *testing.T) {
	tests := []struct {
		level         string
		expectedValue logrus.Level
		errorMsg      string
	}{
		{level: "debug", expectedValue: logrus.DebugLevel},
		{level: "info", expectedValue: logrus.InfoLevel},
		{level: "warn", expectedValue: logrus.WarnLevel},
		{level: "error", expectedValue: logrus.ErrorLevel},
		{level: "", expectedValue: logrus.InfoLevel},
		{level: "panic", errorMsg: "unrecognized or unsupported log level \"panic\""},
		{level: "foo", errorMsg: "unrecognized or unsupported log level \"foo\""},
	}

	for _, test := range tests {
		sut := new(Level)
		err := sut.Set(test.level)
		if err == nil {
			assert.Empty(t, test.errorMsg)
			assert.Equal(t, test.expectedValue, sut.l)
		} else {
			assert.EqualError(t, err, test.errorMsg)
		}
	}
}

func TestLogLevelZeroValue(t *testing.T) {
	sut := new(Level)
	assert.Equal(t, logrus.PanicLevel.String(), sut.String())
}

func TestSetLogFormat(t *testing.T) {
	tests := []struct {
		format   string
		errorMsg string
	}{
		{format: "json"},
		{format: "logfmt"},
		{format: ""},
		{format: "foo", errorMsg: "unrecognized log format \"foo\""},
	}

	for _, test := range tests {
		sut := new(Format)
		err := sut.Set(test.format)
		if err == nil {
			assert.Empty(t, test.errorMsg)
			if test.format == "" {
				assert.Equal(t, DefaultLogFormat, sut.f)
			} else {
				assert.Equal(t, test.format, sut.f)
			}
		} else {
			assert.EqualError(t, err, test.errorMsg)
		}
	}
}
