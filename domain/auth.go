package domain

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

// auth model
type Auth struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type JwtResults struct {
	AccessUUID   string `json:"access_uuid"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	RefreshUUID  string `json:"refresh_uuid"`
	AccessExp    int64  `json:"access_exp"`
	RefreshExp   int64  `json:"refresh_exp"`
}

type JwtCustomClaims struct {
	AccessUUID  string `json:"access_uuid"`
	RefreshUUID string `json:"refresh_uuid"`
	jwt.StandardClaims
}

type VerifyEmail struct {
	ID         int64     `json:"id"`
	Token      string    `json:"token"`
	UserId     int64     `json:"user_id"`
	Verified   string    `json:"verified"`
	CreatedAt  time.Time `json:"created_at"`
	VerifiedAt time.Time `json:"verified_at"`
}

// AuthUsecase represent the authentication usecases
type AuthUsecase interface {
	Authenticate(ctx context.Context, email, password string) (JwtResults, error)
	SignUp(ctx context.Context, user *User) error
	CreateVerifyEmail(ctx context.Context, email string) (encodedString string, err error)
	VerifyEmail(ctx context.Context, token string) error
	// CreateResetPassword(ctx context.Context, email string) error
	// VerifyResetPassword(ctx context.Context, token string) error
	// ResetPassword(ctx context.Context, password, confirm_password, token string) error
	GenerateNewAccessToken(ctx echo.Context) (JwtResults, error)
	Logout(accessToken, refreshToken string) error
}

// AuthRepository represent the authentication repository contract
type AuthRepository interface {
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserByUsername(ctx context.Context, email string) (User, error)
	RegisterUser(ctx context.Context, user *User) error
	IsExistEmail(email string) (result User, err error)
	IsVerifiedEmail(email string) (result bool, err error)
	IsExpiredTokenEmail(token string) (result bool, err error)
	IsExistTokenEmail(token string) (result VerifyEmail, err error)
	CreateVerifyEmail(verifyEmail *VerifyEmail) error
	DeletePreviousVerifyEmail(userId int64) error
	VerifyTokenEmail(ctx context.Context, token string) error
	VerifyTokenAccount(ctx context.Context, userId int64) error
	// CreateResetPassword(ctx context.Context, email string) error
	// VerifyResetPassword(ctx context.Context, token string) error
	// ResetPassword(ctx context.Context, password, confirm_password, token string) error
}
