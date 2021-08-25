package auth

import (
	"net/http"
	"time"

	"github.com/RedLucky/potongin/app/delivery/api/auth"
	"github.com/RedLucky/potongin/domain"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

type AuthMiddleware struct {
}

func (m *AuthMiddleware) Authentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims, err := auth.TokenValid(c.Request())
		if err != nil {
			makeLogEntry(c).Error(domain.ErrorAuthorization)
			// this should be to check refresh token if exist (using redis)
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "Unathorized"})
		}
		c.Set("user", claims)
		return next(c)
	}
}

func makeLogEntry(c echo.Context) *log.Entry {
	if c == nil {
		return log.WithFields(log.Fields{
			"at": time.Now().Format("2006-01-02 15:04:05"),
		})
	}

	return log.WithFields(log.Fields{
		"at":     time.Now().Format("2006-01-02 15:04:05"),
		"method": c.Request().Method,
		"uri":    c.Request().URL.String(),
		"ip":     c.Request().RemoteAddr,
	})
}

// InitMiddleware initialize the middleware
func New() *AuthMiddleware {
	return &AuthMiddleware{}
}
