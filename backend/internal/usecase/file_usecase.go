package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/1sh-repalto/e2ee-file-sharing-platform/internal/domain"
	"github.com/1sh-repalto/e2ee-file-sharing-platform/internal/repository"
	"github.com/google/uuid"
)

type FileUsecase interface {
	Upload(ctx context.Context, file domain.File) error
	GetByID(ctx context.Context, id uuid.UUID) (domain.File, error)
	ListByOwner(ctx context.Context, ownerID uuid.UUID) ([]domain.File, error)
	Delete(ctx context.Context, id uuid.UUID, ownerID uuid.UUID) error 
}

type fileUsecase struct {
	fileRepo repository.FileRepository
}

func NewFileUsecase(fileRepo repository.FileRepository) FileUsecase {
	return &fileUsecase{fileRepo: fileRepo}
}

func (u *fileUsecase) Upload(ctx context.Context, file domain.File) error {
	if file.ID == uuid.Nil {
		file.ID = uuid.New()
	}
	if file.CreatedAt.IsZero() {
		file.CreatedAt = time.Now().UTC()
	}

	return u.fileRepo.Save(ctx, file)
}

func (u *fileUsecase) GetByID(ctx context.Context, id uuid.UUID) (domain.File, error) {
	return u.fileRepo.FindByID(ctx, id)
}

func (u *fileUsecase) ListByOwner(ctx context.Context, ownerID uuid.UUID) ([]domain.File, error) {
	return u.fileRepo.FindByOwner(ctx, ownerID)
}

func (u *fileUsecase) Delete(ctx context.Context, id uuid.UUID, ownerID uuid.UUID) error {
	file, err := u.fileRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if file.OwnerID != ownerID {
		return errors.New("unauthorized: cannot delete someone else's file")
	}

	return u.fileRepo.Delete(ctx, id)
}