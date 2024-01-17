package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"time"

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

These are built upon application startup
*/
func BuildAppLoggers(cfg Config) (*Loggers, error) {

	// create log directory if it doesn't exist
	logDirPath := JoinFilepathToSlash(projectpath.Root, "log")

	dbLogger, err := newLogEnvConfig(
		cfg,
		JoinFilepathToSlash(logDirPath, "db.log"),
	).Build()

	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("Error creating database logger"))
	}

	appLogger, err := newLogEnvConfig(
		cfg,
		JoinFilepathToSlash(logDirPath, "app.log"),
	).Build()

	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("Error creating application logger"))
	}

	return &Loggers{
		DBLogger: *dbLogger.Sugar(),
		AppLogger: SerenLogger{
			appLogger.Sugar(),
		},
	}, nil
}

/*
BuildOperationLogger returns a logger for use in operations

This is built before each operation is started
*/
func BuildOperationLogger(cfg Config, termSink *TermSink) SerenLogger {

	// create log directory if it doesn't exist
	logDirPath := JoinFilepathToSlash(projectpath.Root, "log")
	CreateDirIfNotExists(logDirPath)

	opLogger, err := newLogEnvConfig(
		cfg,
		JoinFilepathToSlash(logDirPath, "op.log"),
	).Build()

	if err != nil {
		panic(fault.Wrap(err, fmsg.With("Error creating operation logger")))
	}

	termLogger, err := newTerminalConfig(termSink).Build()

	if err != nil {
		panic(fault.Wrap(err, fmsg.With("Error building terminal logger")))
	}

	coreOpt := zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.NewTee(core, termLogger.Core())
	})

	defer opLogger.Sync()
	defer termLogger.Sync()

	return SerenLogger{
		opLogger.Sugar().WithOptions(
			coreOpt,
		),
	}
}

/*
newLogEnvConfig returns a zap.Config that writes to the provided log paths,
depending on the development flag, this will use a development or production config
*/
func newLogEnvConfig(cfg Config, logPaths ...string) zap.Config {
	if cfg.Development {
		return newDevelopmentConfig(
			JoinFilepathToSlash(logPaths...),
		)
	} else {
		return newProductionConfig(
			JoinFilepathToSlash(logPaths...),
		)
	}
}

/*
newDevelopmentConfig returns a zap.Config that writes stderr as well
as any other log paths provided
*/
func newDevelopmentConfig(logPaths ...string) zap.Config {

	cfg := zap.NewDevelopmentConfig()

	cfg.OutputPaths = logPaths
	cfg.OutputPaths = append(cfg.OutputPaths, "stderr")

	return cfg
}

/*
newProductionConfig returns a zap.Config that writes to the provided log paths,
unlike the development config, this does not write to stderr
*/
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
func newTerminalConfig(s *TermSink) zap.Config {

	zap.RegisterSink("termsink", func(u *url.URL) (zap.Sink, error) {
		return s, nil
	})

	cfg := zap.NewProductionConfig()

	cfg.OutputPaths = []string{"termsink:term"}
	cfg.Encoding = "json"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

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

	// get err issues
	issues := fmsg.GetIssues(err)

	s.Errorw(
		chain[0].Message,
		"ctx", ectx,
		"chain", chain,
		"issues", issues,
	)
}

func (s SerenLogger) FatalError(err error) {
	// get err chain
	chain := fault.Flatten(err)

	// get err context
	ectx := fctx.Unwrap(err)

	// get err issues
	issues := fmsg.GetIssues(err)

	s.Fatalw(
		chain[0].Message,
		"context", ectx,
		"chain", chain,
		"issues", issues,
	)

	os.Exit(1)
}

/*
TermSink complies with type Sink interface (zapcore.WriteSyncer, io.Closer)

This allows us to add a terminal embedded into the GUI as an additional display for
logs via registering it as a custom sink with zap
*/

type TermSink struct {
	Writer io.Writer
	Reader io.Reader
}

func (t *TermSink) Sync() error  { return nil }
func (t *TermSink) Close() error { return nil }
func (t *TermSink) Write(p []byte) (int, error) {

	// TODO: only write if the terminal is visible??
	// possibly by checking reflect.TypeOf(t.Writer) == reflect.TypeOf(*io.PipeWriter) ?? (pseudo code)

	// unmarshal the bytes into JSON
	var data termData
	var i int
	err := json.Unmarshal(p, &data)
	if err != nil {
		// If unmarshaling fails, write the bytes to the terminal
		i, err = t.Writer.Write(p)
		if err != nil {
			return 0, err
		}
	} else {
		// If unmarshaling succeeds, write the formatted entry to the terminal
		encoded := []byte(data.String())
		i, err = t.Writer.Write(encoded)
		if err != nil {
			return 0, err
		}
	}

	return i, nil
}

func (t *TermSink) Register() error {
	err := zap.RegisterSink("termsink", func(u *url.URL) (zap.Sink, error) {
		return t, nil
	})
	if err != nil {
		return fault.Wrap(err, fmsg.With("Error registering terminal sink"))
	}
	return nil
}

func (t *TermSink) SetIO(r io.Reader, w io.Writer) {
	t.Reader = r
	t.Writer = w
}

func (t *TermSink) SetDiscard() {
	t.Writer = NewDiscardCloser()
	t.Reader = io.NopCloser(nil)
}

func BuildTermSink(w io.Writer) *TermSink {
	return &TermSink{
		Writer: w,
	}
}

type termData struct {
	Time    time.Time     `json:"ts"`
	Level   zapcore.Level `json:"level"`
	Message string        `json:"msg"`
	Issues  []string      `json:"issues"`
}

func (t termData) String() string {
	issues := ""
	for _, i := range t.Issues {
		issues += fmt.Sprintf("\r\n\t\t\t\t%s", i)
	}

	return fmt.Sprintf(
		"[%s] %s %s %s\r\n",
		t.Time.Format("2006-01-02 15:04:05"),
		formatLevel(t.Level),
		t.Message,
		issues,
	)
}

func formatLevel(level zapcore.Level) string {
	switch level {
	case zapcore.DebugLevel:
		return fmt.Sprintf("\u001b[34m%-5s\u001b[0m", level.CapitalString())
	case zapcore.InfoLevel:
		return fmt.Sprintf("\u001b[32m%-5s\u001b[0m", level.CapitalString())
	case zapcore.WarnLevel:
		return fmt.Sprintf("\u001b[33m%-5s\u001b[0m", level.CapitalString())
	case zapcore.ErrorLevel:
		return fmt.Sprintf("\u001b[31m%-5s\u001b[0m", level.CapitalString())
	case zapcore.DPanicLevel:
		return fmt.Sprintf("\u001b[37m%-5s\u001b[0m", level.CapitalString())
	case zapcore.PanicLevel:
		return fmt.Sprintf("\u001b[37m%-5s\u001b[0m", level.CapitalString())
	case zapcore.FatalLevel:
		return fmt.Sprintf("\u001b[37m%-5s\u001b[0m", level.CapitalString())
	default:
		return "UNKNOWN"
	}
}

/*
DiscardCloser is a wrapper around io.Discard that implements the io.Closer interface

This serves as an io.Writer that can be used in place of io.Discard, but also implements
the io.Closer interface
*/
type DiscardCloser struct {
	io.Writer
}

func (d DiscardCloser) Close() error {
	return nil
}

func NewDiscardCloser() *DiscardCloser {
	return &DiscardCloser{
		Writer: io.Discard,
	}
}
