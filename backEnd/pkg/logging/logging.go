package logging

import (
	"github.com/413ksz/BlueFox/backEnd/pkg/models"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func LogLevelHelperForError(customError *models.CustomError) *zerolog.Event {
	var logEvent *zerolog.Event
	switch customError.LogLevel {
	case zerolog.DebugLevel:
		logEvent = log.Debug()
	case zerolog.InfoLevel:
		logEvent = log.Info()
	case zerolog.WarnLevel:
		logEvent = log.Warn()
	case zerolog.ErrorLevel:
		logEvent = log.Error()
	default:
		logEvent = log.Error() // Default to error if level is unset/unknown
	}

	return logEvent
}
