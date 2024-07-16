package apm

import (
	"context"

	"go.elastic.co/apm/v2"

	"github.com/silvan-talos/tlp/transaction"
)

type Recorder struct {
	// create specific APM tracer instead of using the default
}

func NewRecorder() *Recorder {
	return &Recorder{}
}

func (r *Recorder) RecordTransaction(ctx context.Context, name, transactionType string) (*transaction.Transaction, context.Context) {
	tx := apm.DefaultTracer().StartTransaction(name, transactionType)
	return &transaction.Transaction{
		TraceID: tx.TraceContext().Trace.String(),
	}, ctx
}
