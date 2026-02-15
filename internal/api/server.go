package api

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/utkarsh/unfollow-tracker/internal/api/handlers"
	"github.com/utkarsh/unfollow-tracker/internal/api/middleware"
	"github.com/utkarsh/unfollow-tracker/internal/config"
	"github.com/utkarsh/unfollow-tracker/internal/repository"
	"github.com/utkarsh/unfollow-tracker/internal/service"
)

type Server struct {
	config      *config.Config
	db          *sql.DB
	authSvc     *service.AuthService
	accountSvc  *service.AccountService
	unfollowSvc *service.UnfollowService
	router      *chi.Mux
	httpSrv     *http.Server
}

func NewServer(cfg *config.Config, db *sql.DB) *Server {
	// Initialize repositories
	userRepo := repository.NewPostgresUserRepository(db)
	accountRepo := repository.NewPostgresAccountRepository(db)
	snapshotRepo := repository.NewPostgresSnapshotRepository(db)
	unfollowRepo := repository.NewPostgresUnfollowRepository(db)

	// Initialize services
	authSvc := service.NewAuthService(userRepo, cfg.JWT)
	accountSvc := service.NewAccountService(accountRepo, userRepo)
	unfollowSvc := service.NewUnfollowService(unfollowRepo, accountRepo)

	s := &Server{
		config:      cfg,
		db:          db,
		authSvc:     authSvc,
		accountSvc:  accountSvc,
		unfollowSvc: unfollowSvc,
		router:      chi.NewRouter(),
	}

	s.setupMiddleware()
	s.setupRoutes()

	s.httpSrv = &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      s.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Prevent unused variable warnings
	_ = snapshotRepo

	return s
}

func (s *Server) setupMiddleware() {
	// Built-in middleware
	s.router.Use(chimiddleware.RequestID)
	s.router.Use(chimiddleware.RealIP)
	s.router.Use(chimiddleware.Logger)
	s.router.Use(chimiddleware.Recoverer)
	s.router.Use(chimiddleware.Timeout(30 * time.Second))

	// Custom middleware
	s.router.Use(middleware.CORS())
}

func (s *Server) setupRoutes() {
	// Health check
	s.router.Get("/health", handlers.Health)

	// API v1 routes
	s.router.Route("/api/v1", func(r chi.Router) {
		// Public routes
		r.Group(func(r chi.Router) {
			r.Post("/auth/register", s.Register)
			r.Post("/auth/login", s.Login)
		})

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(s.config.JWT.Secret))

			// User
			r.Get("/me", s.GetCurrentUser)

			// Accounts
			r.Route("/accounts", func(r chi.Router) {
				r.Get("/", s.ListAccounts)
				r.Post("/", s.CreateAccount)
				r.Get("/{accountID}", s.GetAccount)
				r.Delete("/{accountID}", s.DeleteAccount)
				r.Get("/{accountID}/stats", s.GetAccountStats)
			})

			// Unfollows
			r.Route("/unfollows", func(r chi.Router) {
				r.Get("/", s.ListUnfollows)
				r.Get("/summary", s.GetUnfollowSummary)
			})
		})
	})
}

func (s *Server) Start() error {
	// Graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("🚀 Server starting on port %s", s.config.Server.Port)
		if err := s.httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-done
	log.Println("Server stopping...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.httpSrv.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	log.Println("Server stopped gracefully")
	return nil
}

func (s *Server) Router() *chi.Mux {
	return s.router
}
