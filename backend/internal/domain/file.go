package domain

import (
	"time"

	"github.com/google/uuid"
)

type File struct {
	ID        		uuid.UUID	`db:"id"`
	OwnerID   		uuid.UUID	`db:"owner_id"`
	Filename 		string      `db:"name"`
	MimeType  		string	    `db:"mime_type"`
	Size      		int64       `db:"size"`
	IV        		[]byte      `db:"iv"`
	Encrypted_key	[]byte	    `db:"encrypted_key"`
	CreatedAt 		time.Time   `db:"created_at"`
}
