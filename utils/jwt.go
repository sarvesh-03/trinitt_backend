package utils

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/trinitt/config"
)

type JwtCustomClaims struct {
	ID uint `json:"id"`
	jwt.RegisteredClaims
}

var JWTConfig = echojwt.Config{
	NewClaimsFunc: func(c echo.Context) jwt.Claims {
		return new(JwtCustomClaims)
	},
	ErrorHandler: func(c echo.Context, err error) error {
		if err == echojwt.ErrJWTMissing {
			return c.JSON(http.StatusBadRequest, "token not found")
		}
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return c.JSON(http.StatusBadRequest, "token is malformed")
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				return c.JSON(http.StatusUnauthorized, "token is expired or not valid yet")
			} else {
				return c.JSON(http.StatusBadRequest, "token is invalid")
			}
		}
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, "something went wrong")
	},
	SigningKey: []byte("secret"),
}

func GenerateToken(userID uint) (string, error) {
	// Set custom claims
	claims := &JwtCustomClaims{
		ID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte("secret"))
}

func GetUserID(c echo.Context) uint {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JwtCustomClaims)
	return claims.ID
}

func GetUserIDFromToken(token string) (uint, error) {
	claims := &JwtCustomClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWTSecret), nil
	})
	if err != nil {
		return 0, err
	}
	return claims.ID, nil
}
