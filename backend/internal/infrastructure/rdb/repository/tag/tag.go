package tag

import (
	"context"
	"fmt"

	"github.com/aarondl/opt/omit"
	"github.com/gofrs/uuid/v5"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/dm"
	"github.com/stephenafamo/bob/dialect/psql/sm"

	tagentity "github.com/jimi024rion/todo-go/backend/internal/domain/tag/model/entity"
	tagrepo "github.com/jimi024rion/todo-go/backend/internal/domain/tag/repository"
	"github.com/jimi024rion/todo-go/backend/internal/infrastructure/rdb"
	"github.com/jimi024rion/todo-go/backend/internal/infrastructure/rdb/models"
)

// Ensure Repository implements TagRepository.
var _ tagrepo.TagRepository = (*Repository)(nil)

type Repository struct {
	db bob.DB
}

func NewRepository(db bob.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) exec(ctx context.Context) bob.Executor {
	if e, ok := rdb.GetExecutor(ctx); ok {
		return e
	}
	return r.db
}

// Create はタグを新規作成します。
func (r *Repository) Create(ctx context.Context, tag *tagentity.Tag) error {
	id, err := uuid.FromString(tag.ID())
	if err != nil {
		return fmt.Errorf("invalid tag id: %w", err)
	}
	userID, err := uuid.FromString(tag.UserID())
	if err != nil {
		return fmt.Errorf("invalid user id: %w", err)
	}

	setter := &models.TagSetter{
		ID:        omit.From(id),
		UserID:    omit.From(userID),
		Name:      omit.From(tag.Name()),
		CreatedAt: omit.From(tag.CreatedAt()),
		UpdatedAt: omit.From(tag.UpdatedAt()),
	}
	_, err = models.Tags.Insert(setter).Exec(ctx, r.exec(ctx))
	if err != nil {
		return fmt.Errorf("failed to insert tag: %w", err)
	}
	return nil
}

// FindByUserID はユーザーに紐づくタグ一覧を返します。
func (r *Repository) FindByUserID(ctx context.Context, userID string) ([]*tagentity.Tag, error) {
	uid, err := uuid.FromString(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}

	rows, err := models.Tags.Query(
		sm.Where(models.Tags.Columns.UserID.EQ(psql.Arg(uid))),
	).All(ctx, r.exec(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to find tags by user_id: %w", err)
	}

	tags := make([]*tagentity.Tag, len(rows))
	for i, row := range rows {
		tags[i] = toEntity(row)
	}
	return tags, nil
}

// FindByTodoIDs は複数のTodoIDに紐づくタグをまとめて返します。
func (r *Repository) FindByTodoIDs(ctx context.Context, todoIDs []string) (map[string][]*tagentity.Tag, error) {
	result := make(map[string][]*tagentity.Tag)
	if len(todoIDs) == 0 {
		return result, nil
	}

	exec := r.exec(ctx)

	for _, todoID := range todoIDs {
		todoUID, err := uuid.FromString(todoID)
		if err != nil {
			return nil, fmt.Errorf("invalid todo id: %w", err)
		}

		// 1. todo_tags からこの todoID の tag_id 一覧を取得
		todoTags, err := models.TodoTags.Query(
			sm.Where(models.TodoTags.Columns.TodoID.EQ(psql.Arg(todoUID))),
		).All(ctx, exec)
		if err != nil {
			return nil, fmt.Errorf("failed to find todo_tags: %w", err)
		}

		if len(todoTags) == 0 {
			continue
		}

		// 2. 各 tag_id でタグを取得
		for _, tt := range todoTags {
			tag, err := models.FindTag(ctx, exec, tt.TagID)
			if err != nil {
				continue // タグが削除されていた場合はスキップ
			}
			result[todoID] = append(result[todoID], toEntity(tag))
		}
	}

	return result, nil
}

// Delete はタグを削除します（所有者チェック含む）。
func (r *Repository) Delete(ctx context.Context, id, userID string) error {
	uid, err := uuid.FromString(id)
	if err != nil {
		return fmt.Errorf("invalid tag id: %w", err)
	}
	uUserID, err := uuid.FromString(userID)
	if err != nil {
		return fmt.Errorf("invalid user id: %w", err)
	}

	_, err = models.Tags.Delete(
		dm.Where(models.Tags.Columns.ID.EQ(psql.Arg(uid))),
		dm.Where(models.Tags.Columns.UserID.EQ(psql.Arg(uUserID))),
	).Exec(ctx, r.exec(ctx))
	if err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}
	return nil
}

// SetTodoTags はTodoに紐づくタグを一括更新します（全削除→再挿入）。
func (r *Repository) SetTodoTags(ctx context.Context, todoID string, tagIDs []string) error {
	todoUID, err := uuid.FromString(todoID)
	if err != nil {
		return fmt.Errorf("invalid todo id: %w", err)
	}

	exec := r.exec(ctx)

	// 既存の todo_tags を全削除
	_, err = models.TodoTags.Delete(
		dm.Where(models.TodoTags.Columns.TodoID.EQ(psql.Arg(todoUID))),
	).Exec(ctx, exec)
	if err != nil {
		return fmt.Errorf("failed to delete todo_tags: %w", err)
	}

	// 新しい tag_ids を挿入
	for _, tagID := range tagIDs {
		tagUID, err := uuid.FromString(tagID)
		if err != nil {
			return fmt.Errorf("invalid tag id: %w", err)
		}
		setter := &models.TodoTagSetter{
			TodoID: omit.From(todoUID),
			TagID:  omit.From(tagUID),
		}
		if _, err = models.TodoTags.Insert(setter).Exec(ctx, exec); err != nil {
			return fmt.Errorf("failed to insert todo_tag: %w", err)
		}
	}
	return nil
}

func toEntity(m *models.Tag) *tagentity.Tag {
	return tagentity.Reconstruct(
		m.ID.String(),
		m.UserID.String(),
		m.Name,
		m.CreatedAt,
		m.UpdatedAt,
	)
}
