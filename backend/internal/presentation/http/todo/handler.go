package todo

// Handler is a collection of todo handlers.
type Handler struct {
	List   *ListHandler
	Create *CreateHandler
	Get    *GetHandler
	Update *UpdateHandler
	Delete *DeleteHandler
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
		List:   listHandler,
		Create: createHandler,
		Get:    getHandler,
		Update: updateHandler,
		Delete: deleteHandler,
	}
}
