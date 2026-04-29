package tag

// Handler はタグ関連のハンドラーをまとめた構造体です。
type Handler struct {
	List        *ListHandler
	Create      *CreateHandler
	Delete      *DeleteHandler
	SetTodoTags *SetTodoTagsHandler
}

func NewHandler(
	list *ListHandler,
	create *CreateHandler,
	delete *DeleteHandler,
	setTodoTags *SetTodoTagsHandler,
) *Handler {
	return &Handler{
		List:        list,
		Create:      create,
		Delete:      delete,
		SetTodoTags: setTodoTags,
	}
}
