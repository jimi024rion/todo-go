package middleware

import (
	"github.com/gin-gonic/gin"

	userrepo "github.com/jimi024rion/todo-go/backend/internal/domain/user/repository"
	firebaseinfra "github.com/jimi024rion/todo-go/backend/internal/infrastructure/firebase"
)

type Middleware struct {
	FirebaseAuth gin.HandlerFunc
}

func NewMiddleware(verifier firebaseinfra.TokenVerifier, userRepo userrepo.UserRepository) *Middleware {
	return &Middleware{
		FirebaseAuth: FirebaseAuth(verifier, userRepo),
	}
}
