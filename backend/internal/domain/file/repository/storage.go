package repository

import (
	"context"
	"time"
)

// StorageRepository はファイルストレージ操作を担うリポジトリインターフェースです。
type StorageRepository interface {
	// GetDownloadURL はオブジェクトへの署名付きダウンロードURLを返します。
	GetDownloadURL(ctx context.Context, objectName string, expiry time.Duration) (string, error)
}
