package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/RedLucky/potongin/domain"
	jwt "github.com/golang-jwt/jwt"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

var (
	AccessToken  string = "access"
	RefreshToken string = "refresh"
)

func CreateToken(user *domain.User) (jwtResults domain.JwtResults, err error) {
	// access_token
	jwtResults.AccessUUID = uuid.New().String()
	jwtResults.AccessExp = time.Now().Add(time.Duration(viper.GetInt32(`authentication.duration_access`)) * time.Minute).Unix()
	Accessclaims := domain.JwtCustomClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    viper.GetString(`server.application_name`),
			ExpiresAt: jwtResults.AccessExp,
		},
		AccessUUID: jwtResults.AccessUUID,
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, Accessclaims)
	jwtResults.AccessToken, err = accessToken.SignedString([]byte(viper.GetString(`authentication.jwt_signature_access_key`)))
	if err != nil {
		return domain.JwtResults{}, err
	}

	// refresh_token
	jwtResults.RefreshUUID = uuid.New().String()
	jwtResults.RefreshExp = time.Now().Add(time.Duration(viper.GetInt32(`authentication.duration_refresh`)) * time.Hour).Unix() //18 hours
	Refreshclaims := domain.JwtCustomClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    viper.GetString(`server.application_name`),
			ExpiresAt: jwtResults.RefreshExp,
		},
		RefreshUUID: jwtResults.RefreshUUID,
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, Refreshclaims)
	jwtResults.RefreshToken, err = refreshToken.SignedString([]byte(viper.GetString(`authentication.jwt_signature_refresh_key`)))
	if err != nil {
		return domain.JwtResults{}, err
	}

	return
}

func TokenValid(tokenString, tipe string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		var key string
		if tipe == AccessToken {
			key = viper.GetString(`authentication.jwt_signature_access_key`)
		} else if tipe == RefreshToken {
			key = viper.GetString(`authentication.jwt_signature_refresh_key`)
		}
		return []byte(key), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

func ExtractToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	if !strings.Contains(bearerToken, "Bearer") {
		return ""
	}

	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

// save token to redis
func SaveToken(redisConn redis.Conn, user domain.User, jwt domain.JwtResults) error {
	redisConn.Send("MULTI")
	redisConn.Send("HSET", jwt.AccessUUID, "id", user.ID)
	redisConn.Send("HSET", jwt.AccessUUID, "flag", "access_token")
	redisConn.Send("EXPIRE", jwt.AccessUUID, viper.GetInt32(`authentication.duration_access`)*60) //minutes
	redisConn.Send("HSET", jwt.RefreshUUID, "id", user.ID)
	redisConn.Send("HSET", jwt.RefreshUUID, "flag", "refresh_token")
	redisConn.Send("EXPIRE", jwt.RefreshUUID, viper.GetInt32(`authentication.duration_refresh`)*60*60) //hours
	_, err := redisConn.Do("EXEC")
	if err != nil {
		return err
	}
	return nil
}

func GetTokenFromRedis(redisConn redis.Conn, uuid string) (userId int64, err error) {
	userId, err = redis.Int64(redisConn.Do("HGET", uuid, "id"))
	return
}

func DeleteTokenRedis(redisConn redis.Conn, uuid string) (err error) {
	_, err = redisConn.Do("DEL", uuid)
	return
}
