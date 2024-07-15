package log_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/silvan-talos/tlp"
	"github.com/silvan-talos/tlp/log"
	"github.com/silvan-talos/tlp/mock"
	"github.com/silvan-talos/tlp/transaction"
)

func TestLogger_Log(t *testing.T) {
	t.Run("correctly ignore levels", checkLevelFunctionality)
	t.Run("logger with attrs", checkLoggerAttributes)
	t.Run("logger with args and transaction details", checkLoggerArgsAndTransaction)
	t.Run("arg with undefined key", argWithUndefinedKey)
}

func checkLevelFunctionality(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		setLogLevel   log.Level
		testLevel     log.Level
		expectedToLog bool
	}{
		"debug_debug_log": {
			setLogLevel:   log.LevelDebug,
			testLevel:     log.LevelDebug,
			expectedToLog: true,
		},
		"debug_info_log": {
			setLogLevel:   log.LevelDebug,
			testLevel:     log.LevelInfo,
			expectedToLog: true,
		},
		"debug_trace_skipLogging": {
			setLogLevel:   log.LevelDebug,
			testLevel:     log.Level(-5),
			expectedToLog: false,
		},
		"error_debug_skipLogging": {
			setLogLevel:   log.LevelError,
			testLevel:     log.LevelDebug,
			expectedToLog: false,
		},
		"error_info_skipLogging": {
			setLogLevel:   log.LevelError,
			testLevel:     log.LevelInfo,
			expectedToLog: false,
		},
		"error_error_log": {
			setLogLevel:   log.LevelError,
			testLevel:     log.LevelError,
			expectedToLog: true,
		},
		"warn_info_skipLogging": {
			setLogLevel:   log.LevelWarn,
			testLevel:     log.LevelInfo,
			expectedToLog: false,
		},
		"warn_error_log": {
			setLogLevel:   log.LevelWarn,
			testLevel:     log.LevelError,
			expectedToLog: true,
		},
	}
	for testName, tc := range tests {
		t.Run(testName, func(t *testing.T) {
			driver := &mock.Driver{}
			logger := log.NewLogger(driver, tc.setLogLevel)
			logger.Log(context.Background(), tc.testLevel, "test message to be logged")
			if tc.expectedToLog && driver.Count == 0 {
				t.Error("was expected to log")
			}
			if !tc.expectedToLog && driver.Count > 0 {
				t.Error("was expected to not log")
			}
		})
	}
}

func checkLoggerAttributes(t *testing.T) {
	t.Parallel()

	attrs := []tlp.Attr{
		tlp.NewAttr("env", "dev"),
		tlp.NewAttr("test-type", "unit"),
	}
	driver := &mock.Driver{
		LogFn: func(ctx context.Context, entry log.Entry) {
			// append traceID and args attributes
			attrs = append(attrs, tlp.NewAttr("TraceID", ""), tlp.NewAttr("reason", "test"))
			require.Equal(t, entry.Attrs, attrs, "logger attributes should be passed")
		},
	}
	logger := log.NewLogger(driver, log.LevelDebug).WithAttrs(attrs...)
	logger.Log(context.Background(), log.LevelInfo, "test message to be logged", "reason", "test")
}

func checkLoggerArgsAndTransaction(t *testing.T) {
	t.Parallel()

	attrs := []tlp.Attr{
		tlp.NewAttr("env", "dev"),
		tlp.NewAttr("test-type", "unit"),
	}
	driver := &mock.Driver{
		LogFn: func(ctx context.Context, entry log.Entry) {
			// append traceID and args attributes
			attrs = append(attrs, tlp.NewAttr("TraceID", "test-trace"), tlp.NewAttr("reason", "test"))
			require.ElementsMatch(t, entry.Attrs, attrs, "logger attributes should be passed")
		},
	}
	tracer := transaction.NewTracer(&mock.TransactionRecorder{
		RecordTransactionFn: func(ctx context.Context, name, transactionType string) (*transaction.Transaction, context.Context) {
			return &transaction.Transaction{
				TraceID: "test-trace",
				Attrs:   nil,
			}, ctx
		},
	})
	tx, ctx := tracer.StartTransaction(context.Background(), "test", "context-test", attrs...)
	defer tx.End()
	logger := log.NewLogger(driver, log.LevelDebug)
	logger.Log(ctx, log.LevelInfo, "test message to be logged", "reason", "test")
}

func argWithUndefinedKey(t *testing.T) {
	t.Parallel()

	attrs := []tlp.Attr{
		tlp.NewAttr("TraceID", ""),
		tlp.NewAttr("reason", "test"),
		tlp.NewAttr("undefKey", 3),
	}
	driver := &mock.Driver{
		LogFn: func(ctx context.Context, entry log.Entry) {
			// append traceID and args attributes
			require.ElementsMatch(t, entry.Attrs, attrs, "logger attributes should be passed")
		},
	}
	logger := log.NewLogger(driver, log.LevelDebug)
	logger.Log(context.Background(), log.LevelInfo, "test message to be logged", "reason", "test", 3)
}
