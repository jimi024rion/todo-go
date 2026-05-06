package entity

import (
	"fmt"
	"time"

	"github.com/jimi024rion/todo-go/backend/internal/domain/user/model/valueobject"
)

type User struct {
	id          valueobject.UserID
	name        string
	email       string
	firebaseUID string
	createdAt   time.Time
	updatedAt   time.Time
}

func NewUser(id, name, email string, now time.Time) (*User, error) {
	userID, err := valueobject.UserIDFromString(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}

	return &User{
		id:        userID,
		name:      name,
		email:     email,
		createdAt: now,
		updatedAt: now,
	}, nil
}

func NewUserWithFirebaseUID(id, name, email, firebaseUID string, now time.Time) (*User, error) {
	userID, err := valueobject.UserIDFromString(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}

	return &User{
		id:          userID,
		name:        name,
		email:       email,
		firebaseUID: firebaseUID,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

// Getter
func (u *User) ID() valueobject.UserID { return u.id }
func (u *User) Name() string           { return u.name }
func (u *User) Email() string          { return u.email }
func (u *User) FirebaseUID() string    { return u.firebaseUID }
func (u *User) CreatedAt() time.Time   { return u.createdAt }
func (u *User) UpdatedAt() time.Time   { return u.updatedAt }

func (u *User) UpdateName(name string, now time.Time) {
	u.name = name
	u.updatedAt = now
}

func (u *User) UpdateEmail(email string, now time.Time) {
	u.email = email
	u.updatedAt = now
}
