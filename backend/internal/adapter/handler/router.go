package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/vsssp/birthday-app/backend/internal/port"
	jwtpkg "github.com/vsssp/birthday-app/backend/internal/pkg/jwt"
)

// NewRouter sets up all HTTP routes and middleware.
func NewRouter(
	authService port.AuthService,
	userService port.UserService,
	recipientService port.RecipientService,
	jwtService *jwtpkg.Service,
) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	authHandler := NewAuthHandler(authService)
	userHandler := NewUserHandler(userService)
	recipientHandler := NewRecipientHandler(recipientService)
	authMiddleware := NewAuthMiddleware(jwtService)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	r.Route("/api", func(r chi.Router) {
		// Public auth routes
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authHandler.Register)
			r.Post("/login", authHandler.Login)
			r.Post("/google", authHandler.GoogleLogin)
			r.Post("/apple", authHandler.AppleLogin)
			r.Post("/refresh", authHandler.RefreshToken)
		})

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.Authenticate)

			r.Get("/auth/me", userHandler.GetCurrentUser)

			r.Route("/recipients", func(r chi.Router) {
				r.Post("/", recipientHandler.Create)
				r.Get("/", recipientHandler.List)
				r.Delete("/", recipientHandler.BulkDelete)
				r.Get("/{id}", recipientHandler.GetByID)
				r.Put("/{id}", recipientHandler.Update)
				r.Delete("/{id}", recipientHandler.Delete)
			})
		})
	})

	return r
}
