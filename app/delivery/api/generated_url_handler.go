package api

import (
	"net/http"

	"github.com/RedLucky/potongin/app/delivery/api/response"
	"github.com/RedLucky/potongin/domain"
	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type GeneratedUrlHandler struct {
	GeneratedUrlUsecase domain.GeneratedUrlUsecase
	Response            *response.JsonResponse
}

func NewGeneratedUrlHandler(e *echo.Group, guu domain.GeneratedUrlUsecase, response *response.JsonResponse) {
	handlers := &GeneratedUrlHandler{
		GeneratedUrlUsecase: guu,
		Response:            response,
	}

	e.POST("/createUrl", handlers.CreateUrl)
	e.GET("/urls", handlers.GetUrlByUserId)
	e.GET("/url/:url_id", handlers.GetUrlById)

}

func (handler *GeneratedUrlHandler) CreateUrl(c echo.Context) (err error) {

	var generateUrl domain.GeneratedUrl
	err = c.Bind(&generateUrl)
	if err != nil {
		return handler.Response.Error(c, err)
	}

	var ok bool
	if ok, err = validateCreateUrl(&generateUrl); !ok {
		return handler.Response.Error(c, err)
	}

	data := c.Get("user").(jwt.MapClaims)
	idP := data["id"].(float64)
	id := int64(idP)
	generateUrl.UserId = id
	ctx := c.Request().Context()
	err = handler.GeneratedUrlUsecase.CreateUrl(ctx, &generateUrl)
	if err != nil {
		return handler.Response.Error(c, err)
	}

	return handler.Response.Success(c, "success", http.StatusCreated, map[string]interface{}{"generated_url": generateUrl})

}

func (handler *GeneratedUrlHandler) GetUrlByUserId(c echo.Context) (err error) {
	var generateUrl []domain.GeneratedUrl

	data := c.Get("user").(jwt.MapClaims)
	idP := data["id"].(float64)
	id := int64(idP)

	ctx := c.Request().Context()
	generateUrl, err = handler.GeneratedUrlUsecase.GetUrlByUserId(ctx, id)
	if err != nil {
		return handler.Response.Error(c, err)
	}

	return handler.Response.Success(c, "success", http.StatusOK, map[string]interface{}{"generated_url": generateUrl})

}

func (handler *GeneratedUrlHandler) GetUrlById(c echo.Context) (err error) {
	var generateUrl domain.GeneratedUrl
	id := c.Param("url_id")
	ctx := c.Request().Context()
	generateUrl, err = handler.GeneratedUrlUsecase.GetUrlById(ctx, id)
	if err != nil {
		return handler.Response.Error(c, err)
	}

	return handler.Response.Success(c, "success", http.StatusOK, map[string]interface{}{"generated_url": generateUrl})

}

// private function
func validateCreateUrl(m *domain.GeneratedUrl) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}
