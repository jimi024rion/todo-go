package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/google/uuid"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/sm"

	"github.com/jimi024rion/todo-go/backend/internal/config/errs"
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
	uid, err := uuid.Parse(user.ID().String())
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

	if uid := user.FirebaseUID(); uid != "" {
		setter.FirebaseUID = omitnull.From(uid)
	}

	_, err = models.Users.Insert(setter).Exec(ctx, r.getExecutor(ctx))
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	return nil
}

func (r *Repository) FindByFirebaseUID(ctx context.Context, firebaseUID string) (*userentity.User, error) {
	m, err := models.Users.Query(
		sm.Where(models.Users.Columns.FirebaseUID.EQ(psql.Arg(firebaseUID))),
	).One(ctx, r.getExecutor(ctx))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewErr(errs.InternalCodeNotFound, fmt.Errorf("user not found: firebase_uid=%s", firebaseUID))
		}
		return nil, fmt.Errorf("failed to find user by firebase_uid: %w", err)
	}

	return userentity.NewUserWithFirebaseUID(
		m.ID.String(),
		m.Name,
		m.Email,
		m.FirebaseUID.GetOrZero(),
		m.CreatedAt,
	)
}

func (r *Repository) FindOrCreateByFirebaseUID(ctx context.Context, firebaseUID, email, name string) (*userentity.User, error) {
	user, err := r.FindByFirebaseUID(ctx, firebaseUID)
	if err == nil {
		return user, nil
	}
	if !errs.IsNotFound(err) {
		return nil, err
	}

	// 初回ログイン：新規ユーザーを作成
	now := time.Now()
	newID, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("failed to generate user id: %w", err)
	}

	newUser, err := userentity.NewUserWithFirebaseUID(newID.String(), name, email, firebaseUID, now)
	if err != nil {
		return nil, fmt.Errorf("failed to create user entity: %w", err)
	}

	if err := r.Save(ctx, newUser); err != nil {
		return nil, fmt.Errorf("failed to save new user: %w", err)
	}

	return newUser, nil
}
