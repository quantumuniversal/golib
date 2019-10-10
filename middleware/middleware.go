package middleware

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

// Cors --
func Cors(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println("Running Cors")
		return next(c)
	}
}

// GetAccountNumber --
func GetAccountNumber(c echo.Context, pk []byte) string {
	var accountNumber string

	accountNumber = extractAccountNumber(c, pk)
	if accountNumber == "" {
		accountNumber = c.Request().Header.Get("accountNumber")
	}

	return accountNumber
}

// extractAccountNumber --
func extractAccountNumber(c echo.Context, pk []byte) string {
	bearer := c.Request().Header.Get("Authorization")
	token := strings.Split(bearer, " ")

	if len(token) > 1 {
		tokenString := token[1]
		parsedJWT, _ := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			return pk, nil
		})

		var claims map[string]interface{}

		bytes, _ := json.Marshal(parsedJWT.Claims)
		json.Unmarshal(bytes, &claims)

		return claims["accountNumber"].(string)
	}

	return ""
}
