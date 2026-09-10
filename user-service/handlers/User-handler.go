package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"user-service/auth"
	"user-service/events"
	model "user-service/models"
	"user-service/repositries"

	"github.com/nats-io/nats.go"
	"golang.org/x/crypto/bcrypt"
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
	if user.Name == "" || user.Email == "" || user.Password == "" {
		http.Error(w, "Name  email , password are reqired are required", http.StatusBadRequest)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		http.Error(w, "Failed to process password", http.StatusInternalServerError)
		return
	}
	// Database mein user save karna
	userID, err := h.Repo.CreateUser(
		r.Context(),
		user.Name,
		user.Email,
		string(passwordHash),
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
	response := map[string]string{
		"id":    userID,
		"name":  user.Name,
		"email": user.Email,
	}

	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	userID, name, passwordHash, err := h.Repo.GetUserByEmail(
		r.Context(),
		req.Email,
	)

	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(passwordHash),
		[]byte(req.Password),
	)

	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}
	token, err := auth.GenerateJWT(userID)

	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"id":    userID,
		"name":  name,
		"token": token,
	})
}
