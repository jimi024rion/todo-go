package valueobject

import (
	"fmt"

	"github.com/google/uuid"
)

type UserID uuid.UUID

func NewUserID() UserID {
	return UserID(uuid.New())
}

func UserIDFromString(id string) (UserID, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return UserID{}, fmt.Errorf("invalid user id format: %w", err)
	}
	return UserID(uid), nil
}

func (id UserID) String() string {
	return uuid.UUID(id).String()
}

func (id UserID) Value() string {
	return id.String()
}
