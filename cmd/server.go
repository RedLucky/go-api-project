package main

import (
	_delivery "github.com/RedLucky/potongin/app/delivery/api"
	_customMiddleware "github.com/RedLucky/potongin/app/delivery/api/middleware"
	_AuthMiddleware "github.com/RedLucky/potongin/app/delivery/api/middleware/auth"
	"github.com/RedLucky/potongin/app/delivery/api/response"
	_repo "github.com/RedLucky/potongin/app/repository"
	_uc "github.com/RedLucky/potongin/app/usecase"
	"github.com/RedLucky/potongin/config/db"

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

	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	// user
	userRepo := _repo.NewUserRepository(mysql)
	userUc := _uc.NewUserUsecase(userRepo, timeoutContext)

	// auth
	authRepo := _repo.NewAuthRepository(mysql)
	authUc := _uc.NewAuthUsecase(authRepo, timeoutContext)

	// generated url
	generatedUrlRepo := _repo.NewGeneratedUrlRepository(mysql)
	generatedUrlUc := _uc.NewGeneratedUrlUsecase(generatedUrlRepo, timeoutContext)

	r := echo.New()
	middL := _customMiddleware.New()
	authMiddl := _AuthMiddleware.New()
	response := response.New()
	r.Use(echo.WrapMiddleware(middL.CorsMiddleware.Handler))
	r.Use(middL.MiddlewareLogging)

	r.GET("/stats", middL.Handle)
	_delivery.NewAuthHandler(r, authUc, response)
	_delivery.NewHitUrlHandler(r, generatedUrlUc, response)
	apiProtect := r.Group("")

	apiProtect.Use(authMiddl.Authentication)
	_delivery.NewUserHandler(apiProtect, userUc, response)
	_delivery.NewGeneratedUrlHandler(apiProtect, generatedUrlUc, response)

	// Configure middleware with the custom claims type
	// config := middleware.JWTConfig{
	// 	Claims:     &domain.JwtCustomClaims{},
	// 	SigningKey: []byte(viper.GetString(`authentication.jwt_signature_key`)),
	// }
	// r.Use(middleware.JWTWithConfig(config))

	r.Logger.Fatal(r.Start(viper.GetString("server.address")))
}
