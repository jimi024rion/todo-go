package todo

// Handler is a collection of todo handlers.
type Handler struct {
	ListHandler   *ListHandler
	CreateHandler *CreateHandler
	GetHandler    *GetHandler
	UpdateHandler *UpdateHandler
	DeleteHandler *DeleteHandler
}

// NewHandler creates a new todo handler collection.
func NewHandler(
	listHandler *ListHandler,
	createHandler *CreateHandler,
	getHandler *GetHandler,
	updateHandler *UpdateHandler,
	deleteHandler *DeleteHandler,
) *Handler {
	return &Handler{
		ListHandler:   listHandler,
		CreateHandler: createHandler,
		GetHandler:    getHandler,
		UpdateHandler: updateHandler,
		DeleteHandler: deleteHandler,
	}
}
