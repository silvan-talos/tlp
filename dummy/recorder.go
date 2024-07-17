package dummy

import (
	"context"

	"github.com/google/uuid"

	"github.com/silvan-talos/tlp/transaction"
)

type Recorder struct {
	// dummy recorder that sets an uuid as TraceID
}

func NewRecorder() *Recorder {
	return &Recorder{}
}

func (r *Recorder) RecordTransaction(ctx context.Context, name, transactionType string) (*transaction.Transaction, context.Context) {
	return &transaction.Transaction{
		TraceID: uuid.New().String(),
	}, ctx
}
