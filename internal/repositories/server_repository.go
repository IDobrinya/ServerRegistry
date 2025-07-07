package repositories

import (
	"database/sql"
	"errors"
	"server-registry/internal/models"
)

type ServerRepository struct {
	db *sql.DB
}

func NewServerRepository(db *sql.DB) *ServerRepository {
	return &ServerRepository{db: db}
}

func (r *ServerRepository) CreateServer(token string, bridgeURL *string) (*models.Server, error) {
	query := `
		INSERT INTO servers (token, bridge_url) 
		VALUES ($1, $2)`

	result, err := r.db.Exec(query, token, bridgeURL)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return r.GetServerByID(int(id))
}

func (r *ServerRepository) GetServerByToken(token string) (*models.Server, error) {
	query := `SELECT id, token, bridge_url, created_at, updated_at FROM servers WHERE token = $1`

	server := &models.Server{}
	err := r.db.QueryRow(query, token).Scan(
		&server.ID,
		&server.Token,
		&server.BridgeURL,
		&server.CreatedAt,
		&server.UpdatedAt,
	)

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
	err := r.db.QueryRow(query, serverID).Scan(
		&server.ID,
		&server.Token,
		&server.BridgeURL,
		&server.CreatedAt,
		&server.UpdatedAt,
	)

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
