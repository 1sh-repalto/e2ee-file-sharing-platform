package repository

import (
	"context"
	"fmt"

	"github.com/1sh-repalto/e2ee-file-sharing-platform/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FileRepository interface {
	Save(ctx context.Context, file domain.File) error
	FindByID(ctx context.Context, id uuid.UUID) (domain.File, error)
	FindByOwner(ctx context.Context, ownerID uuid.UUID) ([]domain.File, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type fileRepository struct {
	db *pgxpool.Pool
}

func NewFileRepository(db *pgxpool.Pool) FileRepository {
	return &fileRepository{db: db}
}

func (r *fileRepository) Save(ctx context.Context, file domain.File) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO files (id, owner_id, filename, mime_type, size, iv, encrypted_key, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, file.ID, file.OwnerID, file.Filename, file.MimeType, file.Size, file.IV, file.EncryptedKey, file.CreatedAt)

	return err
}

func (r *fileRepository) FindByID(ctx context.Context, id uuid.UUID) (domain.File, error) {
	var f domain.File
	err := r.db.QueryRow(ctx, `
		SELECT id, owner_id, filename, mime_type, size, iv, encrypted_key, created_at
		FROM files WHERE id = $1
	`, id).Scan(&f.ID, &f.OwnerID, &f.Filename, &f.MimeType, &f.Size, &f.IV, &f.EncryptedKey, &f.CreatedAt)

	return f, err
}

func (r *fileRepository) FindByOwner(ctx context.Context, ownerID uuid.UUID) ([]domain.File, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, owner_id, filename, mime_type, size, iv, encrypted_key, created_at
		FROM files WHERE owner_id = $1
		ORDER BY created_at DESC
	`, ownerID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var files []domain.File
	for rows.Next() {
		var f domain.File
		if err := rows.Scan(&f.ID, &f.OwnerID, &f.Filename, &f.MimeType, &f.Size, &f.IV, &f.EncryptedKey, &f.CreatedAt); err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return files, nil
}

func (r *fileRepository) Delete(ctx context.Context, id uuid.UUID) error {
	cmd, err := r.db.Exec(ctx, `
		DELETE FROM files WHERE id = $1
	`, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("file not found")
	}
	return nil
}
