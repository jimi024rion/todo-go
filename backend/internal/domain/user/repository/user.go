package repository

import (
	"context"

	"github.com/jimi024rion/todo-go/backend/internal/domain/user/model/entity"
)

type UserRepository interface {
	Save(ctx context.Context, user *entity.User) error
}
