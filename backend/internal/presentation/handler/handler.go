package handler

import (
	"github.com/jimi024rion/todo-go/backend/internal/presentation/http/health"
	"github.com/jimi024rion/todo-go/backend/internal/presentation/http/todo"
)

type Handler struct {
	HealthHandler *health.Handler
	TodoHandler   *todo.Handler
}

func NewHandler(
	healthHandler *health.Handler,
	todoHandler *todo.Handler,
) *Handler {
	return &Handler{
		HealthHandler: healthHandler,
		TodoHandler:   todoHandler,
	}
}
