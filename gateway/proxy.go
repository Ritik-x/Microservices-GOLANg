package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func userServiceProxy() http.Handler {
	targetURL := os.Getenv("USER_SERVICE_URL")
	target, err := url.Parse(targetURL)
	if err != nil {
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	originalDirector := proxy.Director

	proxy.Director = func(req *http.Request) {

		originalDirector(req)

		req.URL.Path = "/users"
	}

	return proxy
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

func authServiceProxy() http.Handler {
	targetURL := os.Getenv("USER_SERVICE_URL")

	target, err := url.Parse(targetURL)
	if err != nil {
		panic(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = "/login"
	}
	return proxy
}
