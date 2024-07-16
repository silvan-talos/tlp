package logging_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/silvan-talos/tlp/logging"
)

func TestLevel_String(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		level    logging.Level
		expected string
	}{
		"uninitialized": {
			expected: "INFO",
		},
		"debug": {
			level:    logging.LevelDebug,
			expected: "DEBUG",
		},
		"info": {
			level:    logging.LevelInfo,
			expected: "INFO",
		},
		"warn": {
			level:    logging.LevelWarn,
			expected: "WARN",
		},
		"error": {
			level:    logging.LevelError,
			expected: "ERROR",
		},
		"trace": {
			level:    logging.Level(-8),
			expected: "LEVEL(-8)",
		},
		"betweenExistingLevels": {
			level:    logging.Level(2),
			expected: "LEVEL(2)",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if actual := tc.level.String(); tc.expected != actual {
				t.Errorf("expected %s, got %s", tc.expected, actual)
			}
		})
	}
}

func TestParseLevel(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		level     string
		want      logging.Level
		shouldErr bool
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
		lvl, err := logging.ParseLevel(tc.level)
		if tc.shouldErr {
			require.Error(t, err, "test case should return an error")
		} else {
			require.NoError(t, err, "test case should not return an error")
		}
		require.Equal(t, tc.want, lvl, "level should be set correctly")
	}
}
