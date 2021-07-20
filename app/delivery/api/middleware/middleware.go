package middleware

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

// GoMiddleware represent the data-struct for middleware
type CustomMiddleware struct {
	// another stuff , may be needed by middleware
	Uptime         string     `json:"uptime"`
	CorsMiddleware *cors.Cors `json:"-"`
}

// CORS will handle the CORS middleware
// func (m *CustomMiddleware) CORS(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		c.Response().Header().Set("Access-Control-Allow-Origin", "*")
// 		return next(c)
// 	}
// }

func (m *CustomMiddleware) MiddlewareLogging(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		makeLogEntry(c).Info("incoming request")
		return next(c)
	}
}

// Handle is the endpoint to get stats.
func (s *CustomMiddleware) Handle(c echo.Context) error {
	return c.JSON(http.StatusOK, s)
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
func New() *CustomMiddleware {
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"OPTIONS", "GET", "POST", "PUT"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
		Debug:          true,
	})
	return &CustomMiddleware{
		Uptime:         time.Now().Format("2006-01-02 15:04:05"),
		CorsMiddleware: corsMiddleware,
	}
}
