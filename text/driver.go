// Package text provides a plain-text driver that can be used on any io.Writer output (file, stdout etc.).
package text

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/silvan-talos/tlp/logging"
)

const dateFormat = "2006-01-02 15:04:05.000"

type Driver struct {
	writer *bufio.Writer
}

func NewDriver(output io.Writer) *Driver {
	if output == nil {
		output = os.Stdout
	}
	return &Driver{
		writer: bufio.NewWriter(output),
	}
}

func (d *Driver) Log(ctx context.Context, entry logging.Entry) {
	// log format times - LEVEL: msg	traceID=123 details=[key1='value 1', composed-key='value 2'] transactionDetails=[userID='123', requestPath='/users/1/details']
	_, _ = fmt.Fprintf(d.writer, "%s - %s: %s",
		entry.Time.Format(dateFormat),
		entry.Level,
		entry.Message,
	)
	if entry.TraceID != "" {
		_, _ = fmt.Fprintf(d.writer, "\ttraceID=%s", entry.TraceID)
	}
	if len(entry.Attrs) > 0 {
		_, _ = fmt.Fprintf(d.writer, " details=[%s]", textFormatAttrs(entry.Attrs))
	}
	if len(entry.TransactionAttrs) > 0 {
		_, _ = fmt.Fprintf(d.writer, " transactionDetails=[%s]", textFormatAttrs(entry.TransactionAttrs))
	}
	_ = d.writer.WriteByte('\n')
	_ = d.writer.Flush()
}

func textFormatAttrs(attrs []logging.Attr) string {
	parts := make([]string, len(attrs))
	for i, attr := range attrs {
		parts[i] = fmt.Sprintf("%s='%v'", attr.Key, attr.Value)
	}
	return strings.Join(parts, ", ")
}
