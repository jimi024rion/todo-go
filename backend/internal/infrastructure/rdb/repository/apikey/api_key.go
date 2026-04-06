package apikey

import (
	"context"
	"fmt"

	"github.com/aarondl/opt/omit"
	"github.com/gofrs/uuid/v5"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/dm"
	"github.com/stephenafamo/bob/dialect/psql/sm"

	apikeyentity "github.com/jimi024rion/todo-go/backend/internal/domain/apikey/model/entity"
	apikeyvo "github.com/jimi024rion/todo-go/backend/internal/domain/apikey/model/valueobject"
	apikeyrepo "github.com/jimi024rion/todo-go/backend/internal/domain/apikey/repository"
	"github.com/jimi024rion/todo-go/backend/internal/infrastructure/rdb"
	"github.com/jimi024rion/todo-go/backend/internal/infrastructure/rdb/models"
)

type Repository struct {
	db bob.DB
}

func NewRepository(db bob.DB) apikeyrepo.APIKeyRepository {
	return &Repository{db: db}
}

func (r *Repository) getExecutor(ctx context.Context) bob.Executor {
	if exec, ok := rdb.GetExecutor(ctx); ok {
		return exec
	}
	return r.db
}

func (r *Repository) Save(ctx context.Context, apiKey *apikeyentity.APIKey) error {
	uid, err := uuid.FromString(apiKey.ID().String())
	if err != nil {
		return fmt.Errorf("invalid api key id format: %w", err)
	}
	userUID, err := uuid.FromString(apiKey.UserID())
	if err != nil {
		return fmt.Errorf("invalid user id format: %w", err)
	}

	setter := &models.APIKeySetter{
		ID:        omit.From(uid),
		UserID:    omit.From(userUID),
		KeyHash:   omit.From(apiKey.KeyHash()),
		Name:      omit.From(apiKey.Name()),
		CreatedAt: omit.From(apiKey.CreatedAt()),
	}

	_, err = models.APIKeys.Insert(setter).Exec(ctx, r.getExecutor(ctx))
	if err != nil {
		return fmt.Errorf("failed to insert api key: %w", err)
	}

	return nil
}

func (r *Repository) FindByKeyHash(ctx context.Context, keyHash string) (*apikeyentity.APIKey, error) {
	m, err := models.APIKeys.Query(
		sm.Where(models.APIKeys.Columns.KeyHash.EQ(psql.Arg(keyHash))),
	).One(ctx, r.getExecutor(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to find api key by hash: %w", err)
	}

	ak, err := apikeyentity.NewAPIKey(
		m.ID.String(),
		m.UserID.String(),
		m.KeyHash,
		m.Name,
		m.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to reconstruct api key entity: %w", err)
	}

	return ak, nil
}

func (r *Repository) Delete(ctx context.Context, id apikeyvo.APIKeyID) error {
	uid, err := uuid.FromString(id.String())
	if err != nil {
		return fmt.Errorf("invalid api key id format: %w", err)
	}

	_, err = models.APIKeys.Delete(
		dm.Where(models.APIKeys.Columns.ID.EQ(psql.Arg(uid))),
	).Exec(ctx, r.getExecutor(ctx))
	if err != nil {
		return fmt.Errorf("failed to delete api key: %w", err)
	}

	return nil
}
