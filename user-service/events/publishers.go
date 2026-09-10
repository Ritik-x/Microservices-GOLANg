package events

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/uuid"

	"github.com/nats-io/nats.go"
)

func ConnectNats() (*nats.Conn, error) {

	natsURL := os.Getenv("NATS_URL")
	nc, err := nats.Connect(natsURL)

	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to NATS")
	return nc, nil
}

func PublishUserCreated(
	js nats.JetStreamContext,
	userID string,
	name string,
	email string,
) error {

	eventID := uuid.New()

	event := struct {
		EventID string `json:"event_id"`
		UserID  string `json:"user_id"`
		Name    string `json:"name"`
		Email   string `json:"email"`
	}{
		EventID: eventID.String(),
		UserID:  userID,
		Name:    name,
		Email:   email,
	}

	message, err := json.Marshal(event)
	if err != nil {
		return err
	}

	fmt.Println("Publishing event:", string(message))

	_, err = js.Publish(
		"user.created",
		message,
	)

	return err
}
