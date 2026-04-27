package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/jimi024rion/todo-go/backend/internal/config/errs"
	apikeyCache "github.com/jimi024rion/todo-go/backend/internal/domain/apikey/cache"
	apikeyentity "github.com/jimi024rion/todo-go/backend/internal/domain/apikey/model/entity"
	apikeyrepo "github.com/jimi024rion/todo-go/backend/internal/domain/apikey/repository"
	"github.com/jimi024rion/todo-go/backend/internal/presentation/http/response"
)

func APIKeyAuth(repo apikeyrepo.APIKeyRepository, cache apikeyCache.APIKeyCache) gin.HandlerFunc {
	return func(c *gin.Context) {
		rawKey := c.GetHeader("X-API-Key")
		if rawKey == "" {
			c.JSON(
				errs.IsUnauthorizedRequest.HTTPStatus(),
				response.Fail(errs.IsUnauthorizedRequest.Message()),
			)
			c.Abort()
			return
		}

		ctx := c.Request.Context()

		if apiKey, ok := cache.Get(ctx, rawKey); ok {
			c.Set("user_id", apiKey.UserID())
			c.Next()
			return
		}

		keyHash := apikeyentity.HashKey(rawKey)
		apiKey, err := repo.FindByKeyHash(ctx, keyHash)
		if err != nil {
			c.JSON(
				errs.IsUnauthorizedRequest.HTTPStatus(),
				response.Fail(errs.IsUnauthorizedRequest.Message()),
			)
			c.Abort()
			return
		}

		cache.Set(ctx, rawKey, apiKey)

		c.Set("user_id", apiKey.UserID())
		c.Next()
	}
}
