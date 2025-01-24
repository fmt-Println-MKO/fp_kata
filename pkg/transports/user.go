package transports

import (
	"fmt"
	"fp_kata/internal/models"
	"math/rand"
	"strings"
	"time"
)

type UserResponse struct {
	ID       int               `json:"id"`
	Username string            `json:"username"`
	Email    string            `json:"email"`
	Password string            `json:"password"`
	Orders   []OrderResponse   `json:"orders,omitempty"`
	Payments []PaymentResponse `json:"payments,omitempty"`
}

type UserCreateRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ToUser maps a UserCreateRequest to a User model
func (ur *UserCreateRequest) ToUser() *models.User {

	if ur.Email == "" || ur.Password == "" {
		return nil
	}

	return &models.User{
		Username: extractName(ur),
		Email:    ur.Email,
		Password: ur.Password,
	}
}

func extractName(userRequest *UserCreateRequest) (username string) {
	if userRequest != nil && userRequest.Email != "" {
		parts := strings.Split(userRequest.Email, "@")
		if len(parts) > 1 {
			username = parts[0]
		} else {
			randGen := rand.New(rand.NewSource(time.Now().UnixNano()))
			username = fmt.Sprintf("User%d", randGen.Intn(999)+100)
		}
	} else {
		randGen := rand.New(rand.NewSource(time.Now().UnixNano()))
		username = fmt.Sprintf("User%d", randGen.Intn(90000)+10000)
	}
	return username
}

func MapToUserResponse(user models.User) *UserResponse {
	return &UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	}
}
