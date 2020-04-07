package service

import "net/http"

// Defines a single route, e.g. a human readable name, HTTP method, pattern the function that will execute when the route is called.
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Defines the type Routes which is just an array (slice) of Route structs.
type Routes []Route

// Initialize our routes
var routes = Routes{

	Route{
		"GetQuote", // Name
		"GET",      // HTTP method
		"/quotes/", // Route pattern
		GetQuote,
	},
	Route{
		"HealthCheck",
		"GET",
		"/healthz",
		HealthCheck,
	},
	Route{
		"Testability",
		"GET",
		"/testability/healthy/{state}",
		SetHealthyState,
	},
}
