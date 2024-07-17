package log

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/silvan-talos/tlp/apm"
	"github.com/silvan-talos/tlp/config"
	"github.com/silvan-talos/tlp/dummy"
	"github.com/silvan-talos/tlp/json"
	"github.com/silvan-talos/tlp/logging"
	"github.com/silvan-talos/tlp/text"
	"github.com/silvan-talos/tlp/transaction"
)

var defaultLogger atomic.Pointer[Logger]

func init() {
	var cfg config.Config
	err := config.LoadFromYAML("log-config.yml", &cfg)
	if err != nil {
		fmt.Println("load config error", "err", err)
		transaction.SetDefaultTracer(transaction.NewTracer(dummy.NewRecorder()))
		defaultLogger.Store(NewLogger(text.NewDriver(nil), logging.LevelInfo))
		return
	}
	interpretConfig(cfg)
}

func interpretConfig(cfg config.Config) {
	var recorder transaction.Recorder
	if cfg.Transaction.RecorderType != "" && strings.EqualFold(cfg.Transaction.RecorderType, "apm") {
		recorder = apm.NewRecorder()
	} else {
		recorder = dummy.NewRecorder()
	}
	transaction.SetDefaultTracer(transaction.NewTracer(recorder))
	defaultLogger.Store(NewLoggerFromConfig(cfg.Log))
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

func NewLoggerFromConfig(cfg config.LogConfig) *Logger {
	output := os.Stdout
	if cfg.OutputFile != "" {
		f, err := os.OpenFile(cfg.OutputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			output = f
		} else {
			fmt.Println("open file", err)
		}
	}
	var driver Driver
	switch cfg.ProcessingType {
	case "json":
		driver = json.NewDriver(output)
	default:
		driver = text.NewDriver(output)
	}
	lvl := logging.LevelInfo
	if cfg.Level != "" {
		if l, err := logging.ParseLevel(cfg.Level); err == nil {
			lvl = l
		}
	}
	logger := NewLogger(driver, lvl)
	if cfg.PermanentAttributes != nil {
		attrs := make([]logging.Attr, 0, 1)
		for _, item := range cfg.PermanentAttributes {
			for k, v := range item {
				attrs = append(attrs, logging.Attr{Key: k, Value: v})
			}
		}
		logger = logger.WithAttrs(attrs...)
	}
	return logger
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
