package helpers

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/charmbracelet/log"
	sqldblogger "github.com/simukti/sqldb-logger"
)

type AppLogger struct {
	DBLogger  *log.Logger
	GUILogger *log.Logger
	CLILogger *log.Logger
	OPLogger  *log.Logger
}

func BuildAppLogger(c Config) (*AppLogger, error) {

	dbF, err := os.OpenFile("log/db.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return nil, err
	}

	return &AppLogger{
		DBLogger: newLogger(io.MultiWriter(
			dbF,
		), "DB"),
		GUILogger: newLogger(os.Stderr, "GUI"),
		CLILogger: newLogger(os.Stderr, "CLI"),
		OPLogger:  newLogger(os.Stderr, "OP"),
	}, nil
}

func newLogger(w io.Writer, prefix string) *log.Logger {
	logger := log.NewWithOptions(w, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		TimeFormat:      time.DateTime,
	})
	if prefix != "" {
		logger.SetPrefix(prefix)
	}
	return logger
}

type CharmLogAdapter struct {
	log.Logger
}

func (c *CharmLogAdapter) Log(_ context.Context, level sqldblogger.Level, msg string, data map[string]interface{}) {

	switch level {
	case sqldblogger.LevelError:
		c.Logger.Error(msg, data)
	case sqldblogger.LevelInfo:
		c.Logger.Info(msg, data)
	case sqldblogger.LevelDebug:
		c.Logger.Debug(msg, data)
	default:
		c.Logger.Debug(msg, data)
	}
}
