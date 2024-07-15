package mock

import (
	"context"

	"github.com/silvan-talos/tlp/log"
)

type Driver struct {
	LogFn func(ctx context.Context, entry log.Entry)

	Count int
}

func (d *Driver) Log(ctx context.Context, entry log.Entry) {
	d.Count++
	if d.LogFn != nil {
		d.LogFn(ctx, entry)
	}
}
