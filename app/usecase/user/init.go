package user

import (
	"context"
	"go-api-project/domain"
	"time"
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
		return
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
	res, err = uc.UserRepo.GetByUsername(ctx, email)
	if err != nil {
		return
	}

	return
}

func (uc *UserUsecase) Store(c context.Context, m *domain.User) (err error) {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	existUser, _ := uc.GetByEmail(ctx, m.Email)
	if existUser != (domain.User{}) {
		return domain.ErrConflict
	}

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
