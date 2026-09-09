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
func (repo *UserRepository) CreateUser(ctx context.Context, name string, email string) (string, error) {
	var userID string

	err := repo.DB.QueryRow(
		ctx,
		`INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id`,
		name,
		email,
	).Scan(&userID)

	return userID, err

}
