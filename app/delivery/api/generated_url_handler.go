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

type RequestParam struct {
	urlGenerated string
}

func NewGeneratedUrlHandler(e *echo.Group, guu domain.GeneratedUrlUsecase, response *response.JsonResponse) {
	handlers := &GeneratedUrlHandler{
		GeneratedUrlUsecase: guu,
		Response:            response,
	}

	e.POST("/createUrl", handlers.CreateUrl)

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

func (handler *GeneratedUrlHandler) HitUrl(c echo.Context) (err error) {
	ctx := c.Request().Context()
	var param = RequestParam{}
	if err = c.Bind(param); err != nil {
		return
	}
	results, err := handler.GeneratedUrlUsecase.HitUrl(ctx, param.urlGenerated)
	if err != nil {
		return handler.Response.Error(c, err)
	}

	return handler.Response.Success(c, "success", http.StatusOK, map[string]interface{}{"origin_url": results})

}

func validateCreateUrl(m *domain.GeneratedUrl) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}
