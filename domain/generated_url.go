package domain

import (
	"context"
	"time"

	"github.com/gomodule/redigo/redis"
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
	GetUrlByUserId(ctx context.Context, userId int64) ([]GeneratedUrl, error)
	GetUrlById(ctx context.Context, urlId string) (GeneratedUrl, error)
	HitUrl(ctx context.Context, generateUrl string) (originUrl string, err error)
}

type GeneratedUrlRepository interface {
	InsertUrl(ctx context.Context, url *GeneratedUrl) error
	UpdateUrl(ctx context.Context, url *GeneratedUrl) error
	GetUrlByUserId(ctx context.Context, userId int64) ([]GeneratedUrl, error)
	GetUrlById(ctx context.Context, urlId string) (GeneratedUrl, error)
	GetUrlByUrl(ctx context.Context, url string) (GeneratedUrl, error)
	IsExistUrlOrigin(ctx context.Context, urlOrigin string) (bool, error)
	IsExistUrlGenerated(ctx context.Context, urlGenerated string) (bool, error)
	CheckDoubleNameByUserId(ctx context.Context, name string, userId int64) (bool, error)
	HitUrl(ctx context.Context, urlId, total int64) error
	// using redis
	GetUrlFromCache(redisCon redis.Conn, generatedUrl string) (string, error)
	SetUrlToCache(redisCon redis.Conn, generatedUrl, sourceUrl string) error
	SetUrlExpCache(redisCon redis.Conn, generatedUrl string, duration int) error
}
