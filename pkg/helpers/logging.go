package helpers

import (
	"context"
	"io"
	"net/url"
	"os"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/billiem/seren-management/pkg/projectpath"
	"github.com/charmbracelet/log"
	sqldblogger "github.com/simukti/sqldb-logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

/*
Loggers holds loggers for the application and database

We manage the operation logger seperately, as it is initialised only when an operation is started,
and also because it may require an additional io.Writer depending on the operation
*/
type Loggers struct {
	DBLogger  zap.SugaredLogger
	AppLogger SerenLogger
}

/*
BuildAppLoggers returns a Loggers struct containing loggers for the interface of the application,
and a database logger

These are build upon application startup
*/
func BuildAppLoggers(cfg Config) (*Loggers, error) {

	// create log directory if it doesn't exist
	logDirPath := JoinFilepathToSlash(projectpath.Root, "log")

	dbLogger, err := newLogger(
		cfg,
		JoinFilepathToSlash(logDirPath, "db.log"),
	)

	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("Error creating database logger"))
	}

	appLogger, err := newLogger(
		cfg,
		JoinFilepathToSlash(logDirPath, "app.log"),
	)

	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("Error creating application logger"))
	}

	return &Loggers{
		DBLogger: *dbLogger,
		AppLogger: SerenLogger{
			appLogger,
		},
	}, nil
}

/*
BuildOperationLogger returns a logger for use in operations

This is built before each operation is started
*/
func BuildOperationLogger(cfg Config) SerenLogger {

	// create log directory if it doesn't exist
	logDirPath := JoinFilepathToSlash(projectpath.Root, "log")

	opLogger, err := newLogger(
		cfg,
		JoinFilepathToSlash(logDirPath, "op.log"),
	)

	if err != nil {
		panic(fault.Wrap(err, fmsg.With("Error creating operation logger")))
	}

	return SerenLogger{
		opLogger,
	}
}

func newLogger(cfg Config, logPaths ...string) (*zap.SugaredLogger, error) {

	var logger *zap.Logger
	var err error

	if cfg.Development {
		logger, err = newDevelopmentConfig(logPaths...).Build()
	} else {
		logger, err = newProductionConfig(logPaths...).Build()
	}

	defer logger.Sync()

	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("Error building logger"))
	}

	return logger.Sugar(), nil
}

func newDevelopmentConfig(logPaths ...string) zap.Config {

	cfg := zap.NewDevelopmentConfig()

	cfg.OutputPaths = logPaths
	cfg.OutputPaths = append(cfg.OutputPaths, "stderr")

	return cfg
}

func newProductionConfig(logPaths ...string) zap.Config {

	cfg := zap.NewProductionConfig()

	cfg.OutputPaths = logPaths

	cfg.Encoding = "console"

	return cfg
}

/*
newTerminalConfig returns a zap.Config that writes to the provided io.Writer

We use this to add the terminal as an additional display for logs, this also
allows us to provide a custom zap Config for the terminal
*/
func newTerminalConfig(w io.Writer) zap.Config {

	zap.RegisterSink("termsink", func(u *url.URL) (zap.Sink, error) {
		return TermSink{
			w,
		}, nil
	})

	cfg := zap.NewProductionConfig()

	cfg.OutputPaths = []string{"termsink:term"}

	cfg.Encoding = "console"

	return cfg
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
	*zap.SugaredLogger
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

/*
AddTermCore adds a terminal writer to the logger
*/
func (s *SerenLogger) AddTermCore(w io.Writer, writeCallback func()) error {

	l, err := newTerminalConfig(w).Build()

	if err != nil {
		return fault.Wrap(err, fmsg.With("Error building terminal logger"))
	}

	coreOpt := zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.NewTee(core, l.Core())
	})

	hookOpt := zap.Hooks(func(entry zapcore.Entry) error {
		writeCallback()
		return nil
	})

	s.SugaredLogger = s.WithOptions(
		coreOpt,
		hookOpt,
	)

	return nil
}

/*
JSONUnmarshallWriter is a writer that attempts to unmarshal the bytes written to it into JSON

# If the unmarshaling fails, the bytes are written directly to the io.Writer
# If the unmarshaling succeeds, the unmarshaled JSON is written to the io.Writer
*/
// type JSONUnmarshallWriter struct {
// 	io.Writer
// }

// func NewJSONUnmarshallWriter(w io.Writer) JSONUnmarshallWriter {
// 	return JSONUnmarshallWriter{
// 		w,
// 	}
// }

// func (t JSONUnmarshallWriter) Write(p []byte) (n int, err error) {
// 	var data interface{}

// 	fmt.Println(string(p))

// 	// Attempt to unmarshal the array of bytes into JSON
// 	err = json.Unmarshal(p, &data)
// 	if err != nil {
// 		// If unmarshaling fails, write the bytes to the terminal
// 		_, err = t.Write(p)
// 		if err != nil {
// 			return 0, err
// 		}
// 	} else {
// 		// If unmarshaling succeeds, write the JSON to the terminal
// 		encoded, err := json.MarshalIndent(data, "", "  ")
// 		if err != nil {
// 			return 0, err
// 		}
// 		_, err = t.Write(encoded)
// 		if err != nil {
// 			return 0, err
// 		}
// 	}

// 	return len(p), nil
// }

// func (t JSONUnmarshallWriter) Close() error {
// 	return nil
// }

// Complies with type Sink interface (zapcore.WriteSyncer, io.Closer)
type TermSink struct {
	w io.Writer
}

func (t TermSink) Sync() error  { return nil }
func (t TermSink) Close() error { return nil }

func (t TermSink) Write(p []byte) (int, error) {
	t.w.Write(p)
	return len(p), nil
}
