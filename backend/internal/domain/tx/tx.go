package tx

import "context"

// TxManager はトランザクションを管理するインターフェースです。
type TxManager interface {
	// Do はトランザクション内で関数を実行します。
	// トランザクション内でエラーが発生した場合はロールバックし、そうでなければコミットします。
	Do(ctx context.Context, fn func(ctx context.Context) (any, error)) (any, error)
}
