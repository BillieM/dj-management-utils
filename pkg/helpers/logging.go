package helpers

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/billiem/seren-management/pkg/projectpath"
	"github.com/charmbracelet/log"
	sqldblogger "github.com/simukti/sqldb-logger"
)

type Loggers struct {
	DBLogger  log.Logger
	AppLogger SerenLogger
}

type logWriters struct {
	dbW  io.Writer
	appW io.Writer
}

func BuildAppLoggers(c Config) (*Loggers, error) {

	// create log directory if it doesn't exist
	logDirPath := JoinFilepathToSlash(projectpath.Root, "log")

	logWriters, err := getLogWriters(c, logDirPath)

	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("Error getting log writers"))
	}

	// create loggers
	dbLogger := newLogger(
		logWriters.dbW,
		log.JSONFormatter,
		"DB",
	)

	appLogger := newLogger(
		logWriters.appW,
		log.TextFormatter,
		"APP",
	)

	return &Loggers{
		DBLogger: *dbLogger,
		AppLogger: SerenLogger{
			*appLogger,
		},
	}, nil
}

func getLogWriters(c Config, logDirPath string) (logWriters, error) {

	// open db log file
	dbW, err := os.OpenFile(
		JoinFilepathToSlash(logDirPath, "db.jsonl"),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644,
	)

	if err != nil {
		return logWriters{}, fault.Wrap(err, fmsg.With("Error opening database log file"))
	}

	appW, err := os.OpenFile(
		JoinFilepathToSlash(logDirPath, "app.jsonl"),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644,
	)

	if err != nil {
		return logWriters{}, fault.Wrap(err, fmsg.With("Error opening app log file"))
	}

	return logWriters{
		dbW:  dbW,
		appW: appW,
	}, nil
}

func newLogger(w io.Writer, formatter log.Formatter, prefix string) *log.Logger {
	logger := log.NewWithOptions(w, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		TimeFormat:      time.DateTime,
		Formatter:       formatter,
		Level:           log.DebugLevel,
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

type SerenLogger struct {
	log.Logger
}

func (s SerenLogger) NonFatalError(err error) {
	// get err chain
	chain := fault.Flatten(err)

	// get err context
	ectx := fctx.Unwrap(err)

	s.Error(
		chain[0].Message,
		"caller", chain[0].Location,
		"error context", ectx,
		"error chain", chain,
	)
}

func (s SerenLogger) FatalError(err error) {
	// get err chain
	chain := fault.Flatten(err)

	// get err context
	ectx := fctx.Unwrap(err)

	s.Fatal(
		chain[0].Message,
		"caller", chain[0].Location,
		"error context", ectx,
		"error chain", chain,
	)

	os.Exit(1)
}
