package entity

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/jimi024rion/todo-go/backend/internal/domain/apikey/model/valueobject"
)

type APIKey struct {
	id        valueobject.APIKeyID
	userID    string
	keyHash   string
	name      string
	createdAt time.Time
}

func NewAPIKey(id, userID, keyHash, name string, now time.Time) (*APIKey, error) {
	apiKeyID, err := valueobject.APIKeyIDFromString(id)
	if err != nil {
		return nil, fmt.Errorf("invalid api key id: %w", err)
	}
	return &APIKey{
		id:        apiKeyID,
		userID:    userID,
		keyHash:   keyHash,
		name:      name,
		createdAt: now,
	}, nil
}

func (a *APIKey) ID() valueobject.APIKeyID { return a.id }
func (a *APIKey) UserID() string           { return a.userID }
func (a *APIKey) KeyHash() string          { return a.keyHash }
func (a *APIKey) Name() string             { return a.name }
func (a *APIKey) CreatedAt() time.Time     { return a.createdAt }

// HashKey はrawKeyのSHA-256ハッシュを返します。
func HashKey(rawKey string) string {
	h := sha256.Sum256([]byte(rawKey))
	return hex.EncodeToString(h[:])
}

// VerifyKey はrawKeyがこのAPIキーに対応するか検証します。
func (a *APIKey) VerifyKey(rawKey string) bool {
	return HashKey(rawKey) == a.keyHash
}
