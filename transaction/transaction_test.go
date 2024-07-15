package transaction_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/silvan-talos/tlp"
	"github.com/silvan-talos/tlp/mock"
	"github.com/silvan-talos/tlp/transaction"
)

func TestTracer_StartTransaction(t *testing.T) {
	t.Run("must succeed", startTransactionSuccessfully)
}

func startTransactionSuccessfully(t *testing.T) {
	t.Parallel()

	tracer := transaction.NewTracer(&mock.TransactionRecorder{})
	tx, ctx := tracer.StartTransaction(context.Background(), "test", "unit-test", tlp.NewAttr("env", "test"))
	require.NotNil(t, tx, "expected transaction to be not nil")
	require.Contains(t, tx.Attrs, tlp.NewAttr("env", "test"), "returned transaction should contain the custom attrs")
	ctxTx := transaction.FromContext(ctx)
	require.Equal(t, tx, ctxTx, "retrieved transaction should be the same")
	tx.End()
	require.NotEmpty(t, tx.GetDuration(), "duration should be different from 0")
}

func TestFromContext(t *testing.T) {
	t.Run("retrieve existing transaction", retrieveExistingTransaction)
	t.Run("return empty transaction if none available", retrieveWhenNoneInContext)
}

func retrieveExistingTransaction(t *testing.T) {
	t.Parallel()

	tracer := transaction.NewTracer(&mock.TransactionRecorder{})
	tx, ctx := tracer.StartTransaction(context.Background(), "test", "unit-test")
	require.NotNil(t, tx, "expected transaction to be not nil")
	ctxTx := transaction.FromContext(ctx)
	require.Equal(t, tx, ctxTx, "retrieved transaction should be the same")
}

func retrieveWhenNoneInContext(t *testing.T) {
	t.Parallel()

	ctxTx := transaction.FromContext(context.Background())
	require.NotNil(t, ctxTx, "transaction must not be nil")
}
