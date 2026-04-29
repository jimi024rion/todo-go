package tag

import (
	"context"
	"time"

	tagrepo "github.com/jimi024rion/todo-go/backend/internal/domain/tag/repository"
)

type ListUseCase struct {
	tagRepo tagrepo.TagRepository
}

func NewListUseCase(tagRepo tagrepo.TagRepository) *ListUseCase {
	return &ListUseCase{tagRepo: tagRepo}
}

type ListInput struct {
	UserID string
}

type ListOutput struct {
	Tags []TagItem
}

type TagItem struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (uc *ListUseCase) Execute(ctx context.Context, input *ListInput) (*ListOutput, error) {
	tags, err := uc.tagRepo.FindByUserID(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	items := make([]TagItem, len(tags))
	for i, t := range tags {
		items[i] = TagItem{
			ID:        t.ID(),
			Name:      t.Name(),
			CreatedAt: t.CreatedAt(),
			UpdatedAt: t.UpdatedAt(),
		}
	}
	return &ListOutput{Tags: items}, nil
}
