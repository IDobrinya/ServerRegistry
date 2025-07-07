package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID           string        `json:"id" db:"id"`
	LinkedServer sql.NullInt32 `json:"linked_server" db:"linked_server"`
	CreatedAt    time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at" db:"updated_at"`
}
