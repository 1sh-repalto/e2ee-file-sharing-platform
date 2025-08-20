package domain

import (
	"time"

	"github.com/google/uuid"
)

type Share struct {
	ID          uuid.UUID	`db:"id"`
	FileID      uuid.UUID	`db:"file_id"`
	RecipientID uuid.UUID   `db:"recipient_id"`
	WrappedKey  []byte      `db:"wrapped_key"`
	CreatedAt   time.Time   `db:"created_at"`
}
