package utils

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	TokenID string `json:"id"`
	UserID  string `json:"subject"`
	jwt.RegisteredClaims
}

func ParseToken(token string, secret []byte) (*Claims, error) {
	parsedToken, err := jwt.ParseWithClaims(
		token, &Claims{},
		func(t *jwt.Token) (any, error) { return secret, nil },
		jwt.WithValidMethods([]string{"HS256"}),
	)

	if err != nil {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(*Claims)
	if !ok || !parsedToken.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
