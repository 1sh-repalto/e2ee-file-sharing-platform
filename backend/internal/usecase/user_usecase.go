package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/1sh-repalto/e2ee-file-sharing-platform/internal/domain"
	"github.com/1sh-repalto/e2ee-file-sharing-platform/internal/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
)

type UserUsecase interface {
	Register(ctx context.Context, username, password string, publicKey, encryptedPrivateKey []byte) (domain.User, error)
	Login(ctx context.Context, username, password string) (domain.User, []byte, error)
	GetByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	GetByUsername(ctx context.Context, username string) (domain.User, error)
}

type userUsecase struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &userUsecase{userRepo: userRepo}
}

func hashPassword(password string) []byte {
	salt := []byte("static-salt")
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	return hash
}

func (uc *userUsecase) Register(ctx context.Context, username, password string, publicKey, encryptedPrivateKey []byte) (domain.User, error) {
	_, err := uc.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return domain.User{}, errors.New("username already exists")
	}

	user := domain.User{
		ID:                  uuid.New(),
		Username:            username,
		PasswordHash:        string(hashPassword(password)),
		PublicKey:           string(publicKey),
		EncryptedPrivateKey: encryptedPrivateKey,
		CreatedAt:           time.Now(),
	}

	err = uc.userRepo.Save(ctx, user)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (uc *userUsecase) Login(ctx context.Context, username, password string) (domain.User, []byte, error) {
	user, err := uc.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return domain.User{}, nil, errors.New("invalid credentials")
	}

	hash := string(hashPassword(password))
	if hash != user.PasswordHash {
		return domain.User{}, nil, errors.New("invalid credentials")
	}

	return user, user.EncryptedPrivateKey, nil
}

func (uc *userUsecase) GetByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	return uc.userRepo.FindByID(ctx, id)
}

func (uc *userUsecase) GetByUsername(ctx context.Context, username string) (domain.User, error) {
	return uc.userRepo.FindByUsername(ctx, username)
}
