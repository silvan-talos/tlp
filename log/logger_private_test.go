package log

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/silvan-talos/tlp/logging"
)

func TestLogger_WithAttrs(t *testing.T) {
	t.Parallel()

	var tests = map[string]struct {
		attrs    []logging.Attr
		expected []logging.Attr
	}{
		"nilSliceOfAttrs_onlyCloneTheLogger": {
			attrs:    nil,
			expected: nil,
		},
		"specifiedAttrs_shouldBeInNewLoggerWithoutAffectingBase": {
			attrs: []logging.Attr{
				{
					Key:   "env",
					Value: "test",
				},
				{
					Key:   "test-type",
					Value: "unit",
				},
			},
			expected: []logging.Attr{{
				Key:   "env",
				Value: "test",
			},
				{
					Key:   "test-type",
					Value: "unit",
				},
			},
		},
	}
	for testName, tc := range tests {
		t.Run(testName, func(t *testing.T) {
			defLogger := Default()
			newLogger := defLogger.WithAttrs(tc.attrs...)
			require.ElementsMatch(t, newLogger.attrs, tc.expected, "logger attrs should match")
			require.Empty(t, defLogger.attrs, "default logger should not be affected")
		})
	}
}

func TestLogger_WithLevel(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		level       string
		want        logging.Level
		shouldError bool
	}{
		{"debug", logging.LevelDebug, false},
		{"info", logging.LevelInfo, false},
		{"warn", logging.LevelWarn, false},
		{"ERROR", logging.LevelError, false},
		{"LEVEL(-2)", logging.Level(-2), false},
		{"LEVEL(9)", logging.Level(9), false},
		{"LEVEL(3)", logging.Level(3), false},
		{"LeVeL(-7)", logging.Level(-7), false},
		{"LeVeLl(-7)", logging.LevelDebug, true},
		{"LEVEL-7", logging.LevelDebug, true},
		{"UNPARSABLE", logging.LevelDebug, true},
	} {
		logger, err := Default().WithLevel(tc.level)
		if !tc.shouldError {
			require.NoError(t, err, "test case shouldn't error")
			require.Equal(t, tc.want, logger.level, "logger should have the correct level")
		} else {
			require.Error(t, err, "test case should return error")
		}
	}
}
