package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vsssp/birthday-app/backend/internal/adapter/handler"
	"github.com/vsssp/birthday-app/backend/internal/adapter/repository/postgres"
	"github.com/vsssp/birthday-app/backend/internal/adapter/social"
	"github.com/vsssp/birthday-app/backend/internal/config"
	jwtpkg "github.com/vsssp/birthday-app/backend/internal/pkg/jwt"
	"github.com/vsssp/birthday-app/backend/internal/usecase"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Database connection
	pool, err := postgres.NewPool(context.Background(), cfg.Database.URL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Repositories
	userRepo := postgres.NewUserRepository(pool)
	providerRepo := postgres.NewAuthProviderRepository(pool)
	tokenRepo := postgres.NewRefreshTokenRepository(pool)
	recipientRepo := postgres.NewRecipientRepository(pool)

	// Services
	jwtService := jwtpkg.NewService(
		cfg.JWT.AccessSecret,
		cfg.JWT.RefreshSecret,
		cfg.JWT.AccessExpiry,
		cfg.JWT.RefreshExpiry,
	)

	googleVerifier := social.NewGoogleVerifier(cfg.Google.ClientID)
	appleVerifier := social.NewAppleVerifier(cfg.Apple.ClientID)
	socialVerifier := social.NewCompositeVerifier(googleVerifier, appleVerifier)

	// Use cases
	authUseCase := usecase.NewAuthUseCase(userRepo, providerRepo, tokenRepo, jwtService, socialVerifier)
	userUseCase := usecase.NewUserUseCase(userRepo)
	recipientUseCase := usecase.NewRecipientUseCase(recipientRepo)

	// Router
	router := handler.NewRouter(authUseCase, userUseCase, recipientUseCase, jwtService)

	// Server
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		log.Printf("server starting on port %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}
	log.Println("server exited")
}
