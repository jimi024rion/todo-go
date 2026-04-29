package tag

import (
	"context"

	tagrepo "github.com/jimi024rion/todo-go/backend/internal/domain/tag/repository"
)

type DeleteUseCase struct {
	tagRepo tagrepo.TagRepository
}

func NewDeleteUseCase(tagRepo tagrepo.TagRepository) *DeleteUseCase {
	return &DeleteUseCase{tagRepo: tagRepo}
}

type DeleteInput struct {
	ID     string
	UserID string
}

func (uc *DeleteUseCase) Execute(ctx context.Context, input *DeleteInput) error {
	return uc.tagRepo.Delete(ctx, input.ID, input.UserID)
}
