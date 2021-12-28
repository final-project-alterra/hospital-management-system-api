package middleware

import (
	"errors"
	"net/http"

	"github.com/final-project-alterra/hospital-management-system-api/config"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func IsAuth() echo.MiddlewareFunc {
	parseToken := func(auth string, c echo.Context) (interface{}, error) {
		keyFunction := func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(config.ENV.JWT_SECRET), nil
		}

		token, err := jwt.Parse(auth, keyFunction)
		if err != nil {
			return nil, err
		}
		if !token.Valid {
			return nil, errors.New("Invalid token")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return nil, errors.New("Invalid claims")
		}
		return claims, nil
	}

	errorHandlerWithContext := func(err error, c echo.Context) error {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": map[string]interface{}{
				"code":    http.StatusUnauthorized,
				"message": "Unauthorized user!",
			},
		})
	}

	jwtConfig := middleware.JWTConfig{
		ParseTokenFunc:          parseToken,
		ErrorHandlerWithContext: errorHandlerWithContext,
	}

	return middleware.JWTWithConfig(jwtConfig)
}

func IsAdmin() echo.MiddlewareFunc {
	parseToken := func(auth string, c echo.Context) (interface{}, error) {
		keyFunction := func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(config.ENV.JWT_SECRET), nil
		}

		token, err := jwt.Parse(auth, keyFunction)
		if err != nil {
			return nil, err
		}
		if !token.Valid {
			return nil, errors.New("Invalid token")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return nil, errors.New("Invalid claims")
		}

		if claims["role"] != "admin" {
			return nil, errors.New("Unauthorized user!")
		}

		return claims, nil
	}

	errorHandlerWithContext := func(err error, c echo.Context) error {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": map[string]interface{}{
				"code":    http.StatusUnauthorized,
				"message": "Unauthorized user!",
			},
		})
	}

	jwtConfig := middleware.JWTConfig{
		ParseTokenFunc:          parseToken,
		ErrorHandlerWithContext: errorHandlerWithContext,
	}

	return middleware.JWTWithConfig(jwtConfig)
}
