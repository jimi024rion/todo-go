package tag

import (
	"context"
	"fmt"
	"time"

	"github.com/gofrs/uuid/v5"

	"github.com/jimi024rion/todo-go/backend/internal/config/clock"
	"github.com/jimi024rion/todo-go/backend/internal/domain/tag/model/entity"
	tagrepo "github.com/jimi024rion/todo-go/backend/internal/domain/tag/repository"
)

type CreateUseCase struct {
	tagRepo tagrepo.TagRepository
	clock   clock.Clock
}

func NewCreateUseCase(tagRepo tagrepo.TagRepository, clock clock.Clock) *CreateUseCase {
	return &CreateUseCase{tagRepo: tagRepo, clock: clock}
}

type CreateInput struct {
	UserID string
	Name   string
}

type CreateOutput struct {
	ID        string
	UserID    string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (uc *CreateUseCase) Execute(ctx context.Context, input *CreateInput) (*CreateOutput, error) {
	if input.Name == "" {
		return nil, fmt.Errorf("tag name is required")
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("failed to generate tag id: %w", err)
	}

	now := uc.clock.Now(ctx)
	tag := entity.NewTag(id.String(), input.UserID, input.Name, now)

	if err := uc.tagRepo.Create(ctx, tag); err != nil {
		return nil, err
	}

	return &CreateOutput{
		ID:        tag.ID(),
		UserID:    tag.UserID(),
		Name:      tag.Name(),
		CreatedAt: tag.CreatedAt(),
		UpdatedAt: tag.UpdatedAt(),
	}, nil
}
