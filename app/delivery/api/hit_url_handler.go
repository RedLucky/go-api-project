package api

import (
	"net/http"

	"github.com/RedLucky/potongin/app/delivery/api/response"
	"github.com/RedLucky/potongin/domain"
	"github.com/labstack/echo/v4"
)

type HitUrlHandler struct {
	GeneratedUrlUsecase domain.GeneratedUrlUsecase
	Response            *response.JsonResponse
}

type RequestParam struct {
	UrlGenerated string `json:"url_generated"`
}

func NewHitUrlHandler(e *echo.Echo, guu domain.GeneratedUrlUsecase, response *response.JsonResponse) {
	handlers := &GeneratedUrlHandler{
		GeneratedUrlUsecase: guu,
		Response:            response,
	}

	e.POST("/accessUrl", handlers.HitUrl)

}

func (handler *GeneratedUrlHandler) HitUrl(c echo.Context) (err error) {
	ctx := c.Request().Context()
	var param RequestParam
	if err = c.Bind(&param); err != nil {
		return
	}
	results, err := handler.GeneratedUrlUsecase.HitUrl(ctx, param.UrlGenerated)
	if err != nil {
		return handler.Response.Error(c, err)
	}

	return handler.Response.Success(c, "success", http.StatusOK, map[string]interface{}{"origin_url": results})

}
