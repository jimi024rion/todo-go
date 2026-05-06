package repository

import (
	"context"

	"github.com/jimi024rion/todo-go/backend/internal/domain/user/model/entity"
)

type UserRepository interface {
	Save(ctx context.Context, user *entity.User) error
	FindByFirebaseUID(ctx context.Context, firebaseUID string) (*entity.User, error)
	FindOrCreateByFirebaseUID(ctx context.Context, firebaseUID, email, name string) (*entity.User, error)
}
