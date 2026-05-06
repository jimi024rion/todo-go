package firebase

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"

	"github.com/jimi024rion/todo-go/backend/internal/config/env"
)

// TokenVerifier は Firebase ID Token を検証するインターフェース。
// テスト時にモックに差し替えられるよう interface で定義する。
type TokenVerifier interface {
	VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error)
}

type firebaseVerifier struct {
	client *auth.Client
}

// NewTokenVerifier は Firebase Admin SDK を初期化して TokenVerifier を返す。
// GOOGLE_APPLICATION_CREDENTIALS が設定されていればそのキーファイルを使用し、
// Cloud Run では ADC（Application Default Credentials）を自動使用する。
func NewTokenVerifier(ctx context.Context) (TokenVerifier, error) {
	app, err := firebase.NewApp(ctx, &firebase.Config{
		ProjectID: env.Cfg.FirebaseProjectID,
	})
	if err != nil {
		return nil, fmt.Errorf("firebase app init: %w", err)
	}

	client, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("firebase auth client init: %w", err)
	}

	return &firebaseVerifier{client: client}, nil
}

func (v *firebaseVerifier) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	token, err := v.client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, fmt.Errorf("verify id token: %w", err)
	}
	return token, nil
}
