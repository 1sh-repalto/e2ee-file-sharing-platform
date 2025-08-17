package domain

import (
	"time"

	"github.com/google/uuid"
)

type File struct {
	ID			uuid.UUID	`db:"id"`
	OwnerID		uuid.UUID	`db:"owner_id"`
	Name		string		`db:"name"`
	Size		int64		`db:"size"`
	IV			[]byte		`db:"iv"`
	CreatedAt	time.Time	`db:"created_at"`
}