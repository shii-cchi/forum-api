package auth

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type TokenManager interface {
	NewJWT(userId string, ttl time.Duration) (string, error)
	Parse(Token string) (string, error)
	IsValidToken(refreshToken string) (string, error)
}

type Manager struct {
	signingKey string
}

func NewManager(signingKey string) (*Manager, error) {
	if signingKey == "" {
		return nil, errors.New("Empty signing key")
	}

	return &Manager{
		signingKey: signingKey,
	}, nil
}

func (m *Manager) NewJWT(userId string, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(ttl).Unix(),
		Subject:   userId,
	})

	return token.SignedString([]byte(m.signingKey))
}

func (m *Manager) Parse(accessToken string) (string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(m.signingKey), nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return "", fmt.Errorf("Error get user claims from token")
	}

	return claims["sub"].(string), nil
}

func (m *Manager) IsValidToken(validatedToken string) (bool, error) {
	token, err := jwt.Parse(validatedToken, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(m.signingKey), nil
	})

	if err != nil {
		return false, err
	}

	if !token.Valid {
		return false, fmt.Errorf("Token is not valid")
	}

	return true, nil
}
