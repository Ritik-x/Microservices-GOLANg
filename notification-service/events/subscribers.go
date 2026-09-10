package events

import (
	"context"
	"encoding/json"
	"fmt"
	"notification-service/repositries"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

type UserCreatedEvent struct {
	EventID string `json:"event_id"`
	UserID  string `json:"user_id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
}

func ConnectNats() (*nats.Conn, error) {
	natsURL := os.Getenv("NATS_URL")

	nc, err := nats.Connect(natsURL)

	if err != nil {
		return nil, err
	}

	fmt.Println("Notification Service connected to NATS")

	return nc, nil
}

func StartConsumer(nc *nats.Conn, repo *repositries.NotificationRepository) error {
	js, err := nc.JetStream()
	if err != nil {
		return err
	}
	_, err = js.Subscribe("user.created", func(msg *nats.Msg) {

		var event UserCreatedEvent
		err := json.Unmarshal(msg.Data, &event)
		if err != nil {
			fmt.Println("Invalid event:", err)
			fmt.Println("RAW EVENT:", string(msg.Data))
			return
		}

		fmt.Println("Received events")
		fmt.Println(string(msg.Data))
		message := fmt.Sprintf(
			"Welcome %s! Your account has been created.",
			event.Name,
		)
		err = repo.CreateNotification(
			context.Background(),
			event.EventID,
			event.UserID,
			message,
		)

		if err != nil {
			fmt.Println("Failed to save notification:", err)
			return
		}

		fmt.Println("Notification saved successfully")

		msg.Ack()
	},
		nats.Durable("notification-service"),
		nats.ManualAck(),
		nats.AckWait(30*time.Second),
		nats.MaxDeliver(5),
	)

	return err
}
