package repository

import (
	"context"
	"errors"
	"time"

	"github.com/RedLucky/potongin/domain"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type AuthRepository struct {
	Mysql *gorm.DB
}

func NewAuthRepository(Conn *gorm.DB) domain.AuthRepository {
	return &AuthRepository{Conn}
}

func (m *AuthRepository) RegisterUser(ctx context.Context, a *domain.User) (err error) {
	err = m.Mysql.Create(&a).Error
	if err != nil {
		return err
	}

	return
}
func (r *AuthRepository) GetUserByEmail(ctx context.Context, email string) (user domain.User, err error) {

	err = r.Mysql.Model(&domain.User{}).Where("email = ?", email).First(&user).Error
	if err != nil {
		logrus.Error(err)
		return domain.User{}, err
	}

	return user, nil
}

func (r *AuthRepository) GetUserByUsername(ctx context.Context, username string) (user domain.User, err error) {

	err = r.Mysql.Model(&domain.User{}).Where("username = ?", username).First(&user).Error
	if err != nil {
		logrus.Error(err)
		return domain.User{}, err
	}

	return user, nil
}

func (r *AuthRepository) IsExistEmail(email string) (result domain.User, err error) {
	err = r.Mysql.Model(&domain.User{}).Where("email = ?", email).First(&result).Error
	if err != nil {
		logrus.Error(err)
		return result, err
	}

	return
}

func (r *AuthRepository) IsVerifiedEmail(email string) (results bool, err error) {
	var user domain.User
	err = r.Mysql.Model(&domain.User{}).Where("email = ? and email_verified = ?", email, "Y").First(&user).Error
	if err != nil {
		logrus.Error(err)
		return false, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.Error(err)
		return false, err
	}

	return
}

func (r *AuthRepository) IsExpiredTokenEmail(token string) (result bool, err error) {
	return
}

func (r *AuthRepository) IsExistTokenEmail(token string) (result domain.VerifyEmail, err error) {
	var verifyEmail domain.VerifyEmail
	err = r.Mysql.Model(&domain.VerifyEmail{}).Where("token = ? ", token).First(&verifyEmail).Error
	if err != nil {
		logrus.Error(err)
		return verifyEmail, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.Error(err)
		return verifyEmail, err
	}
	return verifyEmail, nil
}

func (r *AuthRepository) VerifyTokenEmail(ctx context.Context, token string) error {
	err := r.Mysql.Model(&domain.VerifyEmail{}).Where("token = ?", token).Updates(
		domain.VerifyEmail{Verified: "Y", VerifiedAt: time.Now()}).Error

	return err
}

func (r *AuthRepository) VerifyTokenAccount(ctx context.Context, userId int64) error {
	err := r.Mysql.Model(&domain.User{}).Where("id = ?", userId).Updates(
		domain.User{EmailVerified: "Y"}).Error

	return err
}

func (r *AuthRepository) DeletePreviousVerifyEmail(userId int64) error {
	err := r.Mysql.Model(&domain.VerifyEmail{}).Where("user_id = ?", userId).Delete(&domain.VerifyEmail{}).Error
	return err
}

func (r *AuthRepository) CreateVerifyEmail(verifyEmail *domain.VerifyEmail) error {
	err := r.Mysql.Create(&verifyEmail).Error
	if err != nil {
		return err
	}
	return nil
}
