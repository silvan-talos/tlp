package transaction

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/silvan-talos/tlp/logging"
)

var defaultTracer atomic.Pointer[Tracer]

type Transaction struct {
	TraceID string
	Attrs   []logging.Attr

	start    time.Time
	duration time.Duration
}

type transactionKey struct{}

func FromContext(ctx context.Context) *Transaction {
	tx, ok := ctx.Value(transactionKey{}).(*Transaction)
	if !ok {
		return &Transaction{}
	}
	return tx
}

func NewContext(ctx context.Context, tx *Transaction) context.Context {
	return context.WithValue(ctx, transactionKey{}, tx)
}

type Recorder interface {
	RecordTransaction(ctx context.Context, name, transactionType string) (*Transaction, context.Context)
}

type Tracer struct {
	recorder Recorder
}

func DefaultTracer() *Tracer {
	return defaultTracer.Load()
}

func SetDefaultTracer(t *Tracer) {
	defaultTracer.Store(t)
}

func NewTracer(recorder Recorder) *Tracer {
	return &Tracer{recorder: recorder}
}

func (t *Tracer) StartTransaction(ctx context.Context, name, transactionType string, attrs ...logging.Attr) (*Transaction, context.Context) {
	tx, ctx := t.recorder.RecordTransaction(ctx, name, transactionType)
	tx.Attrs = append(tx.Attrs, attrs...)
	tx.start = time.Now()
	ctx = NewContext(ctx, tx)
	return tx, ctx
}

func (t *Transaction) End() {
	t.duration = time.Now().Sub(t.start)
}

func (t *Transaction) GetDuration() time.Duration {
	return t.duration
}
