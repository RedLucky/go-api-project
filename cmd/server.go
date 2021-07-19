package main

import (
	_deliveryUser "go-api-project/app/delivery/api/user"
	_repoUser "go-api-project/app/repository/user"
	_ucUser "go-api-project/app/usecase/user"
	"go-api-project/config/db"

	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type M map[string]interface{}

// CORS will handle the CORS middleware
func cors(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")
		return next(c)
	}
}

func middlewareLogging(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		makeLogEntry(c).Info("incoming request")
		return next(c)
	}
}

func makeLogEntry(c echo.Context) *logrus.Entry {
	if c == nil {
		return logrus.WithFields(logrus.Fields{
			"at": time.Now().Format("2006-01-02 15:04:05"),
		})
	}

	return logrus.WithFields(logrus.Fields{
		"at":     time.Now().Format("2006-01-02 15:04:05"),
		"method": c.Request().Method,
		"uri":    c.Request().URL.String(),
		"ip":     c.Request().RemoteAddr,
	})
}

func init() {
	viper.SetConfigFile(`config/config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		log.Println("Service RUN on DEBUG mode")
	}
}

func main() {

	mysql := db.New().Conn

	defer func() {
		err := mysql.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	userRepo := _repoUser.NewUserRepository(mysql)
	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	userUc := _ucUser.NewUserUsecase(userRepo, timeoutContext)

	r := echo.New()
	// middL := middleware.InitMiddleware()
	r.Use(middlewareLogging)
	r.Use(cors)

	r.GET("/", func(c echo.Context) error {
		data := M{"Message": "Hello", "Counter": 2, "age": "27"}
		return c.JSON(http.StatusOK, data)
	})
	_deliveryUser.NewUserHandler(r, userUc)

	r.Logger.Fatal(r.Start(viper.GetString("server.address")))
}
