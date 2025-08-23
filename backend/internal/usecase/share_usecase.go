package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/1sh-repalto/e2ee-file-sharing-platform/internal/domain"
	"github.com/1sh-repalto/e2ee-file-sharing-platform/internal/repository"
	"github.com/google/uuid"
)

type ShareUsecase interface {
	ShareFile(ctx context.Context, fileID, ownerID, recipientID uuid.UUID, wrappedKey []byte) error
	GetSharesForRecipient(ctx context.Context, recipientID uuid.UUID) ([]domain.Share, error)
	GetShare(ctx context.Context, fileID, recipientID uuid.UUID) (domain.Share, error)
	Unshare(ctx context.Context, shareID, ownerID uuid.UUID) error
}

type shareUsecase struct {
	shareRepo repository.ShareRepository
	fileRepo repository.FileRepository
}

func NewShareUsecase(shareRepo repository.ShareRepository, fileRepo repository.FileRepository) ShareUsecase {
	return &shareUsecase{shareRepo: shareRepo, fileRepo: fileRepo}
}

func (u *shareUsecase) ShareFile(ctx context.Context, fileID, ownerID, recipientID uuid.UUID, wrappedKey []byte) error {
	file, err := u.fileRepo.FindByID(ctx, fileID)
	if err != nil {
		return err
	}

	if(file.OwnerID != ownerID) {
		return errors.New("unauthorized action")
	}

	share := domain.Share{
		ID:				uuid.New(),
		FileID:			fileID,
		RecipientID: 	recipientID,
		WrappedKey: 	wrappedKey,
		CreatedAt:		time.Now().UTC(),
	}

	return u.shareRepo.Save(ctx, share)
}

func (u *shareUsecase) GetSharesForRecipient(ctx context.Context, recipientID uuid.UUID) ([]domain.Share, error) {
	return u.shareRepo.FindByRecipient(ctx, recipientID)
}

func (u *shareUsecase) GetShare(ctx context.Context, fileID, recipientID uuid.UUID) (domain.Share, error) {
	return u.shareRepo.FindByFileAndRecipient(ctx, fileID, recipientID)
}

func (u *shareUsecase) Unshare(ctx context.Context, shareID, ownerID uuid.UUID) error {
	share, err := u.shareRepo.FindByID(ctx, shareID)
	if err != nil {
		return err
	}

	file, err := u.fileRepo.FindByID(ctx, share.FileID)
	if err != nil {
		return err
	}

	if file.OwnerID != ownerID {
		return errors.New("unauthorized action")
	}

	return u.shareRepo.Delete(ctx, shareID)
}