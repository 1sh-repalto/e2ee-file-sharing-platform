package usecase

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/1sh-repalto/e2ee-file-sharing-platform/internal/domain"
	"github.com/1sh-repalto/e2ee-file-sharing-platform/internal/repository"
	"github.com/1sh-repalto/e2ee-file-sharing-platform/internal/storage"
	"github.com/google/uuid"
)

type FileUsecase interface {
	Upload(ctx context.Context, file domain.File, content io.ReadCloser) error
	Download(ctx context.Context, id, recipientID uuid.UUID) (io.ReadCloser, domain.File, []byte, error)
	GetByID(ctx context.Context, id uuid.UUID) (domain.File, error)
	ListByOwner(ctx context.Context, ownerID uuid.UUID) ([]domain.File, error)
	Delete(ctx context.Context, id uuid.UUID, ownerID uuid.UUID) error
}

type fileUsecase struct {
	fileRepo repository.FileRepository
	shareRepo repository.ShareRepository
	storage storage.Storage
}

func NewFileUsecase(fileRepo repository.FileRepository,shareRepo repository.ShareRepository, storage storage.Storage) FileUsecase {
	return &fileUsecase{fileRepo: fileRepo, shareRepo: shareRepo, storage: storage}
}

func (u *fileUsecase) Upload(ctx context.Context, file domain.File, content io.ReadCloser) error {
	if file.ID == uuid.Nil {
		file.ID = uuid.New()
	}
	if file.CreatedAt.IsZero() {
		file.CreatedAt = time.Now().UTC()
	}

	if err := u.storage.Upload(ctx, "files", file.ID.String(), content, file.Size, file.MimeType); err != nil {
		return err
	}

	return u.fileRepo.Save(ctx, file)
}

func (u *fileUsecase) Download(ctx context.Context, id uuid.UUID, recipientID uuid.UUID) (io.ReadCloser, domain.File, []byte, error) {
	file, err := u.fileRepo.FindByID(ctx, id)
	if err != nil {
		return nil, domain.File{}, nil, err
	}

	if file.OwnerID == recipientID {
		content, err := u.storage.Download(ctx, "files", id.String())
		if err != nil {
			return nil, domain.File{}, nil, err
		}
		return content, file, nil, nil
	}

	share, err := u.shareRepo.FindByFileAndRecipient(ctx, id, recipientID)
	if err != nil || share.ID == uuid.Nil {
		return nil, domain.File{}, nil, errors.New("unauthorized: you don't have access to this file")
	}

	content, err := u.storage.Download(ctx, "files", id.String())
	if err != nil {
		return nil, domain.File{}, nil, err
	}

	return content, file, share.WrappedKey, nil
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

	if err := u.storage.Delete(ctx, "files", id.String()); err != nil {
		return err
	}

	return u.fileRepo.Delete(ctx, id)
}