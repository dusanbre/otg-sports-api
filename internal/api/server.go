package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	_ "github.com/dusanbre/otg-sports-api/internal/api/docs" // Swagger docs
	"github.com/dusanbre/otg-sports-api/internal/api/handlers"
	"github.com/dusanbre/otg-sports-api/internal/api/middleware"
	"github.com/dusanbre/otg-sports-api/internal/database"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// Server represents the API server
type Server struct {
	db     *database.DB
	port   string
	server *http.Server
}

// NewServer creates a new API server
func NewServer(db *database.DB, port string) *Server {
	return &Server{
		db:   db,
		port: port,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	router := s.setupRouter()

	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%s", s.port),
		Handler: router,
	}

	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	if s.server != nil {
		return s.server.Shutdown(ctx)
	}
	return nil
}

// setupRouter configures the Chi router with all routes and middleware
func (s *Server) setupRouter() *chi.Mux {
	r := chi.NewRouter()

	// Global middleware
	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	// CORS middleware
	corsConfig := middleware.CORSConfig{
		AllowedOrigins: s.getAllowedOrigins(),
		AllowedMethods: []string{"GET", "OPTIONS"},
		AllowedHeaders: []string{"Authorization", "X-API-Key", "Content-Type"},
	}
	r.Use(middleware.CORS(corsConfig))

	// Create handlers
	healthHandler := handlers.NewHealthHandler()
	soccerHandler := handlers.NewSoccerHandler(s.db)
	basketballHandler := handlers.NewBasketballHandler(s.db)

	// Create rate limiter
	rateLimiter := middleware.NewRateLimiter(s.getDefaultRateLimit())

	// Public routes
	r.Get("/health", healthHandler.Health)

	// Swagger documentation (available at /swagger/index.html)
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	// Protected API routes
	r.Route("/api/v1", func(r chi.Router) {
		// API key authentication
		r.Use(middleware.APIKeyAuth(s.db))

		// Rate limiting (per API key)
		r.Use(rateLimiter.Middleware)

		// Soccer routes
		r.Route("/soccer", func(r chi.Router) {
			r.Use(middleware.RequireSport("soccer"))
			r.Get("/matches", soccerHandler.GetMatches)
			r.Get("/matches/{id}", soccerHandler.GetMatch)
			r.Get("/matches/live", soccerHandler.GetLiveMatches)
			r.Get("/leagues", soccerHandler.GetLeagues)
		})

		// Basketball routes
		r.Route("/basketball", func(r chi.Router) {
			r.Use(middleware.RequireSport("basketball"))
			r.Get("/matches", basketballHandler.GetMatches)
			r.Get("/matches/{id}", basketballHandler.GetMatch)
			r.Get("/matches/live", basketballHandler.GetLiveMatches)
			r.Get("/leagues", basketballHandler.GetLeagues)
		})
	})

	return r
}

// getAllowedOrigins returns the list of allowed CORS origins
func (s *Server) getAllowedOrigins() []string {
	origins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if origins == "" {
		return []string{"*"}
	}
	return strings.Split(origins, ",")
}

// getDefaultRateLimit returns the default rate limit from env
func (s *Server) getDefaultRateLimit() int {
	// Default 100 requests per minute
	return 100
}
