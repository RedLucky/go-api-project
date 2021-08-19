package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/RedLucky/potongin/domain"
	jwt "github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

func CreateToken(user *domain.User) (string, error) {
	claims := domain.JwtCustomClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    viper.GetString(`server.application_name`),
			ExpiresAt: time.Now().Add(time.Duration(viper.GetInt32(`authentication.duration`)) * time.Hour).Unix(),
		},
		ID: user.ID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(viper.GetString(`authentication.jwt_signature_key`)))

}

func TokenValid(r *http.Request) (jwt.MapClaims, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(viper.GetString(`authentication.jwt_signature_key`)), nil
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
