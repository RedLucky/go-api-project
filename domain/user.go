package domain

import (
	"context"
	"time"
)

// User ...
type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username" validate:"required"`
	Email     string    `json:"email" validate:"required"`
	Name      string    `json:"name" validate:"required"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

// UserUsecase represent the article's usecases
type UserUsecase interface {
	Fetch(ctx context.Context) ([]User, error)
	GetByID(ctx context.Context, id int64) (User, error)
	Update(ctx context.Context, ar *User) error
	GetByUsername(ctx context.Context, username string) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	Store(context.Context, *User) error
	Delete(ctx context.Context, id int64) error
}

// UserRepository represent the User's repository contract
type UserRepository interface {
	Fetch(ctx context.Context) (res []User, err error)
	GetByID(ctx context.Context, id int64) (User, error)
	GetByUsername(ctx context.Context, username string) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	Update(ctx context.Context, ar *User) error
	Store(ctx context.Context, a *User) error
	Delete(ctx context.Context, id int64) error
}
