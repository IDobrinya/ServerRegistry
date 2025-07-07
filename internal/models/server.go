package models

import (
	"database/sql"
	"time"
)

type Server struct {
	ID        int            `json:"id" db:"id"`
	Token     string         `json:"token" db:"token"`
	BridgeURL sql.NullString `json:"bridge_url" db:"bridge_url"`
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" db:"updated_at"`
}
