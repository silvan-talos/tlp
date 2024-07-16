package log

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/silvan-talos/tlp/apm"
	"github.com/silvan-talos/tlp/config"
	"github.com/silvan-talos/tlp/logging"
	"github.com/silvan-talos/tlp/text"
	"github.com/silvan-talos/tlp/transaction"
)

var defaultLogger atomic.Pointer[Logger]

func init() {
	transaction.SetDefaultTracer(transaction.NewTracer(apm.NewRecorder()))
	var cfg config.Config
	err := config.LoadFromYAML("config.yml", &cfg)
	if err != nil {
		defaultLogger.Store(NewLogger(text.NewDriver(nil), logging.LevelInfo))
		return
	}
}

type Driver interface {
	Log(ctx context.Context, entry logging.Entry)
}

type Logger struct {
	driver Driver
	level  logging.Level
	attrs  []logging.Attr
}

func NewLogger(driver Driver, level logging.Level) *Logger {
	return &Logger{
		driver: driver,
		level:  level,
	}
}

func SetDefault(l *Logger) {
	defaultLogger.Store(l)
}

func Default() *Logger {
	return defaultLogger.Load()
}

func (l *Logger) Log(ctx context.Context, level logging.Level, msg string, args ...any) {
	if level < l.level {
		return
	}
	entry := logging.Entry{
		Time:    time.Now(),
		Message: msg,
		Level:   level,
	}
	entry.Attrs = l.attrs
	tx := transaction.FromContext(ctx)
	entry.TraceID = tx.TraceID
	entry.TransactionAttrs = tx.Attrs
	for i := 0; i < len(args); i += 2 {
		if key, ok := args[i].(string); ok && i+1 < len(args) {
			entry.Attrs = append(entry.Attrs, logging.NewAttr(key, args[i+1]))
			continue
		}
		entry.Attrs = append(entry.Attrs, logging.NewAttr("undefKey", args[i]))
		// move i backwards since we only processed one arg
		i--
	}
	l.driver.Log(ctx, entry)
}

func (l *Logger) WithAttrs(attrs ...logging.Attr) *Logger {
	clone := *l
	clone.attrs = append(clone.attrs, attrs...)
	return &clone
}

func (l *Logger) WithLevel(level string) (*Logger, error) {
	lvl, err := logging.ParseLevel(level)
	if err != nil {
		return nil, err
	}
	clone := *l
	clone.level = lvl
	return &clone, nil
}

func (l *Logger) Debug(ctx context.Context, msg string, args ...any) {
	l.Log(ctx, logging.LevelDebug, msg, args...)
}

func (l *Logger) Info(ctx context.Context, msg string, args ...any) {
	l.Log(ctx, logging.LevelInfo, msg, args...)
}

func (l *Logger) Warn(ctx context.Context, msg string, args ...any) {
	l.Log(ctx, logging.LevelWarn, msg, args...)
}

func (l *Logger) Error(ctx context.Context, msg string, args ...any) {
	l.Log(ctx, logging.LevelError, msg, args...)
}

func Debug(ctx context.Context, msg string, args ...any) {
	Default().Log(ctx, logging.LevelDebug, msg, args...)
}

func Info(ctx context.Context, msg string, args ...any) {
	Default().Log(ctx, logging.LevelInfo, msg, args...)
}

func Warn(ctx context.Context, msg string, args ...any) {
	Default().Log(ctx, logging.LevelWarn, msg, args...)
}

func Error(ctx context.Context, msg string, args ...any) {
	Default().Log(ctx, logging.LevelError, msg, args...)
}
