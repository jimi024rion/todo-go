package cache

import (
	"context"

	"github.com/jimi024rion/todo-go/backend/internal/domain/apikey/model/entity"
)

// APIKeyCache はAPIキーのキャッシュインターフェースです。
// ローカルキャッシュ・Valkeyなど複数の実装を差し替え可能にします。
type APIKeyCache interface {
	// Get はrawKeyに対応するAPIKeyを取得します。見つからない場合はfalseを返します。
	Get(ctx context.Context, key string) (*entity.APIKey, bool)
	// Set はrawKeyに対応するAPIKeyをキャッシュします。
	Set(ctx context.Context, key string, apiKey *entity.APIKey)
	// Delete はrawKeyに対応するキャッシュエントリを削除します。
	Delete(ctx context.Context, key string)
}
