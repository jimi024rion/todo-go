package entity

import "time"

// Tag はタグエンティティです。
type Tag struct {
	id        string
	userID    string
	name      string
	createdAt time.Time
	updatedAt time.Time
}

// NewTag は新しいTagを生成します。
func NewTag(id, userID, name string, now time.Time) *Tag {
	return &Tag{
		id:        id,
		userID:    userID,
		name:      name,
		createdAt: now,
		updatedAt: now,
	}
}

// Reconstruct はDBの値からTagを復元します。
func Reconstruct(id, userID, name string, createdAt, updatedAt time.Time) *Tag {
	return &Tag{
		id:        id,
		userID:    userID,
		name:      name,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (t *Tag) ID() string        { return t.id }
func (t *Tag) UserID() string    { return t.userID }
func (t *Tag) Name() string      { return t.name }
func (t *Tag) CreatedAt() time.Time { return t.createdAt }
func (t *Tag) UpdatedAt() time.Time { return t.updatedAt }
