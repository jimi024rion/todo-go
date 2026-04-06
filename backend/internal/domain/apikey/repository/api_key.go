package repository

import (
	"context"

	"github.com/jimi024rion/todo-go/backend/internal/domain/apikey/model/entity"
	"github.com/jimi024rion/todo-go/backend/internal/domain/apikey/model/valueobject"
)

type APIKeyRepository interface {
	Save(ctx context.Context, apiKey *entity.APIKey) error
	FindByKeyHash(ctx context.Context, keyHash string) (*entity.APIKey, error)
	Delete(ctx context.Context, id valueobject.APIKeyID) error
}
