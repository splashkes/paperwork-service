package middleware

import (
	"net/http"

	"github.com/rs/cors"
)

// CORSMiddleware creates a CORS middleware
func CORSMiddleware() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // In production, specify exact origins
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodHead,
			http.MethodOptions,
		},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-CSRF-Token",
		},
		ExposedHeaders: []string{
			"Content-Length",
			"Content-Type",
			"Content-Disposition",
		},
		AllowCredentials: false,
		MaxAge:           300, // 5 minutes
	})
}