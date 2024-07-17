package json

import (
	"bufio"
	"context"
	stdjson "encoding/json"
	"io"
	"os"

	"github.com/silvan-talos/tlp/logging"
)

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
	_ = stdjson.NewEncoder(d.writer).Encode(entry)
	_ = d.writer.Flush()
}
