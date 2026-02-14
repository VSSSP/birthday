package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/vsssp/birthday-app/backend/internal/adapter/repository/postgres"
	"github.com/vsssp/birthday-app/backend/internal/domain"
	"github.com/vsssp/birthday-app/backend/internal/pkg/hash"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	ctx := context.Background()
	pool, err := postgres.NewPool(ctx, databaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	userRepo := postgres.NewUserRepository(pool)
	providerRepo := postgres.NewAuthProviderRepository(pool)
	recipientRepo := postgres.NewRecipientRepository(pool)

	// Create demo user
	passwordHash, _ := hash.HashPassword("password123")
	now := time.Now()
	demoUser := &domain.User{
		ID:           uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		Email:        "demo@example.com",
		Name:         "Demo User",
		PasswordHash: &passwordHash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := userRepo.Create(ctx, demoUser); err != nil {
		log.Printf("user may already exist: %v", err)
	} else {
		log.Println("created demo user: demo@example.com / password123")
	}

	// Create email provider link
	emailLink := &domain.AuthProviderLink{
		ID:          uuid.New(),
		UserID:      demoUser.ID,
		Provider:    domain.AuthProviderEmail,
		ProviderUID: demoUser.Email,
		CreatedAt:   now,
	}
	if err := providerRepo.Create(ctx, emailLink); err != nil {
		log.Printf("provider link may already exist: %v", err)
	}

	// Create sample recipients
	recipients := []domain.Recipient{
		{
			ID:        uuid.New(),
			UserID:    demoUser.ID,
			Name:      "Maria Silva",
			Age:       65,
			Gender:    "female",
			MinBudget: 50,
			MaxBudget: 200,
			Keywords:  []string{"cooking", "reading", "gardening"},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        uuid.New(),
			UserID:    demoUser.ID,
			Name:      "Pedro Santos",
			Age:       12,
			Gender:    "male",
			MinBudget: 30,
			MaxBudget: 100,
			Keywords:  []string{"gaming", "nerd", "tech"},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        uuid.New(),
			UserID:    demoUser.ID,
			Name:      "Ana Costa",
			Age:       30,
			Gender:    "female",
			MinBudget: 100,
			MaxBudget: 500,
			Keywords:  []string{"fashion", "travel", "fitness"},
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	for _, r := range recipients {
		if err := recipientRepo.Create(ctx, &r); err != nil {
			log.Printf("failed to create recipient %s: %v", r.Name, err)
		} else {
			log.Printf("created recipient: %s", r.Name)
		}
	}

	log.Println("seed completed successfully")
}
