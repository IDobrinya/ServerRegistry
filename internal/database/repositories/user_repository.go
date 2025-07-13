package repositories

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"server-registry/internal/database/models"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(userID string, linkedServerID *string) (*models.User, error) {
	query := `
		INSERT INTO users (id, linked_server) 
		VALUES ($1, $2)`

	_, err := r.db.Exec(query, userID, linkedServerID)
	if err != nil {
		return nil, err
	}

	return r.GetUserByID(userID)
}

func (r *UserRepository) GetUserByID(userID string) (*models.User, error) {
	query := `SELECT id, linked_server, created_at, updated_at FROM users WHERE id = $1`

	user := &models.User{}
	err := r.db.Get(user, query, userID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetUserByLinkedServer(linkedServerID string) (*models.User, error) {
	query := `SELECT id, linked_server, created_at, updated_at FROM users WHERE linked_server = $1`

	user := &models.User{}
	err := r.db.Get(user, query, linkedServerID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) UpdateUserLinkedServer(userID string, linkedServerID int) (*models.User, error) {
	query := `
		UPDATE users 
		SET linked_server = $2, updated_at = CURRENT_TIMESTAMP 
		WHERE id = $1`

	_, err := r.db.Exec(query, userID, linkedServerID)
	if err != nil {
		return nil, err
	}

	return r.GetUserByID(userID)
}

func (r *UserRepository) UnlinkAllUsersFromServer(serverToken string) (int64, error) {
	query := `
		UPDATE users 
		SET linked_server = NULL, updated_at = CURRENT_TIMESTAMP 
		WHERE linked_server = $1`

	result, err := r.db.Exec(query, serverToken)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}
