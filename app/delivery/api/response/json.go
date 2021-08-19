package response

import (
	"net/http"

	"github.com/RedLucky/potongin/domain"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type JsonResponse struct {
	Message string                 `json:"message"`
	Code    int                    `json:"code"`
	Data    map[string]interface{} `json:"data"`
}

func New() *JsonResponse {
	return &JsonResponse{}
}

func (response *JsonResponse) Success(ctx echo.Context, message string, status_code int, data map[string]interface{}) error {
	response.Message = message
	response.Code = status_code
	response.Data = data

	return ctx.JSON(response.Code, response)
}

func (response *JsonResponse) Error(ctx echo.Context, err error) error {
	response.Message = err.Error()
	response.Code = getStatusCode(err)
	response.Data = nil

	return ctx.JSON(response.Code, response)
}
func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	logrus.Error(err)
	switch err {
	case domain.ErrInternalServerError:
		return http.StatusInternalServerError
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrConflict:
		return http.StatusConflict
	case domain.ErrorAuthorization:
		return http.StatusUnauthorized
	case domain.ErrPassword:
		return http.StatusUnauthorized
	case domain.ErrEmailNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
