package repository

import (
	"context"
	"errors"
	"time"

	"github.com/RedLucky/potongin/domain"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type GeneratedUrlRepository struct {
	Mysql *gorm.DB
}

func NewGeneratedUrlRepository(conn *gorm.DB) domain.GeneratedUrlRepository {
	return &GeneratedUrlRepository{conn}
}

func (repo *GeneratedUrlRepository) InsertUrl(ctx context.Context, url *domain.GeneratedUrl) (err error) {
	err = repo.Mysql.Create(&url).Error
	if err != nil {
		return err
	}
	return
}

func (repo *GeneratedUrlRepository) UpdateUrl(ctx context.Context, url *domain.GeneratedUrl) (err error) {
	err = repo.Mysql.Model(&domain.GeneratedUrl{}).Where("id = ?", url.ID).Updates(
		domain.GeneratedUrl{Source: url.Source, Generated: url.Generated}).Error
	return
}

func (repo *GeneratedUrlRepository) GetUrlByUserId(ctx context.Context, userId int64) (generateUrls []domain.GeneratedUrl, err error) {
	err = repo.Mysql.Model(&domain.GeneratedUrl{}).Where("user_id = ?", userId).Find(&generateUrls).Error
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return
}

func (repo *GeneratedUrlRepository) GetUrlById(ctx context.Context, urlId string) (generateUrl domain.GeneratedUrl, err error) {
	err = repo.Mysql.Model(&domain.GeneratedUrl{}).Where("id = ?", urlId).First(&generateUrl).Error
	if err != nil {
		logrus.Error(err)
		return domain.GeneratedUrl{}, err
	}
	return
}

func (repo *GeneratedUrlRepository) GetUrlByUrl(ctx context.Context, url string) (generateUrl domain.GeneratedUrl, err error) {
	err = repo.Mysql.Model(&domain.GeneratedUrl{}).Where("generated = ?", url).First(&generateUrl).Error
	if err != nil {
		logrus.Error(err)
		return domain.GeneratedUrl{}, err
	}
	return
}

func (repo *GeneratedUrlRepository) IsExistUrlOrigin(ctx context.Context, urlOrigin string) (result bool, err error) {
	var generateUrl domain.GeneratedUrl
	err = repo.Mysql.Model(&domain.GeneratedUrl{}).First(&generateUrl, "source = ?", urlOrigin).Error
	if err != nil {
		logrus.Error(err)
		return false, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.Error(err)
		return false, err
	}
	result = true
	return
}

func (repo *GeneratedUrlRepository) IsExistUrlGenerated(ctx context.Context, urlGenerated string) (result bool, err error) {
	var generateUrl domain.GeneratedUrl
	err = repo.Mysql.Model(&domain.GeneratedUrl{}).Where(" ? between start_date and end_date and is_active = ? ", time.Now(), "Y").First(&generateUrl, "generated = ?", urlGenerated).Error
	if err != nil {
		logrus.Error(err)
		return false, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.Error(err)
		return false, err
	}
	result = true
	return
}

func (repo *GeneratedUrlRepository) CheckDoubleNameByUserId(ctx context.Context, name string, userId int64) (result bool, err error) {
	var generateUrl domain.GeneratedUrl
	err = repo.Mysql.Model(&domain.GeneratedUrl{}).Where("name = ? and user_id = ?", name, userId).First(&generateUrl).Error

	if err != nil {
		logrus.Error(err)
		return false, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.Error(err)
		return false, err
	}
	result = true
	return
}

func (repo *GeneratedUrlRepository) HitUrl(ctx context.Context, urlId, total int64) (err error) {
	err = repo.Mysql.Model(&domain.GeneratedUrl{}).Where("id = ?", urlId).Updates(
		domain.GeneratedUrl{TotalHits: total}).Error
	return
}
