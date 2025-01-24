package services

import (
	"context"
	"errors"
	"fmt"
	"fp_kata/common/utils"
	"fp_kata/internal/models"
	"fp_kata/pkg/log"
	"sync"
)

const compAuthenticationService = "AuthenticationService"

// AuthService is the interface for the authentication service.
type AuthService interface {
	// GenerateAuthToken takes a User object, generates an auth token, and stores the relationship in memory.
	GenerateAuthToken(ctx context.Context, user models.User) (string, error)

	// GetUserIDByToken takes an auth token and retrieves the associated User ID.
	GetUserIDByToken(ctx context.Context, authToken string) (int, error)
}

// authService is the implementation of AuthService.
type authService struct {
	tokenStorage map[string]int
	mutex        sync.Mutex // Mutex to ensure thread-safe access to the tokenStorage map.
}

// NewAuthService creates a new instance of AuthService.
func NewAuthService() AuthService {
	return &authService{
		tokenStorage: make(map[string]int),
	}
}

// GenerateAuthToken generates an authentication token for the provided user and stores it in memory.
func (s *authService) GenerateAuthToken(ctx context.Context, user models.User) (string, error) {
	logger := log.GetLogger(ctx)
	utils.LogAction(ctx, compAuthenticationService, "GenerateAuthToken")

	if user.ID == 0 {
		return "", errors.New("invalid user")
	}
	// Generate a simple token (in real-world apps, use a more secure method like UUIDs or hashes).
	authToken := generateToken(user.ID)

	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Store the token and associated user ID.
	s.tokenStorage[authToken] = user.ID
	logger.Debug().Str(log.Comp, compAuthenticationService).Str("func", "GenerateAuthToken").Int("user_id", user.ID).Str("auth_token", authToken).Send()
	return authToken, nil
}

// GetUserIDByToken retrieves the user ID associated with the given auth token.
func (s *authService) GetUserIDByToken(ctx context.Context, authToken string) (int, error) {
	logger := log.GetLogger(ctx)
	utils.LogAction(ctx, compAuthenticationService, "GetUserIDByToken")

	if authToken == "" {
		return 0, errors.New("invalid auth token")
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()

	logger.Debug().Str(log.Comp, compAuthenticationService).Str("func", "GetUserIDByToken").Str("auth_token", authToken).Send()

	userID, exists := s.tokenStorage[authToken]
	logger.Debug().Str(log.Comp, compAuthenticationService).Str("func", "GetUserIDByToken").Str("auth_token", authToken).Int("user_id", userID).Bool("exists", exists).Send()
	if !exists {
		return 0, errors.New("auth token not found")
	}

	return userID, nil
}

// Helper function to simulate token generation.
func generateToken(userID int) string {
	// In a real-world application, replace this with secure token generation logic.
	return fmt.Sprintf("token_%d", userID)
}
