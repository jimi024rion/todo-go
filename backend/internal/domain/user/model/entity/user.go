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

func NewUser(ctx context.Context, id, name, email string) *User {
	l := logger.NewLogger(ctx)
	now := time.Now()
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

func (u *User) UpdateName(name string) {
	u.name = name
	u.updatedAt = time.Now()
}

func (u *User) UpdateEmail(email string) {
	u.email = email
	u.updatedAt = time.Now()
}
