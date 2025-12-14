package http

import (
	"github.com/jimi024rion/todo-go/backend/internal/presentation/http/health"
)

type Handler struct {
	Health *health.Handler
}

func NewHandler(
	healthHandler *health.Handler,
) *Handler {
	return &Handler{
		Health: healthHandler,
	}
}
