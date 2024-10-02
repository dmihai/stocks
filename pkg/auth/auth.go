package auth

import (
	"crypto/rand"
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
	Accounts      map[string]string
	refreshTokens map[string]string
	jwtKey        []byte
}

func NewAuth() *Auth {
	return &Auth{
		Accounts: map[string]string{
			"admin": os.Getenv("ADMIN_PASSWORD"),
		},
		refreshTokens: make(map[string]string),
		jwtKey:        []byte(os.Getenv("JWT_SECRET_KEY")),
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

func (a *Auth) GenerateRefreshToken(user string) string {
	b := make([]byte, 32)
	rand.Read(b)

	token := fmt.Sprintf("%x", b)
	a.refreshTokens[user] = token

	return token
}

func (a *Auth) FindUserByRefreshToken(refreshToken string) *string {
	for user, token := range a.refreshTokens {
		if token == refreshToken {
			return &user
		}
	}

	return nil
}
