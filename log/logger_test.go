package log_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/silvan-talos/tlp/log"
	"github.com/silvan-talos/tlp/logging"
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
		setLogLevel   logging.Level
		testLevel     logging.Level
		expectedToLog bool
	}{
		"debug_debug_log": {
			setLogLevel:   logging.LevelDebug,
			testLevel:     logging.LevelDebug,
			expectedToLog: true,
		},
		"debug_info_log": {
			setLogLevel:   logging.LevelDebug,
			testLevel:     logging.LevelInfo,
			expectedToLog: true,
		},
		"debug_trace_skipLogging": {
			setLogLevel:   logging.LevelDebug,
			testLevel:     logging.Level(-5),
			expectedToLog: false,
		},
		"error_debug_skipLogging": {
			setLogLevel:   logging.LevelError,
			testLevel:     logging.LevelDebug,
			expectedToLog: false,
		},
		"error_info_skipLogging": {
			setLogLevel:   logging.LevelError,
			testLevel:     logging.LevelInfo,
			expectedToLog: false,
		},
		"error_error_log": {
			setLogLevel:   logging.LevelError,
			testLevel:     logging.LevelError,
			expectedToLog: true,
		},
		"warn_info_skipLogging": {
			setLogLevel:   logging.LevelWarn,
			testLevel:     logging.LevelInfo,
			expectedToLog: false,
		},
		"warn_error_log": {
			setLogLevel:   logging.LevelWarn,
			testLevel:     logging.LevelError,
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

	attrs := []logging.Attr{
		logging.NewAttr("env", "dev"),
		logging.NewAttr("test-type", "unit"),
	}
	driver := &mock.Driver{
		LogFn: func(ctx context.Context, entry logging.Entry) {
			// append traceID and args attributes
			attrs = append(attrs, logging.NewAttr("reason", "test"))
			require.Equal(t, entry.Attrs, attrs, "logger attributes should be passed")
		},
	}
	logger := log.NewLogger(driver, logging.LevelDebug).WithAttrs(attrs...)
	logger.Log(context.Background(), logging.LevelInfo, "test message to be logged", "reason", "test")
}

func checkLoggerArgsAndTransaction(t *testing.T) {
	t.Parallel()

	attrs := []logging.Attr{
		logging.NewAttr("env", "dev"),
		logging.NewAttr("test-type", "unit"),
	}
	driver := &mock.Driver{
		LogFn: func(ctx context.Context, entry logging.Entry) {
			require.ElementsMatch(t, entry.TransactionAttrs, attrs, "logger attributes should be passed")
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
	logger := log.NewLogger(driver, logging.LevelDebug)
	logger.Log(ctx, logging.LevelInfo, "test message to be logged", "reason", "test")
}

func argWithUndefinedKey(t *testing.T) {
	t.Parallel()

	attrs := []logging.Attr{
		logging.NewAttr("reason", "test"),
		logging.NewAttr("undefKey", 3),
		logging.NewAttr("undefKey", "string test"),
	}
	driver := &mock.Driver{
		LogFn: func(ctx context.Context, entry logging.Entry) {
			// append traceID and args attributes
			require.ElementsMatch(t, entry.Attrs, attrs, "logger attributes should be passed")
		},
	}
	logger := log.NewLogger(driver, logging.LevelDebug)
	logger.Log(context.Background(), logging.LevelInfo, "test message to be logged", "reason", "test", 3, "string test")
}
