// Package dummy is a `Recorder` interface implementation that simply generates a uuid TraceID.
package dummy

import (
	"context"

	"github.com/google/uuid"

	"github.com/silvan-talos/tlp/transaction"
)

// Recorder is a dummy implementation of a trace recorder.
// It simply sets an uuid as TraceID.
type Recorder struct{}

func NewRecorder() *Recorder {
	return &Recorder{}
}

func (r *Recorder) RecordTransaction(ctx context.Context, name, transactionType string) (*transaction.Transaction, context.Context) {
	return &transaction.Transaction{
		TraceID: uuid.New().String(),
	}, ctx
}
