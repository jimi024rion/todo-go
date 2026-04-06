package user

import (
	"context"
	"fmt"

	"github.com/aarondl/opt/omit"
	"github.com/gofrs/uuid/v5"
	"github.com/stephenafamo/bob"

	userentity "github.com/jimi024rion/todo-go/backend/internal/domain/user/model/entity"
	userrepo "github.com/jimi024rion/todo-go/backend/internal/domain/user/repository"
	"github.com/jimi024rion/todo-go/backend/internal/infrastructure/rdb"
	"github.com/jimi024rion/todo-go/backend/internal/infrastructure/rdb/models"
)

type Repository struct {
	db bob.DB
}

func NewRepository(db bob.DB) userrepo.UserRepository {
	return &Repository{db: db}
}

func (r *Repository) getExecutor(ctx context.Context) bob.Executor {
	if exec, ok := rdb.GetExecutor(ctx); ok {
		return exec
	}
	return r.db
}

func (r *Repository) Save(ctx context.Context, user *userentity.User) error {
	uid, err := uuid.FromString(user.ID().String())
	if err != nil {
		return fmt.Errorf("invalid user id format: %w", err)
	}

	setter := &models.UserSetter{
		ID:        omit.From(uid),
		Name:      omit.From(user.Name()),
		Email:     omit.From(user.Email()),
		CreatedAt: omit.From(user.CreatedAt()),
		UpdatedAt: omit.From(user.UpdatedAt()),
	}

	_, err = models.Users.Insert(setter).Exec(ctx, r.getExecutor(ctx))
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	return nil
}
