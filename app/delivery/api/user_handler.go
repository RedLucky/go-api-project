package api

import (
	"net/http"
	"strconv"

	"github.com/RedLucky/potongin/app/delivery/api/response"
	"github.com/RedLucky/potongin/domain"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

// UserHandler  represent the httphandler for article
type UserHandler struct {
	UserUsecase domain.UserUsecase
	Response    *response.JsonResponse
}

var userResponse map[string]interface{}

// NewUserHandler will initialize the articles/ resources endpoint
func NewUserHandler(e *echo.Group, uc domain.UserUsecase, response *response.JsonResponse) {
	handler := &UserHandler{
		UserUsecase: uc,
		Response:    response,
	}
	e.GET("/users", handler.FetchUser)
	e.GET("/me", handler.Me)
	e.POST("/user", handler.Store)
	e.GET("/user/:id", handler.GetByID)
	e.DELETE("/user/:id", handler.Delete)
}

func (handler *UserHandler) Me(c echo.Context) error {
	userId := c.Get("user_id").(int64)
	ctx := c.Request().Context()

	user, err := handler.UserUsecase.GetByID(ctx, userId)
	if err != nil {
		return handler.Response.Error(c, err)
	}
	userResponse = map[string]interface{}{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"name":       user.Name,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	}
	return handler.Response.Success(c, "success", http.StatusOK, map[string]interface{}{"user": userResponse})

}

func (handler *UserHandler) FetchUser(c echo.Context) error {
	ctx := c.Request().Context()

	listUsr, err := handler.UserUsecase.Fetch(ctx)
	if err != nil {
		return handler.Response.Error(c, err)
	}
	return handler.Response.Success(c, "success", http.StatusOK, map[string]interface{}{"users": listUsr})
}

// GetByID will get user by given id
func (handler *UserHandler) GetByID(c echo.Context) error {
	idP, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return handler.Response.Error(c, err)
	}

	id := int64(idP)
	ctx := c.Request().Context()

	user, err := handler.UserUsecase.GetByID(ctx, id)
	if err != nil {
		return handler.Response.Error(c, err)
	}

	userResponse := map[string]interface{}{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"name":       user.Name,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	}
	return handler.Response.Success(c, "success", http.StatusOK, map[string]interface{}{"user": userResponse})

}

func isRequestValid(m *domain.User) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Store will store the user by given request body
func (handler *UserHandler) Store(c echo.Context) (err error) {
	var user domain.User
	err = c.Bind(&user)
	if err != nil {
		return handler.Response.Error(c, err)
	}

	var ok bool
	if ok, err = isRequestValid(&user); !ok {
		return handler.Response.Error(c, err)
	}

	ctx := c.Request().Context()
	err = handler.UserUsecase.Store(ctx, &user)
	if err != nil {
		return handler.Response.Error(c, err)
	}

	return handler.Response.Success(c, "success", http.StatusCreated, map[string]interface{}{"user": user})

}

// Delete will delete user by given param
func (handler *UserHandler) Delete(c echo.Context) error {
	idP, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return handler.Response.Error(c, err)
	}

	id := int64(idP)
	ctx := c.Request().Context()

	err = handler.UserUsecase.Delete(ctx, id)
	if err != nil {
		return handler.Response.Error(c, err)
	}
	return handler.Response.Success(c, "success", http.StatusCreated, map[string]interface{}{})
}
