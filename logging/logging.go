// Package logging contains logging helpers like level, attribute and entry and functions related to them.
package logging

import (
	"time"
)

type Attr struct {
	Key   string
	Value any
}

func NewAttr(name string, value any) Attr {
	return Attr{Key: name, Value: value}
}

type Entry struct {
	Time             time.Time
	Message          string
	Level            Level
	Attrs            []Attr
	TraceID          string
	TransactionAttrs []Attr
}
