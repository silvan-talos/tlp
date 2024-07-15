package log_test

import (
	"testing"

	"github.com/silvan-talos/tlp/log"
)

func TestLevel_String(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		level    log.Level
		expected string
	}{
		"uninitialized": {
			expected: "INFO",
		},
		"debug": {
			level:    log.LevelDebug,
			expected: "DEBUG",
		},
		"info": {
			level:    log.LevelInfo,
			expected: "INFO",
		},
		"warn": {
			level:    log.LevelWarn,
			expected: "WARN",
		},
		"error": {
			level:    log.LevelError,
			expected: "ERROR",
		},
		"trace": {
			level:    log.Level(-8),
			expected: "LEVEL(-8)",
		},
		"betweenExistingLevels": {
			level:    log.Level(2),
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
