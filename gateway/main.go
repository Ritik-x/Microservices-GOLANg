// @title Microservices Assignment API
// @version 1.0
// @description API Gateway for User and Notification microservices.
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

package main

import (
	"fmt"
	_ "gateway/docs"
	"net/http"

	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
)

// HealthCheck godoc
// @Summary Check API Gateway health
// @Description Returns the health status of the API Gateway
// @Tags Health
// @Produce plain
// @Success 200 {string} string "API Gateway is running"
// @Router /health [get]
func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "APi gateway is running")
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return a JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body object{email=string,password=string} true "Login credentials"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Router /api/auth/login [post]
func loginHandler(w http.ResponseWriter, r *http.Request) {
	authServiceProxy().ServeHTTP(w, r)
}

// CreateUser godoc
// @Summary Create a new user
// @Description Creates a new user account
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body object{name=string,email=string,password=string} true "User details"
// @Success 201 {object} map[string]string
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Router /api/users [post]
func createUserHandler(w http.ResponseWriter, r *http.Request) {
	userServiceProxy().ServeHTTP(w, r)
}

func main() {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env file found, using environment variables")
	}

	http.HandleFunc("/health", healthHandler)

	http.Handle("/swagger/", httpSwagger.WrapHandler)

	fmt.Println("API Gateway running on port 8080")

	http.HandleFunc(
		"/api/auth/login",
		loginHandler,
	)

	http.Handle(
		"/api/users",
		AuthMiddleware(http.HandlerFunc(createUserHandler)),
	)

	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("Gateway error:", err)
	}
}
