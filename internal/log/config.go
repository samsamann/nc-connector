package log

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	// DefaultLogLevel defines the default log level
	DefaultLogLevel string = "info"
	// DefaultLogFormat defines the default log format
	DefaultLogFormat string = "logfmt"
)

// Level defines the minimal log level.
type Level struct {
	l logrus.Level
}

func (lv *Level) String() string {
	return lv.l.String()
}

// Set updates the value of the log level.
func (lv *Level) Set(s string) error {
	switch s {
	case "debug":
		lv.l = logrus.DebugLevel
	case "info":
		lv.l = logrus.InfoLevel
	case "warn":
		lv.l = logrus.WarnLevel
	case "error":
		lv.l = logrus.ErrorLevel
	case "":
		return lv.Set(DefaultLogLevel)
	default:
		return errors.Errorf("unrecognized or unsupported log level %q", s)
	}

	return nil
}

// Format defines the log format.
type Format struct {
	f string
}

func (f *Format) String() string {
	return f.f
}

// Set updates the value of the log format.
func (f *Format) Set(format string) error {
	switch format {
	case "logfmt", "json":
		f.f = format
	case "":
		return f.Set(DefaultLogFormat)
	default:
		return errors.Errorf("unrecognized log format %q", format)
	}
	return nil
}

// Config is a struct containing log specific configuration values.
type Config struct {
	Level  *Level
	Format *Format
}
