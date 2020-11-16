package log

import (
	"fmt"
	"io"

	log "github.com/sirupsen/logrus"
)

type logger struct {
	entry *log.Entry
}

// Error logs a message at level Error on the standard logger.
func (l logger) Error(args ...interface{}) {
	if len(args) > 2 && len(args[1:])%2 == 0 {
		addFields(Logger(&l), args[1:]...).Error(args[0])
	} else {
		l.entry.Error(args...)
	}
}

func (l logger) With(key string, value interface{}) Logger {
	return logger{l.entry.WithField(key, value)}
}

// Logger is the interface for loggers used in this program.
type Logger interface {
	Error(...interface{})

	With(key string, value interface{}) Logger
}

// NewLogger returns a new instance of Logger.
func NewLogger(w io.Writer, config *Config) (Logger, error) {
	l := log.New()
	l.Out = w
	l.SetLevel(config.Level.l)
	if err := setFormat(l, config.Format.String()); err != nil {
		return nil, err
	}

	return &logger{entry: log.NewEntry(l)}, nil
}

func setFormat(l *log.Logger, format string) error {
	switch format {
	case "json":
		l.SetFormatter(new(log.JSONFormatter))
	case "logfmt":
		l.SetFormatter(&log.TextFormatter{
			DisableColors: true,
		})
	default:
		return fmt.Errorf("unsupported logger %q", format)
	}
	return nil
}

func addFields(l Logger, args ...interface{}) Logger {
	for i := 0; i < len(args); i += 2 {
		l = l.With(args[i].(string), args[i+1])
	}
	return l
}
