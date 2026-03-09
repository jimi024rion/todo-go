package entity

import (
	"context"
	"time"

	"github.com/jimi024rion/todo-go/backend/internal/config/logger"
	"github.com/jimi024rion/todo-go/backend/internal/domain/user/model/valueobject"
)

type User struct {
	id        valueobject.UserID
	name      string
	email     string
	createdAt time.Time
	updatedAt time.Time
}

func NewUser(ctx context.Context, id, name, email string, now time.Time) *User {
	l := logger.NewLogger(ctx)
	userID, err := valueobject.UserIDFromString(id)
	if err != nil {
		l.ErrorLog(err)
	}

	return &User{
		id:        userID,
		name:      name,
		email:     email,
		createdAt: now,
		updatedAt: now,
	}
}

// Getter
func (u *User) ID() valueobject.UserID { return u.id }
func (u *User) Name() string           { return u.name }
func (u *User) Email() string          { return u.email }
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
