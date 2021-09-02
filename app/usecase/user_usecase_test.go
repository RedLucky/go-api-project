package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/RedLucky/potongin/app/usecase"
	"github.com/RedLucky/potongin/domain"
	"github.com/RedLucky/potongin/domain/mocks"
	"github.com/stretchr/testify/assert"
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

func TestUserRepository_Store(t *testing.T) {

}
