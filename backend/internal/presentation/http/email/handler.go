package email

type Handler struct {
	SendWelcome *SendWelcomeHandler
}

func NewHandler(sendWelcome *SendWelcomeHandler) *Handler {
	return &Handler{
		SendWelcome: sendWelcome,
	}
}
