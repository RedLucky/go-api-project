package usecase

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/RedLucky/potongin/domain"
	"github.com/gomodule/redigo/redis"
	"github.com/spf13/viper"
)

type GeneratedUrlUsecase struct {
	GeneratedRepo  domain.GeneratedUrlRepository
	contextTimeout time.Duration
	RedisPool      *redis.Pool
}

type UrlCache struct {
	GenerateUrlId int64  `redis:"generate_url_id"`
	SourceUrl     string `redis:"source_url"`
}

func NewGeneratedUrlUsecase(repo domain.GeneratedUrlRepository, timeout time.Duration, redis *redis.Pool) domain.GeneratedUrlUsecase {
	return &GeneratedUrlUsecase{
		GeneratedRepo:  repo,
		contextTimeout: timeout,
		RedisPool:      redis,
	}
}

func (gu *GeneratedUrlUsecase) CreateUrl(ctx context.Context, url *domain.GeneratedUrl) (err error) {
	ctx, cancel := context.WithTimeout(ctx, gu.contextTimeout)
	defer cancel()
	// check url contains https or http
	if !strings.Contains(url.Source, "https://") && !strings.Contains(url.Source, "http://") {
		return domain.ErrUrlNotFound
	}
	// check is url valid
	reader := strings.NewReader(`{}`)
	request, err := http.NewRequest("GET", url.Source, reader)
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return domain.ErrUrlNotFound
	}
	defer resp.Body.Close()

	url.Source = strings.TrimPrefix(url.Source, "https://")
	url.Source = strings.TrimPrefix(url.Source, "http://")
	url.CreatedAt = time.Now()
	url.UpdatedAt = time.Now()
	url.StartDate = time.Now()
	url.IsActive = "Y"

	existOriginUrl, _ := gu.GeneratedRepo.IsExistUrlOrigin(ctx, url.Source)
	if existOriginUrl {
		return domain.ErrUrlOriginExist
	}

	existGeneratedUrl, _ := gu.GeneratedRepo.IsExistUrlGenerated(ctx, url.Generated)
	if existGeneratedUrl {
		return domain.ErrUrlGeneratedExist
	}

	existNameByUserId, _ := gu.GeneratedRepo.CheckDoubleNameByUserId(ctx, url.Name, url.UserId)
	if existNameByUserId {
		return domain.ErrNameIsExist
	}

	err = gu.GeneratedRepo.InsertUrl(ctx, url)

	return
}

func (gu *GeneratedUrlUsecase) UpdateUrl(ctx context.Context, url *domain.GeneratedUrl) (err error) {
	// check url contains https or http
	if !strings.Contains(url.Source, "https://") && !strings.Contains(url.Source, "http://") {
		return domain.ErrUrlNotFound
	}
	// check is url valid
	reader := strings.NewReader(`{}`)
	request, err := http.NewRequest("GET", url.Source, reader)
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return domain.ErrUrlNotFound
	}
	defer resp.Body.Close()

	url.Source = strings.TrimPrefix(url.Source, "https://")
	url.Source = strings.TrimPrefix(url.Source, "http://")

	existOriginUrl, _ := gu.GeneratedRepo.IsExistUrlOrigin(ctx, url.Source)
	if existOriginUrl {
		return domain.ErrUrlOriginExist
	}

	existGeneratedUrl, _ := gu.GeneratedRepo.IsExistUrlGenerated(ctx, url.Generated)
	if existGeneratedUrl {
		return domain.ErrUrlGeneratedExist
	}

	existNameByUserId, _ := gu.GeneratedRepo.CheckDoubleNameByUserId(ctx, url.Name, url.UserId)
	if existNameByUserId {
		return domain.ErrNameIsExist
	}

	err = gu.GeneratedRepo.UpdateUrl(ctx, url)
	// check on redis cache pool
	conn := gu.RedisPool.Get()
	defer conn.Close()
	_, err = conn.Do("HDEL", url.Generated, "source_url")
	if err != nil {
		return domain.ErrInternalServerError
	}
	return
}

func (gu *GeneratedUrlUsecase) GetUrlByUserId(ctx context.Context, userId int64) (results []domain.GeneratedUrl, err error) {
	ctx, cancel := context.WithTimeout(ctx, gu.contextTimeout)
	defer cancel()

	results, err = gu.GeneratedRepo.GetUrlByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}
	return
}

func (gu *GeneratedUrlUsecase) GetUrlById(ctx context.Context, urlId string) (results domain.GeneratedUrl, err error) {
	ctx, cancel := context.WithTimeout(ctx, gu.contextTimeout)
	defer cancel()

	results, err = gu.GeneratedRepo.GetUrlById(ctx, urlId)
	if err != nil {
		return domain.GeneratedUrl{}, err
	}
	return
}

func (gu *GeneratedUrlUsecase) HitUrl(ctx context.Context, generateUrl string) (results string, err error) {
	ctx, cancel := context.WithTimeout(ctx, gu.contextTimeout)
	defer cancel()

	// check on redis cache pool
	conn := gu.RedisPool.Get()
	defer conn.Close()

	res, err := gu.GeneratedRepo.GetUrlFromCache(conn, generateUrl)
	if err == redis.ErrNil {
		existGeneratedUrl, _ := gu.GeneratedRepo.IsExistUrlGenerated(ctx, generateUrl)
		if !existGeneratedUrl {
			return "", domain.ErrUrlNotFound
		}

		results, err = gu.updateTotalHits(ctx, generateUrl)

		err = gu.GeneratedRepo.SetUrlToCache(conn, generateUrl, results)
		if err != nil {
			return "", domain.ErrInternalServerError
		}
		// 60 is seconds
		err = gu.GeneratedRepo.SetUrlExpCache(conn, generateUrl, viper.GetInt(`redis.exp_hit_url`))
		if err != nil {
			return "", domain.ErrInternalServerError
		}
		return
	} else if err != nil {
		return "", domain.ErrInternalServerError
	} else {
		results, err = gu.updateTotalHits(ctx, generateUrl)
		return res, nil
	}

}

func (gu *GeneratedUrlUsecase) updateTotalHits(ctx context.Context, generateUrl string) (results string, err error) {
	// check url punya siapa
	ownerUrl, err := gu.GeneratedRepo.GetUrlByUrl(ctx, generateUrl)
	if err != nil {
		return "", domain.ErrUrlNotFound
	}

	// update total hit nya
	ownerUrl.TotalHits++
	results = ownerUrl.Source

	err = gu.GeneratedRepo.HitUrl(ctx, ownerUrl.ID, ownerUrl.TotalHits)
	if err != nil {
		return "", domain.ErrUrlGeneratedExist
	}
	return
}
