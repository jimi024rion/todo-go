package handler

import (
	"github.com/jimi024rion/todo-go/backend/internal/presentation/http/health"
	"github.com/jimi024rion/todo-go/backend/internal/presentation/http/todo"
)

type Handler struct {
	Health *health.Handler
	Todo   *todo.Handler
}

func NewHandler(
	healthHandler *health.Handler,
	todoHandler *todo.Handler,
) *Handler {
	return &Handler{
		Health: healthHandler,
		Todo:   todoHandler,
	}
}
