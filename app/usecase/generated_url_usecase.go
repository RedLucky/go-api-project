package usecase

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/RedLucky/potongin/domain"
	"github.com/spf13/viper"
)

type GeneratedUrlUsecase struct {
	GeneratedRepo  domain.GeneratedUrlRepository
	contextTimeout time.Duration
}

func NewGeneratedUrlUsecase(repo domain.GeneratedUrlRepository, timeout time.Duration) domain.GeneratedUrlUsecase {
	return &GeneratedUrlUsecase{
		GeneratedRepo:  repo,
		contextTimeout: timeout,
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

	url.Generated = viper.GetString(`server.url_prefix`) + url.Generated
	existGeneratedUrl, _ := gu.GeneratedRepo.IsExistUrlGenerated(ctx, url.Generated)
	if existGeneratedUrl {
		return domain.ErrUrlGeneratedExist
	}

	err = gu.GeneratedRepo.InsertUrl(ctx, url)

	return
}

func (gu *GeneratedUrlUsecase) UpdateUrl(ctx context.Context, url *domain.GeneratedUrl) (err error) {

	return
}

func (gu *GeneratedUrlUsecase) GetUrlByUserId(ctx context.Context, userId string) (results []domain.GeneratedUrl, err error) {

	return
}

func (gu *GeneratedUrlUsecase) GetUrlById(ctx context.Context, urlId string) (results domain.GeneratedUrl, err error) {

	return
}

func (gu *GeneratedUrlUsecase) HitUrl(ctx context.Context, generateUrl string) (results string, err error) {

	return
}