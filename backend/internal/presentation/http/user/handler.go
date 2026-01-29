package user

type Handler struct {
	Create *CreateUserHandler
}

func NewHandler(
	createUserHandler *CreateUserHandler,
) *Handler {
	return &Handler{
		Create: createUserHandler,
	}
}
