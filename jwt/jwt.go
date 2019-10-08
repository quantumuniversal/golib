package jwt

import "github.com/dgrijalva/jwt-go"

// MapClaims --
func MapClaims(tokenString string, verificationKey string) (*jwt.Token, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(verificationKey), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}
