package rdb

import (
	"context"

	"github.com/jimi024rion/todo-go/backend/internal/domain/tx"
	"github.com/stephenafamo/bob"
)

type executorKey struct{}

type txManager struct {
	db bob.DB
}

// NewTxManager は新しいTxManagerを生成します。
func NewTxManager(db bob.DB) tx.TxManager {
	return &txManager{db: db}
}

func (tm *txManager) Do(ctx context.Context, fn func(ctx context.Context) (any, error)) (_ any, err error) {
	// トランザクションを開始
	transaction, err := tm.db.Begin(ctx)
	if err != nil {
		return nil, err
	}

	// panicが発生した場合でもトランザクションをロールバックする
	defer func() {
		if p := recover(); p != nil {
			_ = transaction.Rollback(ctx)
			panic(p)
		} else if err != nil {
			_ = transaction.Rollback(ctx)
		} else {
			err = transaction.Commit(ctx)
		}
	}()

	// トランザクション（Executor）をコンテキストに含める
	execCtx := context.WithValue(ctx, executorKey{}, transaction)
	result, err := fn(execCtx)
	return result, err
}

// GetExecutor はコンテキストからExecutorを取得します。
func GetExecutor(ctx context.Context) (bob.Executor, bool) {
	exec, ok := ctx.Value(executorKey{}).(bob.Executor)
	return exec, ok
}
