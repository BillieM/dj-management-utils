package helpers

import (
	"context"
	"encoding/json"
	"fmt"
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

/*
Loggers holds loggers for the application and database

We manage the operation logger seperately, as it is initialised only when an operation is started,
and also because it may require an additional io.Writer depending on the operation
*/
type Loggers struct {
	DBLogger  log.Logger
	AppLogger SerenLogger
}

/*
BuildAppLoggers returns a Loggers struct containing loggers for the interface of the application,
and a database logger

These are build upon application startup
*/
func BuildAppLoggers() (*Loggers, error) {

	// create log directory if it doesn't exist
	logDirPath := JoinFilepathToSlash(projectpath.Root, "log")

	dbWriter, err := getLogWriter(JoinFilepathToSlash(logDirPath, "db.jsonl"))

	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("Error getting database log writer"))
	}

	appWriter, err := getLogWriter(JoinFilepathToSlash(logDirPath, "app.jsonl"))

	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("Error getting application log writer"))
	}

	dbLogger := newLogger(
		dbWriter,
		log.JSONFormatter,
		"DB",
	)

	appLogger := newLogger(
		appWriter,
		log.TextFormatter,
		"APP",
	)

	return &Loggers{
		DBLogger: *dbLogger,
		AppLogger: SerenLogger{
			*appLogger,
			appWriter,
		},
	}, nil
}

/*
BuildOperationLogger returns a logger for use in operations

This is built before each operation is started
*/
func BuildOperationLogger() SerenLogger {

	// create log directory if it doesn't exist
	logDirPath := JoinFilepathToSlash(projectpath.Root, "log")

	// main writer to log file
	opWriter, err := getLogWriter(JoinFilepathToSlash(logDirPath, "op.jsonl"))

	if err != nil {
		panic(fault.Wrap(err, fmsg.With("Error getting operation log writer")))
	}

	// secondary writer, converts structured logs to a nice format for the internal terminal log viewer

	opLogger := newLogger(
		opWriter,
		log.TextFormatter,
		"OP",
	)

	return SerenLogger{
		*opLogger,
		opWriter,
	}
}

func getTermWriter() io.Writer {
	return nil
}

/*
getLogWriters returns an io.Writer
*/
func getLogWriter(logPath string) (io.Writer, error) {

	// open db log file
	writer, err := os.OpenFile(
		logPath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644,
	)

	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("Error opening database log file"))
	}

	return writer, nil
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

/*
SerenLogger is a wrapper around the logging library used by the application

# This allows us to add additional helper functions to the logger easily

Additionally we store the io.Writer for the logger in order to simplify adding/ removing
a terminal writer, used to display logs in the application graphically to users
*/
type SerenLogger struct {
	log.Logger
	io.Writer
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

func (s SerenLogger) AddWriter(w io.Writer) {
	s.SetOutput(io.MultiWriter(s.Writer, w))
}

func (s SerenLogger) RemoveWriter() {
	s.SetOutput(s.Writer)
}

/*
JSONUnmarshallWriter is a writer that attempts to unmarshal the bytes written to it into JSON

# If the unmarshaling fails, the bytes are written directly to the io.Writer
# If the unmarshaling succeeds, the unmarshaled JSON is written to the io.Writer
*/
type JSONUnmarshallWriter struct {
	io.Writer
}

func NewJSONUnmarshallWriter(w io.Writer) JSONUnmarshallWriter {
	return JSONUnmarshallWriter{
		w,
	}
}

func (t JSONUnmarshallWriter) Write(p []byte) (n int, err error) {
	var data interface{}

	fmt.Println(string(p))

	// Attempt to unmarshal the array of bytes into JSON
	err = json.Unmarshal(p, &data)
	if err != nil {
		// If unmarshaling fails, write the bytes to the terminal
		_, err = t.Write(p)
		if err != nil {
			return 0, err
		}
	} else {
		// If unmarshaling succeeds, write the JSON to the terminal
		encoded, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return 0, err
		}
		_, err = t.Write(encoded)
		if err != nil {
			return 0, err
		}
	}

	return len(p), nil
}

func (t JSONUnmarshallWriter) Close() error {
	return nil
}
