package repositries

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type NotificationRepository struct {
	DB *pgxpool.Pool
}

func NewNotificationRepository(db *pgxpool.Pool) *NotificationRepository {
	return &NotificationRepository{
		DB: db,
	}
}

func (repo *NotificationRepository) CreateNotification(
	ctx context.Context,
	eventID string,
	userID string,
	message string,
) error {
	_, err := repo.DB.Exec(
		ctx, `INSERT INTO notifications (event_id, user_id, message)
		 VALUES ($1, $2, $3)`,
		eventID,
		userID,
		message,
	)
	return err
}
