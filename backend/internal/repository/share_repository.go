package repository

import (
	"context"
	"fmt"

	"github.com/1sh-repalto/e2ee-file-sharing-platform/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ShareRepository interface {
	Save(ctx context.Context, share domain.Share) error
	FindByID(ctx context.Context,  shareID uuid.UUID) (domain.Share, error)
	FindByRecipient(ctx context.Context, recipientID uuid.UUID) ([]domain.Share, error)
	FindByFileAndRecipient(ctx context.Context, fileID, recipientID uuid.UUID) (domain.Share, error)
	Delete(ctx context.Context, shareID uuid.UUID) error
}

type shareRepository struct {
	db *pgxpool.Pool
}

func NewShareRepository(db *pgxpool.Pool) ShareRepository {
	return &shareRepository{db: db}
}

func (r *shareRepository) Save(ctx context.Context, share domain.Share) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO shares (id, file_id, recipient_id, wrapped_key, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`, share.ID, share.FileID, share.RecipientID, share.WrappedKey, share.CreatedAt)

	return err
}

func (r *shareRepository) FindByID(ctx context.Context, shareID uuid.UUID) (domain.Share, error) {
	var s domain.Share
	err := r.db.QueryRow(ctx, `
		SELECT id, file_id, recipient_id, wrapped_key, created_at
		FROM shares
		WHERE id = $1
	`, shareID).Scan(&s.ID, &s.FileID, &s.RecipientID, &s.WrappedKey, &s.CreatedAt)

	return s, err
}

func (r *shareRepository) FindByRecipient(ctx context.Context, recipientID uuid.UUID) ([]domain.Share, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, file_id, recipient_id, wrapped_key, created_at
		FROM shares WHERE recipient_id = $1
		ORDER BY created_at DESC
	`, recipientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shares []domain.Share
	for rows.Next() {
		var s domain.Share
		if err := rows.Scan(&s.ID, &s.FileID, &s.RecipientID, &s.WrappedKey, &s.CreatedAt); err != nil {
			return nil, err
		}

		shares = append(shares, s)
	}

	return shares, rows.Err()
}

func (r *shareRepository) FindByFileAndRecipient(ctx context.Context, fileID, recipientID uuid.UUID) (domain.Share, error) {
	var s domain.Share
	err := r.db.QueryRow(ctx, `
		SELECT id, file_id, recipient_id, wrapped_key, created_at
		FROM shares WHERE file_id = $1 AND recipient_id = $2	
	`, fileID, recipientID).Scan(&s.ID, &s.FileID, &s.RecipientID, &s.WrappedKey, &s.CreatedAt)

	return s, err
}

func (r *shareRepository) Delete(ctx context.Context, shareID uuid.UUID) error {
	cmd, err := r.db.Exec(ctx, `
		DELETE FROM shares WHERE id = $1
	`, shareID)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("share not found")
	}
	return nil
}