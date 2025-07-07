package repositories

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"server-registry/internal/database/models"
)

type ServerRepository struct {
	db *sqlx.DB
}

func NewServerRepository(db *sqlx.DB) *ServerRepository {
	return &ServerRepository{db: db}
}

func (r *ServerRepository) CreateServer(token string, bridgeURL *string) (*models.Server, error) {
	query := `
		INSERT INTO servers (token, bridge_url) 
		VALUES ($1, $2)`

	var id int
	err := r.db.QueryRow(query+" RETURNING id", token, bridgeURL).Scan(&id)
	if err != nil {
		return nil, err
	}

	return r.GetServerByID(id)
}

func (r *ServerRepository) GetServerByToken(token string) (*models.Server, error) {
	query := `SELECT id, token, bridge_url, created_at, updated_at FROM servers WHERE token = $1`

	server := &models.Server{}
	err := r.db.Get(server, query, token)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return server, nil
}

func (r *ServerRepository) GetServerByID(serverID int) (*models.Server, error) {
	query := `SELECT id, token, bridge_url, created_at, updated_at FROM servers WHERE id = $1`

	server := &models.Server{}
	err := r.db.Get(server, query, serverID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return server, nil
}

func (r *ServerRepository) UpdateServerURL(token string, bridgeURL string) (*models.Server, error) {
	query := `
		UPDATE servers 
		SET bridge_url = $2, updated_at = CURRENT_TIMESTAMP 
		WHERE token = $1`

	result, err := r.db.Exec(query, token, bridgeURL)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, sql.ErrNoRows
	}

	return r.GetServerByToken(token)
}

func (r *ServerRepository) DeleteServer(token string) error {
	query := `DELETE FROM servers WHERE token = $1`

	result, err := r.db.Exec(query, token)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
