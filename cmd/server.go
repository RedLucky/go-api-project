package main

import (
	"go-api-project/app/delivery/api/middleware"
	_deliveryUser "go-api-project/app/delivery/api/user"
	_repoUser "go-api-project/app/repository/user"
	_ucUser "go-api-project/app/usecase/user"
	"go-api-project/config/db"

	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

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
	middL := middleware.New()
	r.Use(echo.WrapMiddleware(middL.CorsMiddleware.Handler))
	r.Use(middL.MiddlewareLogging)

	r.GET("/stats", middL.Handle)
	_deliveryUser.NewUserHandler(r, userUc)

	r.Logger.Fatal(r.Start(viper.GetString("server.address")))
}
