package user

import (
	"context"
	"time"

	"github.com/jimi024rion/todo-go/backend/internal/config/clock"
	"github.com/jimi024rion/todo-go/backend/internal/domain/tx"
	userentity "github.com/jimi024rion/todo-go/backend/internal/domain/user/model/entity"
	uservoid "github.com/jimi024rion/todo-go/backend/internal/domain/user/model/valueobject"
	userrepo "github.com/jimi024rion/todo-go/backend/internal/domain/user/repository"
)

type CreateUserUsecase struct {
	userRepo  userrepo.UserRepository
	txManager tx.TxManager
	clock     clock.Clock
}

func NewCreateUserUsecase(userRepo userrepo.UserRepository, txManager tx.TxManager, clock clock.Clock) *CreateUserUsecase {
	return &CreateUserUsecase{
		userRepo:  userRepo,
		txManager: txManager,
		clock:     clock,
	}
}

type CreateUserInput struct {
	Name  string
	Email string
}

type CreateUserOutput struct {
	ID        string
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *CreateUserUsecase) Execute(ctx context.Context, input *CreateUserInput) (*CreateUserOutput, error) {
	result, err := u.txManager.Do(ctx, func(txCtx context.Context) (any, error) {
		now := u.clock.Now(txCtx)
		newID := uservoid.NewUserID()

		user, err := userentity.NewUser(newID.String(), input.Name, input.Email, now)
		if err != nil {
			return nil, err
		}

		if err := u.userRepo.Save(txCtx, user); err != nil {
			return nil, err
		}

		return &CreateUserOutput{
			ID:        user.ID().String(),
			Name:      user.Name(),
			Email:     user.Email(),
			CreatedAt: user.CreatedAt(),
			UpdatedAt: user.UpdatedAt(),
		}, nil
	})
	if err != nil {
		return nil, err
	}

	return result.(*CreateUserOutput), nil
}
