package log

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/silvan-talos/tlp"
	"github.com/silvan-talos/tlp/config"
	"github.com/silvan-talos/tlp/transaction"
)

var defaultLogger atomic.Pointer[Logger]

func init() {
	defaultLogger.Store(NewLogger(nil, LevelInfo))
}

type Driver interface {
	Log(ctx context.Context, entry Entry)
}

type Entry struct {
	Time             time.Time
	Message          string
	Level            Level
	Attrs            []tlp.Attr
	TraceID          string
	TransactionAttrs []tlp.Attr
}

type Logger struct {
	driver Driver
	level  Level
	attrs  []tlp.Attr
}

func NewLogger(driver Driver, level Level) *Logger {
	return &Logger{
		driver: driver,
		level:  level,
	}
}

func NewLoggerWithConfig(config config.Config) *Logger {
	return &Logger{}
}

func SetDefault(l *Logger) {
	defaultLogger.Store(l)
}

func Default() *Logger {
	return defaultLogger.Load()
}

func (l *Logger) Log(ctx context.Context, level Level, msg string, args ...any) {
	if level < l.level {
		return
	}
	entry := Entry{
		Time:    time.Now(),
		Message: msg,
		Level:   level,
	}
	entry.Attrs = l.attrs
	tx := transaction.FromContext(ctx)
	entry.Attrs = append(entry.Attrs, tlp.NewAttr("TraceID", tx.TraceID))
	entry.Attrs = append(entry.Attrs, tx.Attrs...)
	for i := 0; i < len(args); i += 2 {
		if key, ok := args[i].(string); ok {
			entry.Attrs = append(entry.Attrs, tlp.NewAttr(key, args[i+1]))
			continue
		}
		entry.Attrs = append(entry.Attrs, tlp.NewAttr("undefKey", args[i]))
	}
	l.driver.Log(ctx, entry)
}

func (l *Logger) WithAttrs(attrs ...tlp.Attr) *Logger {
	clone := *l
	clone.attrs = append(clone.attrs, attrs...)
	return &clone
}
