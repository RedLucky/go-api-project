package repository

import (
	"context"

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
