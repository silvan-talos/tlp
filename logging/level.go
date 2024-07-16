package logging

import (
	"fmt"
	"strconv"
	"strings"
)

// A Level represents the severity of a log event.
type Level int

const (
	LevelDebug Level = -4
	LevelInfo  Level = 0
	LevelWarn  Level = 4
	LevelError Level = 8
)

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	}
	return fmt.Sprintf("LEVEL(%d)", int(l))
}

func ParseLevel(l string) (Level, error) {
	l = strings.ToUpper(l)
	switch l {
	case "DEBUG":
		return LevelDebug, nil
	case "INFO":
		return LevelInfo, nil
	case "WARN":
		return LevelWarn, nil
	case "ERROR":
		return LevelError, nil
	}
	level := strings.TrimPrefix(l, "LEVEL(")
	level = strings.TrimSuffix(level, ")")
	severity, err := strconv.Atoi(level)
	if err != nil {
		return LevelDebug, fmt.Errorf("unparsable LEVEL: %s", l)
	}
	return Level(severity), nil
}
