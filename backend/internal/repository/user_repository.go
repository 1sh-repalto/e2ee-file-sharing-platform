package repository

import (
	"context"

	"github.com/1sh-repalto/e2ee-file-sharing-platform/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Save(ctx context.Context, user domain.User) error
	FindByUsername(ctx context.Context, username string) (domain.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (domain.User, error)
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Save(ctx context.Context, user domain.User) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO users (id, username, password_hash, public_key, encrypted_private_key, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, user.ID, user.Username, user.PasswordHash, user.PublicKey, user.EncryptedPrivateKey, user.CreatedAt)

	return err
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (domain.User, error) {
	var u domain.User
	err := r.db.QueryRow(ctx, `
		SELECT id, username, password_hash, public_key, encrypted_private_key, created_at 
		FROM users WHERE username=$1
	`, username).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.PublicKey, &u.EncryptedPrivateKey, &u.CreatedAt)

	return u, err
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	var u domain.User
	err := r.db.QueryRow(ctx, `
		SELECT id, username, password_hash, public_key, encrypted_private_key, created_at 
		FROM users WHERE id=$1
	`, id).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.PublicKey, &u.EncryptedPrivateKey, &u.CreatedAt)

	return u, err
}
