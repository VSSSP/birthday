package handler_test

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/vsssp/birthday-app/backend/internal/domain"
)

// mockUserRepo implements port.UserRepository in memory.
type mockUserRepo struct {
	mu    sync.RWMutex
	users map[uuid.UUID]*domain.User
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{users: make(map[uuid.UUID]*domain.User)}
}

func (r *mockUserRepo) Create(_ context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[user.ID] = user
	return nil
}

func (r *mockUserRepo) GetByID(_ context.Context, id uuid.UUID) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.users[id]
	if !ok {
		return nil, nil
	}
	return u, nil
}

func (r *mockUserRepo) GetByEmail(_ context.Context, email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, u := range r.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, nil
}

func (r *mockUserRepo) Update(_ context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[user.ID] = user
	return nil
}

// mockAuthProviderRepo implements port.AuthProviderRepository in memory.
type mockAuthProviderRepo struct {
	mu    sync.RWMutex
	links map[uuid.UUID]*domain.AuthProviderLink
}

func newMockAuthProviderRepo() *mockAuthProviderRepo {
	return &mockAuthProviderRepo{links: make(map[uuid.UUID]*domain.AuthProviderLink)}
}

func (r *mockAuthProviderRepo) Create(_ context.Context, link *domain.AuthProviderLink) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.links[link.ID] = link
	return nil
}

func (r *mockAuthProviderRepo) GetByProviderUID(_ context.Context, provider domain.AuthProvider, uid string) (*domain.AuthProviderLink, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, l := range r.links {
		if l.Provider == provider && l.ProviderUID == uid {
			return l, nil
		}
	}
	return nil, nil
}

func (r *mockAuthProviderRepo) GetByUserID(_ context.Context, userID uuid.UUID) ([]domain.AuthProviderLink, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.AuthProviderLink
	for _, l := range r.links {
		if l.UserID == userID {
			result = append(result, *l)
		}
	}
	return result, nil
}

// mockRefreshTokenRepo implements port.RefreshTokenRepository in memory.
type mockRefreshTokenRepo struct {
	mu     sync.RWMutex
	tokens map[uuid.UUID]*domain.RefreshTokenRecord
}

func newMockRefreshTokenRepo() *mockRefreshTokenRepo {
	return &mockRefreshTokenRepo{tokens: make(map[uuid.UUID]*domain.RefreshTokenRecord)}
}

func (r *mockRefreshTokenRepo) Create(_ context.Context, token *domain.RefreshTokenRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tokens[token.ID] = token
	return nil
}

func (r *mockRefreshTokenRepo) GetByToken(_ context.Context, token string) (*domain.RefreshTokenRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, t := range r.tokens {
		if t.Token == token {
			return t, nil
		}
	}
	return nil, nil
}

func (r *mockRefreshTokenRepo) RevokeByUserID(_ context.Context, userID uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, t := range r.tokens {
		if t.UserID == userID {
			t.Revoked = true
		}
	}
	return nil
}

func (r *mockRefreshTokenRepo) RevokeByToken(_ context.Context, token string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, t := range r.tokens {
		if t.Token == token {
			t.Revoked = true
		}
	}
	return nil
}

func (r *mockRefreshTokenRepo) DeleteExpired(_ context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for id, t := range r.tokens {
		if t.ExpiresAt.Before(time.Now()) {
			delete(r.tokens, id)
		}
	}
	return nil
}

// mockRecipientRepo implements port.RecipientRepository in memory.
type mockRecipientRepo struct {
	mu         sync.RWMutex
	recipients map[uuid.UUID]*domain.Recipient
}

func newMockRecipientRepo() *mockRecipientRepo {
	return &mockRecipientRepo{recipients: make(map[uuid.UUID]*domain.Recipient)}
}

func (r *mockRecipientRepo) Create(_ context.Context, rec *domain.Recipient) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.recipients[rec.ID] = rec
	return nil
}

func (r *mockRecipientRepo) GetByID(_ context.Context, id uuid.UUID) (*domain.Recipient, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rec, ok := r.recipients[id]
	if !ok {
		return nil, nil
	}
	return rec, nil
}

func (r *mockRecipientRepo) ListByUserID(_ context.Context, userID uuid.UUID) ([]domain.Recipient, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.Recipient
	for _, rec := range r.recipients {
		if rec.UserID == userID {
			result = append(result, *rec)
		}
	}
	return result, nil
}

func (r *mockRecipientRepo) Update(_ context.Context, rec *domain.Recipient) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.recipients[rec.ID] = rec
	return nil
}

func (r *mockRecipientRepo) Delete(_ context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.recipients, id)
	return nil
}

func (r *mockRecipientRepo) BulkDelete(_ context.Context, _ uuid.UUID, ids []uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, id := range ids {
		delete(r.recipients, id)
	}
	return nil
}

// mockSocialVerifier implements port.SocialVerifier.
type mockSocialVerifier struct{}

func (v *mockSocialVerifier) VerifyGoogleToken(_ context.Context, _ string) (string, string, string, error) {
	return "google@example.com", "Google User", "google-sub-123", nil
}

func (v *mockSocialVerifier) VerifyAppleToken(_ context.Context, _ string) (string, string, error) {
	return "apple@example.com", "apple-sub-123", nil
}
