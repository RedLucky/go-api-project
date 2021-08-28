package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/RedLucky/potongin/app/delivery/api/auth"
	"github.com/RedLucky/potongin/domain"
	"github.com/gomodule/redigo/redis"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	AuthRepo       domain.AuthRepository
	contextTimeout time.Duration
	RedisPool      *redis.Pool
}

// NewUserUsecase will create new an USerUsecase object representation of domain.UserUsecase interface
func NewAuthUsecase(repo domain.AuthRepository, timeout time.Duration, redisPool *redis.Pool) domain.AuthUsecase {
	return &AuthUsecase{
		AuthRepo:       repo,
		contextTimeout: timeout,
		RedisPool:      redisPool,
	}
}

func (uc *AuthUsecase) SignUp(c context.Context, user *domain.User) (err error) {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	existEmail, _ := uc.AuthRepo.GetUserByEmail(ctx, user.Email)
	if existEmail != (domain.User{}) {
		return domain.ErrEmailExist
	}

	existUsername, _ := uc.AuthRepo.GetUserByUsername(ctx, user.Username)
	if existUsername != (domain.User{}) {
		return domain.ErrAccountExist
	}

	hashedPassword, err := hash(user.Password)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	err = uc.AuthRepo.RegisterUser(ctx, user)
	return
}

func (uc *AuthUsecase) Authenticate(c context.Context, email, password string) (token domain.JwtResults, err error) {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	user := domain.User{}

	user, err = uc.AuthRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return domain.JwtResults{}, err
	}

	err = verifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return domain.JwtResults{}, domain.ErrPassword
	}
	token, err = auth.CreateToken(&user)
	if err != nil {
		return domain.JwtResults{}, domain.ErrInternalServerError
	}
	// set on redis cache pool
	conn := uc.RedisPool.Get()
	defer conn.Close()
	if err = auth.SaveToken(conn, user, token); err != nil {
		fmt.Println(err)
		return domain.JwtResults{}, domain.ErrInternalServerError
	}
	return
}

func (uc *AuthUsecase) GenerateNewAccessToken(c echo.Context) (token domain.JwtResults, err error) {
	_, cancel := context.WithTimeout(c.Request().Context(), uc.contextTimeout)
	defer cancel()
	var user domain.User

	// validate token
	claims, err := auth.RefreshTokenValid(c.Request())
	if err != nil {
		return domain.JwtResults{}, err
	}

	// get uuid token from jwt
	res, ok := claims["refresh_uuid"].(string)
	if !ok {
		return domain.JwtResults{}, err
	}
	// get user_id to redis
	conn := uc.RedisPool.Get()
	defer conn.Close()
	userId, err := auth.GetTokenFromRedis(conn, res)
	if err != nil {
		return domain.JwtResults{}, err
	}
	user.ID = userId
	// generate new refresh token
	token, err = auth.CreateToken(&user)
	// set on redis cache pool
	if err = auth.SaveToken(conn, user, token); err != nil {
		return domain.JwtResults{}, err
	}
	err = auth.DeleteRefreshTokenRedis(conn, res)
	return
}

// private function
func verifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
