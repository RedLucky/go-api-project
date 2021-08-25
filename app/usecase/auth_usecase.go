package usecase

import (
	"context"
	"time"

	"github.com/RedLucky/potongin/app/delivery/api/auth"
	"github.com/RedLucky/potongin/domain"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	AuthRepo       domain.AuthRepository
	contextTimeout time.Duration
}

// NewUserUsecase will create new an USerUsecase object representation of domain.UserUsecase interface
func NewAuthUsecase(repo domain.AuthRepository, timeout time.Duration) domain.AuthUsecase {
	return &AuthUsecase{
		AuthRepo:       repo,
		contextTimeout: timeout,
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

	return auth.CreateToken(&user)
}

// private function
func verifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
