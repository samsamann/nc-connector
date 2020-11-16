package log

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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
	default:
		return errors.Errorf("unrecognized log level %q", s)
	}

	return nil
}

// Format defines the log format.
type Format struct {
	s string
}

func (f *Format) String() string {
	return f.s
}

// Set updates the value of the log format.
func (f *Format) Set(s string) error {
	switch s {
	case "logfmt", "json":
		f.s = s
	default:
		return errors.Errorf("unrecognized log format %q", s)
	}
	return nil
}

// Config is a struct containing log specific configuration values.
type Config struct {
	Level  *Level
	Format *Format
}
