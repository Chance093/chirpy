// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package database

import (
	"time"

	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID
	Body      string
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
}

type User struct {
	ID        uuid.UUID
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
	Password  string
}
