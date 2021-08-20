package repository

import (
	"context"

	"github.com/RedLucky/potongin/domain"
	"github.com/jinzhu/gorm"

	"github.com/sirupsen/logrus"
)

type UserRepository struct {
	Mysql *gorm.DB
}

// NewUserRepository will create an object that represent the user.Repository interface
func NewUserRepository(Conn *gorm.DB) domain.UserRepository {
	return &UserRepository{Conn}
}

func (m *UserRepository) Fetch(ctx context.Context) (res []domain.User, err error) {

	err = m.Mysql.Model(&domain.User{}).Limit(100).Find(&res).Error
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return
}
func (m *UserRepository) GetByID(ctx context.Context, id int64) (res domain.User, err error) {
	field := []string{"id", "email", "username", "name", "updated_at", "created_at"}
	err = m.Mysql.Model(&domain.User{}).Select(field).Where("id = ?", id).First(&res).Error
	if err != nil {
		logrus.Error(err)
		return domain.User{}, err
	}

	return
}

func (m *UserRepository) GetByEmail(ctx context.Context, email string) (res domain.User, err error) {
	err = m.Mysql.Model(&domain.User{}).Where("email = ?", email).First(&res).Error
	if err != nil {
		logrus.Error(err)
		return domain.User{}, err
	}
	return
}

func (m *UserRepository) GetByUsername(ctx context.Context, username string) (res domain.User, err error) {
	err = m.Mysql.Model(&domain.User{}).Where("username = ?", username).First(&res).Error
	if err != nil {
		logrus.Error(err)
		return domain.User{}, err
	}
	return
}

func (m *UserRepository) Store(ctx context.Context, a *domain.User) (err error) {
	err = m.Mysql.Create(&a).Error
	if err != nil {
		return err
	}

	return
}

func (m *UserRepository) Delete(ctx context.Context, id int64) (err error) {
	err = m.Mysql.Model(&domain.User{}).Where("id = ?", id).Delete(&domain.User{}).Error
	return

}
func (m *UserRepository) Update(ctx context.Context, ar *domain.User) (err error) {
	err = m.Mysql.Model(&domain.User{}).Where("id = ?", ar.ID).Updates(
		domain.User{Name: ar.Name}).Error

	return
}
