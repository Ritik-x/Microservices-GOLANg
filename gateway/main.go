package main

import (
	"fmt"
	"net/http"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "APi gateway is running")
}

func main() {
	http.HandleFunc("/health", healthHandler)
	fmt.Println("API Gateway running on port 8080")
	http.Handle("/api/users", userServiceProxy())
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Gateway error:", err)
	}
}
