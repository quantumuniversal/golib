package token

import (
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

// Token --
type Token struct {
	Context echo.Context
	PubKey  []byte
	Claims  jwt.Claims
}

// Decode --
func (t *Token) Decode() error {
	bearerToken := getBearerTokenFromContext(t.Context)
	err := t.getClaims(bearerToken)

	if err != nil {
		return err
	}
}

func (t *Token) getClaims(bearerToken string) error {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(bearerToken, claims, func(token *jwt.Token) (interface{}, error) {
		return t.PubKey, nil
	})

	if err != nil {
		return err
	}

	t.Claims = token.Claims
	return nil
}

func getBearerTokenFromContext(c echo.Context) string {
	bearer := c.Request().Header.Get("Authorization")
	token := strings.Split(bearer, " ")
	tokenString := token[1]

	return tokenString
}
