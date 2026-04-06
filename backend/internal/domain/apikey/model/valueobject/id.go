package valueobject

import (
	"fmt"

	"github.com/google/uuid"
)

type APIKeyID uuid.UUID

func NewAPIKeyID() APIKeyID {
	return APIKeyID(uuid.New())
}

func APIKeyIDFromString(id string) (APIKeyID, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return APIKeyID{}, fmt.Errorf("failed to parse api key id: %w", err)
	}
	return APIKeyID(uid), nil
}

func (id APIKeyID) String() string {
	return uuid.UUID(id).String()
}
