package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/RedLucky/potongin/domain"
	jwt "github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

func CreateToken(user *domain.User) (jwtResults domain.JwtResults, err error) {
	// access_token
	jwtResults.AccessUUID = uuid.New().String()
	jwtResults.AccessExp = time.Now().Add(time.Duration(viper.GetInt32(`authentication.duration`)) * time.Hour).Unix()
	Accessclaims := domain.JwtCustomClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    viper.GetString(`server.application_name`),
			ExpiresAt: jwtResults.AccessExp,
		},
		ID: user.ID,
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, Accessclaims)
	jwtResults.AccessToken, err = accessToken.SignedString([]byte(viper.GetString(`authentication.jwt_signature_access_key`)))
	if err != nil {
		return domain.JwtResults{}, err
	}

	// refresh_token
	jwtResults.RefreshUUID = uuid.New().String()
	jwtResults.RefreshExp = time.Now().Add(time.Duration(viper.GetInt32(`authentication.duration`)) * time.Hour * 24 * 7).Unix()
	Refreshclaims := domain.JwtCustomClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    viper.GetString(`server.application_name`),
			ExpiresAt: jwtResults.RefreshExp,
		},
		ID: user.ID,
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, Refreshclaims)
	jwtResults.RefreshToken, err = refreshToken.SignedString([]byte(viper.GetString(`authentication.jwt_signature_refresh_key`)))
	if err != nil {
		return domain.JwtResults{}, err
	}

	return

}

func TokenValid(r *http.Request) (jwt.MapClaims, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(viper.GetString(`authentication.jwt_signature_access_key`)), nil
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
