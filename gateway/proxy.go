package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func userServiceProxy() http.Handler {

	target, err := url.Parse("http://localhost:8081")
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

func authServiceProxy() http.Handler {
	target, err := url.Parse("http://localhost:8081")
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
