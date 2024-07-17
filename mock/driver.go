// Package mock contains mocks used in unit tests.
package mock

import (
	"context"

	"github.com/silvan-talos/tlp/logging"
)

type Driver struct {
	LogFn func(ctx context.Context, entry logging.Entry)

	Count int
}

func (d *Driver) Log(ctx context.Context, entry logging.Entry) {
	d.Count++
	if d.LogFn != nil {
		d.LogFn(ctx, entry)
	}
}
