package api

import (
	"encoding/json"
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
	e.POST("/refreshToken", handler.refreshToken)
	e.POST("/signup", handler.Signup)
	e.POST("/createVerifyEmail", handler.createVerifyEmail)
	e.POST("/verifyEmail", handler.verifyEmail)
	e.POST("/logout", handler.Logout)
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

func (handler *AuthHandler) refreshToken(c echo.Context) (err error) {
	jwtResults, err := handler.AuthUsecase.GenerateNewAccessToken(c)
	if err != nil {
		return handler.Response.Error(c, domain.ErrorAuthorization)
	}
	token := map[string]string{
		"access_token":  jwtResults.AccessToken,
		"refresh_token": jwtResults.RefreshToken,
	}

	return handler.Response.Success(c, "success", http.StatusOK, map[string]interface{}{"token": token})
}

func (handler *AuthHandler) createVerifyEmail(c echo.Context) (err error) {
	payload := make(map[string]interface{})
	err = json.NewDecoder(c.Request().Body).Decode(&payload)
	if err != nil {
		return err
	}
	email, ok := payload["email"].(string)
	if !ok {
		return domain.ErrBadParamInput
	}
	_, err = handler.AuthUsecase.CreateVerifyEmail(c.Request().Context(), email)
	return
}

func (handler *AuthHandler) verifyEmail(c echo.Context) (err error) {
	payload := make(map[string]interface{})
	err = json.NewDecoder(c.Request().Body).Decode(&payload)
	if err != nil {
		return domain.ErrBadParamInput
	}
	token, ok := payload["token"].(string)
	if !ok {
		return domain.ErrBadParamInput
	}
	err = handler.AuthUsecase.VerifyEmail(c.Request().Context(), token)
	if err != nil {
		return domain.ErrorTokenNotFound
	}
	return handler.Response.Success(c, "success", http.StatusOK, map[string]interface{}{})
}

func (handler *AuthHandler) Logout(c echo.Context) (err error) {
	payload := make(map[string]interface{})
	err = json.NewDecoder(c.Request().Body).Decode(&payload)
	if err != nil {
		return err
	}
	access_token, ok := payload["access_token"].(string)
	if !ok {
		return domain.ErrBadParamInput
	}

	refresh_token, ok := payload["refresh_token"].(string)
	if !ok {
		return domain.ErrBadParamInput
	}

	err = handler.AuthUsecase.Logout(access_token, refresh_token)
	if err != nil {
		return handler.Response.Error(c, domain.ErrorAuthorization)
	}
	return handler.Response.Success(c, "success", http.StatusOK, map[string]interface{}{})

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
