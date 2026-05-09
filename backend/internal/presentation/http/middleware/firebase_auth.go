package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/jimi024rion/todo-go/backend/internal/config/env"
	"github.com/jimi024rion/todo-go/backend/internal/config/errs"
	userrepo "github.com/jimi024rion/todo-go/backend/internal/domain/user/repository"
	firebaseinfra "github.com/jimi024rion/todo-go/backend/internal/infrastructure/firebase"
	"github.com/jimi024rion/todo-go/backend/internal/presentation/http/response"
)

func FirebaseAuth(verifier firebaseinfra.TokenVerifier, userRepo userrepo.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// non-production 環境のみ X-Skip-Auth によるテスト用認証スキップを許可
		if env.Cfg.AppEnv != "production" {
			if skipUserID := c.GetHeader("X-Skip-Auth"); skipUserID != "" {
				c.Set("user_id", skipUserID)
				c.Next()
				return
			}
		}

		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(
				errs.IsUnauthorizedRequest.HTTPStatus(),
				response.Fail(errs.IsUnauthorizedRequest.Message()),
			)
			c.Abort()
			return
		}

		idToken := strings.TrimPrefix(authHeader, "Bearer ")
		ctx := c.Request.Context()

		token, err := verifier.VerifyIDToken(ctx, idToken)
		if err != nil {
			c.JSON(
				errs.IsUnauthorizedRequest.HTTPStatus(),
				response.Fail(errs.IsUnauthorizedRequest.Message()),
			)
			c.Abort()
			return
		}

		// Firebase Claims から email と name を取得（Google 認証で付与される）
		email, _ := token.Claims["email"].(string)
		name, _ := token.Claims["name"].(string)
		if name == "" {
			name = email
		}

		user, err := userRepo.FindOrCreateByFirebaseUID(ctx, token.UID, email, name)
		if err != nil {
			c.JSON(
				errs.IsUnauthorizedRequest.HTTPStatus(),
				response.Fail(errs.IsUnauthorizedRequest.Message()),
			)
			c.Abort()
			return
		}

		c.Set("user_id", user.ID().String())
		c.Next()
	}
}
