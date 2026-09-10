package repositries

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	DB *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (repo *UserRepository) CreateUser(
	ctx context.Context,
	name string,
	email string,
	passwordHash string,
) (string, error) {

	var userID string

	err := repo.DB.QueryRow(
		ctx,
		`INSERT INTO users (name, email, password_hash)
		 VALUES ($1, $2, $3)
		 RETURNING id`,
		name,
		email,
		passwordHash,
	).Scan(&userID)

	return userID, err
}

func (repo *UserRepository) GetUserByEmail(
	ctx context.Context,
	email string,
) (string, string, string, error) {

	var userID string
	var name string
	var passwordHash string

	err := repo.DB.QueryRow(
		ctx,
		`SELECT id, name, password_hash
		 FROM users
		 WHERE email = $1`,
		email,
	).Scan(&userID, &name, &passwordHash)

	return userID, name, passwordHash, err
}
