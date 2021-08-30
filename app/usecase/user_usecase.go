package usecase

import (
	"context"
	"time"

	"github.com/RedLucky/potongin/domain"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	UserRepo       domain.UserRepository
	contextTimeout time.Duration
}

// NewUserUsecase will create new an USerUsecase object representation of domain.UserUsecase interface
func NewUserUsecase(repo domain.UserRepository, timeout time.Duration) domain.UserUsecase {
	return &UserUsecase{
		UserRepo:       repo,
		contextTimeout: timeout,
	}
}

func (uc *UserUsecase) Fetch(c context.Context) (res []domain.User, err error) {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()

	res, err = uc.UserRepo.Fetch(ctx)
	if err != nil {
		return nil, err
	}

	return
}

func (uc *UserUsecase) GetByID(c context.Context, id int64) (res domain.User, err error) {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()

	res, err = uc.UserRepo.GetByID(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	return
}

func (uc *UserUsecase) Update(c context.Context, ar *domain.User) (err error) {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()

	ar.UpdatedAt = time.Now()
	return uc.UserRepo.Update(ctx, ar)
}

func (uc *UserUsecase) GetByUsername(c context.Context, username string) (res domain.User, err error) {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	res, err = uc.UserRepo.GetByUsername(ctx, username)
	if err != nil {
		return
	}

	return
}

func (uc *UserUsecase) GetByEmail(c context.Context, email string) (res domain.User, err error) {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	res, err = uc.UserRepo.GetByEmail(ctx, email)
	if err != nil {
		return
	}

	return
}

func (uc *UserUsecase) Store(c context.Context, m *domain.User) (err error) {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	existEmail, _ := uc.GetByEmail(ctx, m.Email)
	if existEmail != (domain.User{}) {
		return domain.ErrEmailExist
	}

	existUsername, _ := uc.GetByUsername(ctx, m.Username)
	if existUsername != (domain.User{}) {
		return domain.ErrAccountExist
	}

	hashedPassword, err := hash(m.Password)
	if err != nil {
		return err
	}

	m.Password = string(hashedPassword)
	m.EmailVerified = "N"
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()

	err = uc.UserRepo.Store(ctx, m)
	return
}

func (uc *UserUsecase) Delete(c context.Context, id int64) (err error) {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	existUser, _ := uc.UserRepo.GetByID(ctx, id)
	if err != nil {
		return
	}
	if existUser == (domain.User{}) {
		return domain.ErrNotFound
	}
	return uc.UserRepo.Delete(ctx, id)
}

// private function

func hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}
