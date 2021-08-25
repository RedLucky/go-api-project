package api

import (
	"net/http"

	"github.com/RedLucky/potongin/app/delivery/api/response"
	"github.com/RedLucky/potongin/domain"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	AuthUsecase domain.AuthUsecase
	Response    *response.JsonResponse
}

func NewAuthHandler(e *echo.Echo, uc domain.AuthUsecase, response *response.JsonResponse) {
	handler := &AuthHandler{
		AuthUsecase: uc,
		Response:    response,
	}
	e.POST("/login", handler.Login)
	e.POST("/signup", handler.Signup)
}

func (handler *AuthHandler) Signup(c echo.Context) (err error) {
	var user domain.User
	err = c.Bind(&user)
	if err != nil {
		return handler.Response.Error(c, err)
	}

	var ok bool
	if ok, err = isValidUser(&user); !ok {
		return handler.Response.Error(c, err)
	}

	ctx := c.Request().Context()
	err = handler.AuthUsecase.SignUp(ctx, &user)
	if err != nil {
		return handler.Response.Error(c, err)
	}
	return handler.Response.Success(c, "success", http.StatusCreated, map[string]interface{}{})

}

func (handler *AuthHandler) Login(c echo.Context) (err error) {
	var auth domain.Auth
	err = c.Bind(&auth)
	if err != nil {
		return handler.Response.Error(c, err)
	}

	var ok bool
	if ok, err = validateLogin(&auth); !ok {
		return handler.Response.Error(c, err)
	}

	ctx := c.Request().Context()
	jwtResults, err := handler.AuthUsecase.Authenticate(ctx, auth.Email, auth.Password)
	if err != nil {
		return handler.Response.Error(c, err)
	}
	token := map[string]string{
		"access_token":  jwtResults.AccessToken,
		"refresh_token": jwtResults.RefreshToken,
	}

	return handler.Response.Success(c, "success", http.StatusOK, map[string]interface{}{"token": token})

}

func validateLogin(m *domain.Auth) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

func isValidUser(m *domain.User) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}
