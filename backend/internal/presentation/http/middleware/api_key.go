package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	apikeyCache "github.com/jimi024rion/todo-go/backend/internal/domain/apikey/cache"
	apikeyentity "github.com/jimi024rion/todo-go/backend/internal/domain/apikey/model/entity"
	apikeyrepo "github.com/jimi024rion/todo-go/backend/internal/domain/apikey/repository"
)

func APIKeyAuth(repo apikeyrepo.APIKeyRepository, cache apikeyCache.APIKeyCache) gin.HandlerFunc {
	return func(c *gin.Context) {
		rawKey := c.GetHeader("X-API-Key")
		if rawKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing API key"})
			c.Abort()
			return
		}

		ctx := c.Request.Context()

		// キャッシュを確認
		if apiKey, ok := cache.Get(ctx, rawKey); ok {
			c.Set("user_id", apiKey.UserID())
			c.Next()
			return
		}

		// SHA-256ハッシュで直接DB検索
		keyHash := apikeyentity.HashKey(rawKey)
		apiKey, err := repo.FindByKeyHash(ctx, keyHash)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid API key"})
			c.Abort()
			return
		}

		// キャッシュに保存
		cache.Set(ctx, rawKey, apiKey)

		c.Set("user_id", apiKey.UserID())
		c.Next()
	}
}
