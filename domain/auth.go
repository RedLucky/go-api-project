package domain

import (
	"context"

	"github.com/golang-jwt/jwt"
)

// auth model
type Auth struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type JwtCustomClaims struct {
	ID int64 `json:"id"`
	jwt.StandardClaims
}

// AuthUsecase represent the authentication usecases
type AuthUsecase interface {
	Authenticate(ctx context.Context, email, password string) (string, error)
	SignUp(ctx context.Context, user *User) error
	// CreateVerifyEmail(ctx context.Context, email string) error
	// VerifyEmail(ctx context.Context, token string) error
	// CreateResetPassword(ctx context.Context, email string) error
	// VerifyResetPassword(ctx context.Context, token string) error
	// ResetPassword(ctx context.Context, password, confirm_password, token string) error
	// GenerateAccessToken(ctx context.Context, user *User) error
	// GenerateRefreshToken(ctx context.Context, user *User) error
	// Logout(ctx context.Context, token string) error
}

// AuthRepository represent the authentication repository contract
type AuthRepository interface {
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserByUsername(ctx context.Context, email string) (User, error)
	RegisterUser(ctx context.Context, user *User) error
	// CreateVerifyEmail(ctx context.Context, email string) error
	// VerifyEmail(ctx context.Context, token string) error
	// CreateResetPassword(ctx context.Context, email string) error
	// VerifyResetPassword(ctx context.Context, token string) error
	// ResetPassword(ctx context.Context, password, confirm_password, token string) error
}
