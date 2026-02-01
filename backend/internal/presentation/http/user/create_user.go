package user

import (
	useruc "github.com/jimi024rion/todo-go/backend/internal/usecase/user"
)

type CreateUserHandler struct {
	u *useruc.CreateUserUsecase
}

func NewCreateUserHandler(createUserUsecase *useruc.CreateUserUsecase) *CreateUserHandler {
	return &CreateUserHandler{
		u: createUserUsecase,
	}
}

func (h *CreateUserHandler) Handle() {
	h.u.Execute()
}
