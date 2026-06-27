package file

import (
	"context"
	"fmt"
	"time"

	filerepo "github.com/jimi024rion/todo-go/backend/internal/domain/file/repository"
)

const defaultExpiry = 15 * time.Minute

type GetDownloadURLUseCase struct {
	storageRepo filerepo.StorageRepository
}

func NewGetDownloadURLUseCase(storageRepo filerepo.StorageRepository) *GetDownloadURLUseCase {
	return &GetDownloadURLUseCase{storageRepo: storageRepo}
}

type GetDownloadURLInput struct {
	Path string
}

type GetDownloadURLOutput struct {
	URL       string
	ExpiresAt time.Time
}

func (uc *GetDownloadURLUseCase) Execute(ctx context.Context, input *GetDownloadURLInput) (*GetDownloadURLOutput, error) {
	if input.Path == "" {
		return nil, fmt.Errorf("path is required")
	}

	downloadURL, err := uc.storageRepo.GetDownloadURL(ctx, input.Path, defaultExpiry)
	if err != nil {
		return nil, fmt.Errorf("get download url: %w", err)
	}

	return &GetDownloadURLOutput{
		URL:       downloadURL,
		ExpiresAt: time.Now().Add(defaultExpiry),
	}, nil
}
