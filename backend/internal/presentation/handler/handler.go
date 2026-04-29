package handler

import (
	"github.com/jimi024rion/todo-go/backend/internal/presentation/http/apikey"
	"github.com/jimi024rion/todo-go/backend/internal/presentation/http/health"
	"github.com/jimi024rion/todo-go/backend/internal/presentation/http/tag"
	"github.com/jimi024rion/todo-go/backend/internal/presentation/http/todo"
	"github.com/jimi024rion/todo-go/backend/internal/presentation/http/user"
)

type Handler struct {
	Health *health.Handler
	User   *user.Handler
	Todo   *todo.Handler
	APIKey *apikey.Handler
	Tag    *tag.Handler
}

func NewHandler(
	healthHandler *health.Handler,
	userHandler *user.Handler,
	todoHandler *todo.Handler,
	apiKeyHandler *apikey.Handler,
	tagHandler *tag.Handler,
) *Handler {
	return &Handler{
		Health: healthHandler,
		User:   userHandler,
		Todo:   todoHandler,
		APIKey: apiKeyHandler,
		Tag:    tagHandler,
	}
}
