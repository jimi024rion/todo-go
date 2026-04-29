package tag

import (
	"context"

	tagrepo "github.com/jimi024rion/todo-go/backend/internal/domain/tag/repository"
	"github.com/jimi024rion/todo-go/backend/internal/domain/tx"
)

type SetTodoTagsUseCase struct {
	tagRepo   tagrepo.TagRepository
	txManager tx.TxManager
}

func NewSetTodoTagsUseCase(tagRepo tagrepo.TagRepository, txManager tx.TxManager) *SetTodoTagsUseCase {
	return &SetTodoTagsUseCase{tagRepo: tagRepo, txManager: txManager}
}

type SetTodoTagsInput struct {
	TodoID string
	TagIDs []string
}

func (uc *SetTodoTagsUseCase) Execute(ctx context.Context, input *SetTodoTagsInput) error {
	_, err := uc.txManager.Do(ctx, func(txCtx context.Context) (any, error) {
		return nil, uc.tagRepo.SetTodoTags(txCtx, input.TodoID, input.TagIDs)
	})
	return err
}
