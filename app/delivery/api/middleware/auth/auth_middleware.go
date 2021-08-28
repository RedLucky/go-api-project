package auth

import (
	"net/http"
	"time"

	"github.com/RedLucky/potongin/app/delivery/api/auth"
	"github.com/RedLucky/potongin/domain"
	"github.com/gomodule/redigo/redis"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

type AuthMiddleware struct {
	RedisPool *redis.Pool
}

func (m *AuthMiddleware) Authentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		jwt := auth.ExtractToken(c.Request())
		claims, err := auth.TokenValid(jwt, auth.AccessToken)
		if err != nil {
			makeLogEntry(c).Error(domain.ErrorAuthorization)
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "Unathorized"})
		}
		// check to redis
		conn := m.RedisPool.Get()
		defer conn.Close()
		userId, err := auth.GetTokenFromRedis(conn, claims["access_uuid"].(string))
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "Unathorized"})
		}
		c.Set("user_id", userId)
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
func New(redisPool *redis.Pool) *AuthMiddleware {
	return &AuthMiddleware{
		RedisPool: redisPool,
	}
}
