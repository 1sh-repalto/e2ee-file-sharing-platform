package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                  uuid.UUID `db:"id"`
	Username            string    `db:"username"`
	PasswordHash        string    `db:"password_hash"`
	PublicKey           string    `db:"public_key"`
	EncryptedPrivateKey []byte    `db:"encrypted_private_key"`
	CreatedAt           time.Time `db:"created_at"`
}
