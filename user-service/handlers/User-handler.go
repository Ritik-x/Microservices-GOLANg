package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"user-service/events"
	model "user-service/models"
	"user-service/repositries"

	"github.com/nats-io/nats.go"
)

type UserHandler struct {
	Repo *repositries.UserRepository
	JS   nats.JetStreamContext
}

func NewUserHandler(
	repo *repositries.UserRepository,
	js nats.JetStreamContext,
) *UserHandler {
	return &UserHandler{
		Repo: repo,
		JS:   js,
	}
}

func (h *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {

	var user model.User

	// Request body se JSON read karna
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validation
	if user.Name == "" || user.Email == "" {
		http.Error(w, "Name and email are required", http.StatusBadRequest)
		return
	}

	// Database mein user save karna
	userID, err := h.Repo.CreateUser(
		r.Context(),
		user.Name,
		user.Email,
	)

	if err != nil {
		fmt.Println("CreateUser error:", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}
	err = events.PublishUserCreated(h.JS, userID, user.Name, user.Email)
	if err != nil {
		http.Error(w, "Failed to publish user event", http.StatusInternalServerError)
		return
	}

	// Response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(user)
}
