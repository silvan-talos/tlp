package mock

import (
	"context"

	"github.com/silvan-talos/tlp/logging"
	"github.com/silvan-talos/tlp/transaction"
)

type TransactionRecorder struct {
	RecordTransactionFn func(ctx context.Context, name, transactionType string) (*transaction.Transaction, context.Context)
}

func (tr *TransactionRecorder) RecordTransaction(ctx context.Context, name, transactionType string) (*transaction.Transaction, context.Context) {
	if tr.RecordTransactionFn != nil {
		return tr.RecordTransactionFn(ctx, name, transactionType)
	}
	return &transaction.Transaction{
		TraceID: "test-trace",
		Attrs: []logging.Attr{
			{
				Key:   "name",
				Value: name,
			},
			{
				Key:   "type",
				Value: transactionType,
			},
		},
	}, context.WithValue(ctx, "env", "test")
}
