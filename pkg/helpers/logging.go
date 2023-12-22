package helpers

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/billiem/seren-management/pkg/projectpath"
	"github.com/charmbracelet/log"
	sqldblogger "github.com/simukti/sqldb-logger"
)

type AppLogger struct {
	DBLogger *log.Logger
	UILogger *log.Logger
	OPLogger *log.Logger
}

func BuildAppLogger(c Config) (*AppLogger, error) {

	// create log directory if it doesn't exist
	logDirPath := JoinFilepathToSlash(projectpath.Root, "log")

	// open db log file
	dbF, err := os.OpenFile(
		JoinFilepathToSlash(logDirPath, "db.jsonl"),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644,
	)

	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("Error opening database log file"))
	}

	defer dbF.Close()

	// open ui log file
	uiF, err := os.OpenFile(
		JoinFilepathToSlash(logDirPath, "ui.log"),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644,
	)

	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("Error opening UI log file"))
	}

	defer uiF.Close()

	// open op log file
	opF, err := os.OpenFile(
		JoinFilepathToSlash(logDirPath, "op.log"),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644,
	)

	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("Error opening OP log file"))
	}

	// create loggers

	dbLogger := newLogger(
		dbF,
		log.JSONFormatter,
		"DB",
	)

	uiLogger := newLogger(
		io.MultiWriter(
			os.Stderr,
			uiF,
		),
		log.TextFormatter,
		"UI",
	)

	opLogger := newLogger(
		io.MultiWriter(
			os.Stderr,
			opF,
		),
		log.TextFormatter,
		"OP",
	)

	defer opF.Close()

	return &AppLogger{
		DBLogger: dbLogger,
		UILogger: uiLogger,
		OPLogger: opLogger,
	}, nil
}

func newLogger(w io.Writer, formatter log.Formatter, prefix string) *log.Logger {
	logger := log.NewWithOptions(w, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		TimeFormat:      time.DateTime,
		Formatter:       formatter,
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
