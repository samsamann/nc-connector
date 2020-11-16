package log

import (
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	// LevelFlagName is the flag name to configure the log level.
	LevelFlagName = "log.level"
	// LevelFlagHelp is the help description for the log.level flag.
	LevelFlagHelp = "Only log messages with the given severity or above. One of: [debug, info, warn, error]"
	// FormatFlagName is the flag name to configure the log format.
	FormatFlagName = "log.format"
	// FormatFlagHelp is the help description for the log.format flag.
	FormatFlagHelp = "Output format of log messages. One of: [logfmt, json]"
)

// AddFlags adds the flags used by this package to the Kingpin application.
func AddFlags(a *kingpin.Application, config *Config) {
	config.Level = &Level{}
	a.Flag(LevelFlagName, LevelFlagHelp).
		Default("info").SetValue(config.Level)

	config.Format = &Format{}
	a.Flag(FormatFlagName, FormatFlagHelp).
		Default("logfmt").SetValue(config.Format)
}
