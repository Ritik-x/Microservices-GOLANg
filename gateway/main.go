package main

import (
	"fmt"
	"net/http"

	"github.com/joho/godotenv"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "APi gateway is running")
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env:", err)
		return
	}
	http.HandleFunc("/health", healthHandler)
	fmt.Println("API Gateway running on port 8080")
	http.Handle(
		"/api/auth/login",
		authServiceProxy(),
	)
	http.Handle("/api/users", AuthMiddleware(userServiceProxy()))
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Gateway error:", err)
	}
}
