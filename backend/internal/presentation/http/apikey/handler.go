package apikey

type Handler struct {
	Create *CreateHandler
	Delete *DeleteHandler
}

func NewHandler(create *CreateHandler, delete *DeleteHandler) *Handler {
	return &Handler{
		Create: create,
		Delete: delete,
	}
}
