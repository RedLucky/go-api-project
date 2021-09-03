package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/RedLucky/potongin/app/usecase"
	"github.com/RedLucky/potongin/domain"
	"github.com/RedLucky/potongin/domain/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserUsecase_Fetch(t *testing.T) {
	repository := new(mocks.UserRepository)
	usersMock := []domain.User{
		{ID: 1,
			Username:      "LFR",
			Email:         "lucky@kryptopos.com",
			Name:          "Lucky Fernanda",
			EmailVerified: "N",
			UpdatedAt:     time.Now(),
			CreatedAt:     time.Now(),
		}, {ID: 2,
			Username:      "LFR16",
			Email:         "lucky1@kryptopos.com",
			Name:          "Lucky Fernanda R",
			EmailVerified: "Y",
			UpdatedAt:     time.Now(),
			CreatedAt:     time.Now(),
		},
	}
	repository.On("Fetch").Return(usersMock, nil)

	usecase := usecase.NewUserUsecase(repository, time.Second*5)
	users, err := usecase.Fetch(context.TODO())
	for i := range users {
		assert.Equal(t, users[i].Email, usersMock[i].Email, "user email not valid")
		assert.Equal(t, users[i].Username, usersMock[i].Username, "username not valid")
	}
	assert.NotEmpty(t, users)
	assert.NoError(t, err)
	assert.Len(t, usersMock, len(users))
	repository.AssertCalled(t, "Fetch")
}

func TestUserUsecase_Store(t *testing.T) {
	repository := new(mocks.UserRepository)
	usersMock := domain.User{
		Username: "LFR",
		Email:    "lucky@kryptopos.com",
		Name:     "Lucky Fernanda",
		Password: "123456",
	}

	t.Run("success", func(t *testing.T) {
		tempMockUser := usersMock
		tempMockUser.ID = 0
		repository.On("GetByEmail", mock.AnythingOfType("string")).Return(domain.User{}, nil).Once()
		repository.On("GetByUsername", mock.AnythingOfType("string")).Return(domain.User{}, nil).Once()
		repository.On("Store", mock.AnythingOfType("*domain.User")).Return(nil).Once()

		usecase := usecase.NewUserUsecase(repository, time.Second*5)
		err := usecase.Store(context.TODO(), &usersMock)

		assert.NoError(t, err)
		assert.Equal(t, usersMock.Username, tempMockUser.Username)
		repository.AssertExpectations(t)
	})

	t.Run("existing-email", func(t *testing.T) {
		existingUser := usersMock
		repository.On("GetByEmail", mock.AnythingOfType("string")).Return(existingUser, nil).Once()

		usecase := usecase.NewUserUsecase(repository, time.Second*5)
		err := usecase.Store(context.TODO(), &usersMock)

		assert.Error(t, err)
		repository.AssertExpectations(t)
	})

	t.Run("existing-username", func(t *testing.T) {
		existingUser := usersMock
		repository.On("GetByEmail", mock.AnythingOfType("string")).Return(domain.User{}, nil).Once()
		repository.On("GetByUsername", mock.AnythingOfType("string")).Return(existingUser, nil).Once()

		usecase := usecase.NewUserUsecase(repository, time.Second*5)
		err := usecase.Store(context.TODO(), &usersMock)

		assert.Error(t, err)
		repository.AssertExpectations(t)
	})

}

func TestUserUsecase_GetByID(t *testing.T) {
	repository := new(mocks.UserRepository)
	usersMock := domain.User{ID: 1,
		Username:      "LFR",
		Email:         "lucky@kryptopos.com",
		Name:          "Lucky Fernanda",
		EmailVerified: "N",
		UpdatedAt:     time.Now(),
		CreatedAt:     time.Now(),
	}
	t.Run("success", func(t *testing.T) {
		repository.On("GetByID", mock.AnythingOfType("int64")).Return(usersMock, nil).Once()

		usecase := usecase.NewUserUsecase(repository, time.Second*5)
		user, err := usecase.GetByID(context.TODO(), usersMock.ID)

		assert.NotNil(t, user)
		assert.NoError(t, err)
		repository.AssertExpectations(t)
	})

	t.Run("id-not-found", func(t *testing.T) {
		repository.On("GetByID", mock.AnythingOfType("int64")).Return(domain.User{}, nil).Once()

		usecase := usecase.NewUserUsecase(repository, time.Second*5)
		user, err := usecase.GetByID(context.TODO(), usersMock.ID)

		assert.Equal(t, domain.User{}, user)
		assert.NoError(t, err)
		repository.AssertExpectations(t)
	})
}

func TestUserUsecase_GetByUsername(t *testing.T) {
	repository := new(mocks.UserRepository)
	usersMock := domain.User{ID: 1,
		Username:      "LFR",
		Email:         "lucky@kryptopos.com",
		Name:          "Lucky Fernanda",
		EmailVerified: "N",
		UpdatedAt:     time.Now(),
		CreatedAt:     time.Now(),
	}
	t.Run("success", func(t *testing.T) {
		repository.On("GetByUsername", mock.AnythingOfType("string")).Return(usersMock, nil).Once()

		usecase := usecase.NewUserUsecase(repository, time.Second*5)
		user, err := usecase.GetByUsername(context.TODO(), usersMock.Email)

		assert.NotNil(t, user)
		assert.NoError(t, err)
		repository.AssertExpectations(t)
	})

	t.Run("username-not-found", func(t *testing.T) {
		repository.On("GetByUsername", mock.AnythingOfType("string")).Return(domain.User{}, nil).Once()

		usecase := usecase.NewUserUsecase(repository, time.Second*5)
		user, err := usecase.GetByUsername(context.TODO(), usersMock.Email)

		assert.Equal(t, domain.User{}, user)
		assert.NoError(t, err)
		repository.AssertExpectations(t)
	})
}

func TestUserUsecase_getByEmail(t *testing.T) {
	repository := new(mocks.UserRepository)
	usersMock := domain.User{
		ID:            1,
		Username:      "LFR",
		Email:         "lucky@kryptopos.com",
		Name:          "Lucky Fernanda",
		EmailVerified: "N",
		UpdatedAt:     time.Now(),
		CreatedAt:     time.Now(),
	}
	t.Run("success", func(t *testing.T) {
		repository.On("GetByEmail", mock.AnythingOfType("string")).Return(usersMock, nil).Once()

		usecase := usecase.NewUserUsecase(repository, time.Second*5)
		user, err := usecase.GetByEmail(context.TODO(), usersMock.Email)

		assert.NotNil(t, user)
		assert.NoError(t, err)
		repository.AssertExpectations(t)
	})

	t.Run("email-not-found", func(t *testing.T) {
		repository.On("GetByEmail", mock.AnythingOfType("string")).Return(domain.User{}, nil).Once()

		usecase := usecase.NewUserUsecase(repository, time.Second*5)
		user, err := usecase.GetByEmail(context.TODO(), usersMock.Email)

		assert.Equal(t, domain.User{}, user)
		assert.NoError(t, err)
		repository.AssertExpectations(t)
	})
}

func TestUserUsecase_Update(t *testing.T) {
	repository := new(mocks.UserRepository)
	usersMock := domain.User{
		Username: "LFR",
		Email:    "lucky@kryptopos.com",
		Name:     "Lucky Fernanda RRRR",
		Password: "123456",
	}

	t.Run("success", func(t *testing.T) {
		tempMockUser := usersMock
		tempMockUser.ID = 0

		repository.On("Update", mock.AnythingOfType("*domain.User")).Return(nil).Once()

		usecase := usecase.NewUserUsecase(repository, time.Second*5)
		err := usecase.Update(context.TODO(), &usersMock)

		assert.NoError(t, err)
		assert.Equal(t, usersMock.Name, tempMockUser.Name)
		repository.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		tempMockUser := usersMock
		tempMockUser.ID = 0

		repository.On("Update", mock.AnythingOfType("*domain.User")).Return(errors.New("record not found")).Once()

		usecase := usecase.NewUserUsecase(repository, time.Second*5)
		err := usecase.Update(context.TODO(), &usersMock)

		assert.Error(t, err)
		assert.Equal(t, usersMock.Name, tempMockUser.Name)
		repository.AssertExpectations(t)
	})

}

func TestUserUsecase_Delete(t *testing.T) {
	repository := new(mocks.UserRepository)
	usersMock := domain.User{
		ID:       1,
		Username: "LFR",
		Email:    "lucky@kryptopos.com",
		Name:     "Lucky Fernanda RRRR",
		Password: "123456",
	}

	t.Run("success", func(t *testing.T) {

		repository.On("Delete", mock.AnythingOfType("int64")).Return(nil).Once()
		repository.On("GetByID", mock.AnythingOfType("int64")).Return(usersMock, nil).Once()

		usecase := usecase.NewUserUsecase(repository, time.Second*5)
		err := usecase.Delete(context.TODO(), usersMock.ID)

		assert.NoError(t, err)
		repository.AssertExpectations(t)
	})

	t.Run("user-not-found", func(t *testing.T) {
		repository.On("GetByID", mock.AnythingOfType("int64")).Return(domain.User{}, errors.New("record not found")).Once()

		usecase := usecase.NewUserUsecase(repository, time.Second*5)
		err := usecase.Delete(context.TODO(), usersMock.ID)

		assert.Error(t, err)
		repository.AssertExpectations(t)
	})

	t.Run("something-wrong-db", func(t *testing.T) {

		repository.On("GetByID", mock.AnythingOfType("int64")).Return(usersMock, nil).Once()
		repository.On("Delete", mock.AnythingOfType("int64")).Return(errors.New("unexpected error")).Once()

		usecase := usecase.NewUserUsecase(repository, time.Second*5)
		err := usecase.Delete(context.TODO(), usersMock.ID)

		assert.Error(t, err)
		repository.AssertExpectations(t)
	})
}
