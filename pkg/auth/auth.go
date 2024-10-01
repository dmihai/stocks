package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Username string `json:"user"`
	jwt.RegisteredClaims
}

type Auth struct {
	Accounts map[string]string
	jwtKey   []byte
}

func NewAuth() *Auth {
	return &Auth{
		Accounts: map[string]string{
			"admin": os.Getenv("ADMIN_PASSWORD"),
		},
		jwtKey: []byte(os.Getenv("JWT_SECRET_KEY")),
	}
}

func (a *Auth) GenerateJWT(user string) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)

	claims := &Claims{
		Username: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(a.jwtKey)
}

func (a *Auth) ParseJWT(reqToken string) (string, error) {
	claims := &Claims{}
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return a.jwtKey, nil
	}

	token, err := jwt.ParseWithClaims(reqToken, claims, keyFunc)
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	return claims.Username, nil
}
