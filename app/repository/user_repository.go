package repository

import (
	"github.com/RedLucky/potongin/domain"
	"github.com/jinzhu/gorm"

	"github.com/sirupsen/logrus"
)

type UserRepository struct {
	Mysql *gorm.DB
}

var (
	field []string = []string{"id", "email", "username", "name", "email_verified", "updated_at", "created_at"}
)

// NewUserRepository will create an object that represent the user.Repository interface
func NewUserRepository(Conn *gorm.DB) domain.UserRepository {
	return &UserRepository{Conn}
}

func (m *UserRepository) Fetch() (res []domain.User, err error) {
	err = m.Mysql.Model(&domain.User{}).Select(field).Limit(100).Find(&res).Error
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return
}
func (m *UserRepository) GetByID(id int64) (res domain.User, err error) {
	err = m.Mysql.Model(&domain.User{}).Select(field).Where("id = ?", id).First(&res).Error
	if err != nil {
		logrus.Error(err)
		return domain.User{}, err
	}

	return
}

func (m *UserRepository) GetByEmail(email string) (res domain.User, err error) {
	err = m.Mysql.Model(&domain.User{}).Where("email = ?", email).First(&res).Error
	if err != nil {
		logrus.Error(err)
		return domain.User{}, err
	}
	return
}

func (m *UserRepository) GetByUsername(username string) (res domain.User, err error) {
	err = m.Mysql.Model(&domain.User{}).Where("username = ?", username).First(&res).Error
	if err != nil {
		logrus.Error(err)
		return domain.User{}, err
	}
	return
}

func (m *UserRepository) Store(a *domain.User) (err error) {
	err = m.Mysql.Model(&domain.User{}).Create(&a).Error
	if err != nil {
		return err
	}

	return
}

func (m *UserRepository) Delete(id int64) (err error) {
	err = m.Mysql.Model(&domain.User{}).Where("id = ?", id).Delete(&domain.User{}).Error
	return

}
func (m *UserRepository) Update(ar *domain.User) (err error) {
	err = m.Mysql.Model(&domain.User{}).Where("id = ?", ar.ID).Updates(
		domain.User{Name: ar.Name}).Error

	return
}
