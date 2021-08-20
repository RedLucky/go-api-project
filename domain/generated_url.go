package domain

import (
	"context"
	"time"
)

// define models
type GeneratedUrl struct {
	ID        int64     `json:"id" gorm:"primary_key;auto_increment"`
	UserId    int64     `json:"user_id"`
	Name      string    `json:"name" validate:"required"`
	Source    string    `json:"source_link" validate:"required"`
	Generated string    `json:"generated_link" validate:"required"`
	TotalHits int64     `json:"total_hits"`
	IsActive  string    `json:"is_active"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GeneratedUrlUsecase interface {
	CreateUrl(ctx context.Context, url *GeneratedUrl) error
	UpdateUrl(ctx context.Context, url *GeneratedUrl) error
	GetUrlByUserId(ctx context.Context, userId string) ([]GeneratedUrl, error)
	GetUrlById(ctx context.Context, urlId string) (GeneratedUrl, error)
	HitUrl(ctx context.Context, generateUrl string) (originUrl string, err error)
}

type GeneratedUrlRepository interface {
	InsertUrl(ctx context.Context, url *GeneratedUrl) error
	UpdateUrl(ctx context.Context, url *GeneratedUrl) error
	GetUrlByUserId(ctx context.Context, userId string) ([]GeneratedUrl, error)
	GetUrlById(ctx context.Context, urlId string) (GeneratedUrl, error)
	IsExistUrlOrigin(ctx context.Context, urlOrigin string) (bool, error)
	IsExistUrlGenerated(ctx context.Context, urlGenerated string) (bool, error)
	CheckDoubleNameByUserId(ctx context.Context, name, userId string) (bool, error)
	HitUrl(ctx context.Context, urlId string, total int64) error
}
