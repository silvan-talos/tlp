package log

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/silvan-talos/tlp"
)

func TestLogger_WithAttrs(t *testing.T) {
	t.Parallel()

	var tests = map[string]struct {
		attrs    []tlp.Attr
		expected []tlp.Attr
	}{
		"nilSliceOfAttrs_onlyCloneTheLogger": {
			attrs:    nil,
			expected: nil,
		},
		"specifiedAttrs_shouldBeInNewLoggerWithoutAffectingBase": {
			attrs: []tlp.Attr{
				{
					Key:   "env",
					Value: "test",
				},
				{
					Key:   "test-type",
					Value: "unit",
				},
			},
			expected: []tlp.Attr{{
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
