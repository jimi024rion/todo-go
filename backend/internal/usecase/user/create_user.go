package user

type CreateUserUsecase struct{}

func NewCreateUserUsecase() *CreateUserUsecase {
	return &CreateUserUsecase{}
}

func (u *CreateUserUsecase) Execute() {
}
