package middleware

import (
	"github.com/gin-gonic/gin"

	apikeyCache "github.com/jimi024rion/todo-go/backend/internal/domain/apikey/cache"
	apikeyrepo "github.com/jimi024rion/todo-go/backend/internal/domain/apikey/repository"
)

type Middleware struct {
	APIKeyAuth gin.HandlerFunc
}

func NewMiddleware(apiKeyRepo apikeyrepo.APIKeyRepository, apiKeyCache apikeyCache.APIKeyCache) *Middleware {
	return &Middleware{
		APIKeyAuth: APIKeyAuth(apiKeyRepo, apiKeyCache),
	}
}
