package usecase

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/RedLucky/potongin/app/delivery/api/auth"
	"github.com/RedLucky/potongin/domain"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
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
	if err != nil {
		return err
	}

	encodedString, err := uc.CreateVerifyEmail(ctx, user.Email)
	fmt.Println(encodedString)
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
	// check is verified email?
	ok, err := uc.AuthRepo.IsVerifiedEmail(email)
	if err != nil && !ok {
		return domain.JwtResults{}, domain.ErrorEmailNotVerified
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
	jwt := auth.ExtractToken(c.Request())
	claims, err := auth.TokenValid(jwt, auth.RefreshToken)
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
	err = auth.DeleteTokenRedis(conn, res)
	return
}

func (uc *AuthUsecase) Logout(accessToken string, refreshToken string) (err error) {
	access, err := auth.TokenValid(accessToken, auth.AccessToken)
	if err != nil {
		return err
	}

	refresh, err := auth.TokenValid(refreshToken, auth.RefreshToken)
	if err != nil {
		return err
	}

	conn := uc.RedisPool.Get()
	defer conn.Close()
	err = auth.DeleteTokenRedis(conn, access["access_uuid"].(string))
	err = auth.DeleteTokenRedis(conn, refresh["refresh_uuid"].(string))

	return
}

func (uc *AuthUsecase) CreateVerifyEmail(ctx context.Context, email string) (encodedString string, err error) {
	_, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()
	// check email exist or not
	user, err := uc.AuthRepo.IsExistEmail(email)
	if err != nil {
		return "", err
	}

	if user == (domain.User{}) {
		return "", domain.ErrEmailExist
	}

	// check is verified email?
	ok, err := uc.AuthRepo.IsVerifiedEmail(email)
	if err == nil && ok {
		return "", err // email has verified
	}
	// delete previous token
	err = uc.AuthRepo.DeletePreviousVerifyEmail(user.ID)
	if err != nil {
		return "", err
	}
	// create verify email
	var verifyEmail domain.VerifyEmail
	verifyEmail.Token = uuid.New().String()
	verifyEmail.UserId = user.ID
	verifyEmail.Verified = "N"
	verifyEmail.CreatedAt = time.Now()
	err = uc.AuthRepo.CreateVerifyEmail(&verifyEmail)
	if err != nil {
		return "", err
	}
	encodedString = base64.StdEncoding.EncodeToString([]byte(verifyEmail.Token))
	return
}

func (uc *AuthUsecase) VerifyEmail(ctx context.Context, token string) error {
	_, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()
	var tokenEmail domain.VerifyEmail
	var decodedByte, _ = base64.StdEncoding.DecodeString(token)
	var decodeToken = string(decodedByte)

	tokenEmail, err := uc.AuthRepo.IsExistTokenEmail(decodeToken)
	if err != nil {
		return err
	}

	// soon check the status is verified or not. if verified, u can continue next step.

	err = uc.AuthRepo.VerifyTokenEmail(ctx, decodeToken)
	if err != nil {
		return err
	}

	err = uc.AuthRepo.VerifyTokenAccount(ctx, tokenEmail.UserId)
	if err != nil {
		return err
	}
	return nil
}

// private function
func verifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
