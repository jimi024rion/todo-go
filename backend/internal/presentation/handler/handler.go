package handler

import (
	"github.com/jimi024rion/todo-go/backend/internal/presentation/http/health"
)

type Handler struct {
	HealthHandler *health.Handler
}

func NewHandler(
	healthHandler *health.Handler,
) *Handler {
	return &Handler{
		HealthHandler: healthHandler,
	}
}
