package repository_test

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/RedLucky/potongin/app/repository"
	"github.com/RedLucky/potongin/domain"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestUserRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	gdb, _ := gorm.Open("mysql", db)
	userRepo := repository.NewUserRepository(gdb)

	mock.ExpectQuery(
		"SELECT(.*)").
		WithArgs(5).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "username", "email", "name", "email_verified", "updated_at", "created_at"}).
				AddRow(1, "LFR", "lucky@kryptopos.com", "Lucky Fernanda R", "Y", time.Now(), time.Now()))
	res, err := userRepo.GetByID(5)

	require.NoError(t, err)
	assert.NotNil(t, res)
	// require.Equal(t, res, products[0])
}

func TestUserRepository_Fetch(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	gdb, _ := gorm.Open("mysql", db)
	userRepo := repository.NewUserRepository(gdb)

	mock.ExpectQuery(
		"SELECT id, email, username, name, email_verified, updated_at, created_at FROM `users`").
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "username", "email", "name", "email_verified", "updated_at", "created_at"}).
				AddRow(1, "LFR", "lucky@kryptopos.com", "Lucky Fernanda R", "Y", time.Now(), time.Now()))
	res, err := userRepo.Fetch()

	require.NoError(t, err)
	assert.NotNil(t, res)
}

func TestUserRepository_Store(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	gdb, _ := gorm.Open("mysql", db)
	userRepo := repository.NewUserRepository(gdb)
	queryInsert := "INSERT INTO `users` (`username`,`email`,`password`,`name`,`email_verified`,`updated_at`,`created_at`) VALUES (?,?,?,?,?,?,?)"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("LFR123456"), bcrypt.DefaultCost)
	user := &domain.User{
		Username:      "LFR123",
		Email:         "lucky@kryptopos.com",
		Password:      string(hashedPassword),
		Name:          "Lucky Fernanda R",
		EmailVerified: "N",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	mock.ExpectBegin()
	mock.ExpectExec(queryInsert).WithArgs(
		user.Username, user.Email, user.Password, user.Name, user.EmailVerified, user.UpdatedAt, user.CreatedAt).WillReturnResult(sqlmock.NewResult(12, 1))
	mock.ExpectCommit()

	err = userRepo.Store(user)
	require.NoError(t, err)
	assert.Equal(t, int64(12), user.ID)

}

func TestUserRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	gdb, _ := gorm.Open("mysql", db)
	userRepo := repository.NewUserRepository(gdb)
	queryUpdate := "UPDATE  `users` SET `name` = ?, `updated_at` = ? WHERE (id = ?)"
	user := &domain.User{
		ID:        1,
		Name:      "Lucky Fernanda R",
		UpdatedAt: time.Now(),
	}
	mock.ExpectBegin()
	mock.ExpectExec(queryUpdate).WithArgs(
		user.Name, user.UpdatedAt, user.ID).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = userRepo.Update(user)
	require.NoError(t, err)
	assert.Equal(t, int64(1), user.ID)

}
