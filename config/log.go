package config

import (
	"os"
	"strings"
)

type LoggingConfig struct {
	Level           string
	Format          string
	ExtendedLogging bool
}

func loadLoggingConfig() LoggingConfig {
	return LoggingConfig{
		Level:           os.Getenv("LOG_LEVEL"),
		Format:          os.Getenv("LOG_FORMAT"),
		ExtendedLogging: strings.ToLower(os.Getenv("EXTENDED_LOGGING")) == "true",
	}
}

// IsDebug returns true if log level is debug.
func (l *LoggingConfig) IsDebug() bool {
	return l.Level == "debug"
}
