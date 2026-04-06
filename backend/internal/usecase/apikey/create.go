package apikey

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/jimi024rion/todo-go/backend/internal/config/clock"
	apikeyentity "github.com/jimi024rion/todo-go/backend/internal/domain/apikey/model/entity"
	apikeyvo "github.com/jimi024rion/todo-go/backend/internal/domain/apikey/model/valueobject"
	apikeyrepo "github.com/jimi024rion/todo-go/backend/internal/domain/apikey/repository"
	"github.com/jimi024rion/todo-go/backend/internal/domain/tx"
)

const keyPrefix = "todo_"

type CreateUseCase struct {
	apiKeyRepo apikeyrepo.APIKeyRepository
	txManager  tx.TxManager
	clock      clock.Clock
}

func NewCreateUseCase(apiKeyRepo apikeyrepo.APIKeyRepository, txManager tx.TxManager, clock clock.Clock) *CreateUseCase {
	return &CreateUseCase{
		apiKeyRepo: apiKeyRepo,
		txManager:  txManager,
		clock:      clock,
	}
}

type CreateInput struct {
	UserID string
	Name   string
}

type CreateOutput struct {
	ID        string
	Key       string // 平文（一度限り返す）
	UserID    string
	Name      string
	CreatedAt time.Time
}

func (uc *CreateUseCase) Execute(ctx context.Context, input *CreateInput) (*CreateOutput, error) {
	fullKey, err := generateKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate api key: %w", err)
	}

	keyHash := apikeyentity.HashKey(fullKey)

	result, err := uc.txManager.Do(ctx, func(txCtx context.Context) (any, error) {
		now := uc.clock.Now(txCtx)
		newID := apikeyvo.NewAPIKeyID()

		apiKey, err := apikeyentity.NewAPIKey(newID.String(), input.UserID, keyHash, input.Name, now)
		if err != nil {
			return nil, err
		}

		if err := uc.apiKeyRepo.Save(txCtx, apiKey); err != nil {
			return nil, err
		}

		return &CreateOutput{
			ID:        apiKey.ID().String(),
			Key:       fullKey,
			UserID:    apiKey.UserID(),
			Name:      apiKey.Name(),
			CreatedAt: apiKey.CreatedAt(),
		}, nil
	})
	if err != nil {
		return nil, err
	}

	return result.(*CreateOutput), nil
}

// generateKey は todo_ プレフィックス付きの69文字のAPIキーを生成します。
func generateKey() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return keyPrefix + hex.EncodeToString(b), nil
}
