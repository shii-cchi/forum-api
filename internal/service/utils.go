package service

import (
	"fmt"
	"github.com/shii-cchi/forum-api/pkg/auth"
	"os"
	"time"
)

func CreateToken(key1 string, key2 string, userId string) (string, error) {
	signingKey := os.Getenv(key1)
	if signingKey == "" {
		return "", fmt.Errorf("%s is not found in the environment", key1)
	}

	TTLStr := os.Getenv(key2)
	if TTLStr == "" {
		return "", fmt.Errorf("%s is not found in the environment", key2)
	}

	TTL, err := time.ParseDuration(TTLStr)

	if err != nil {
		return "", err
	}

	m, err := auth.NewManager(signingKey)

	if err != nil {
		return "", err
	}

	token, err := m.NewJWT(userId, TTL)

	if err != nil {
		return "", err
	}

	return token, nil
}
