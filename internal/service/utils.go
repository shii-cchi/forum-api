package service

import (
	"github.com/shii-cchi/forum-api/pkg/auth"
	"time"
)

func (s UserService) CreateTokens(userId string) (string, string, error) {
	m, err := auth.NewManager(s.cfg.AccessSigningKey)

	if err != nil {
		return "", "", err
	}

	TTL, err := time.ParseDuration(s.cfg.AccessTTL)

	accessToken, err := m.NewJWT(userId, TTL)

	if err != nil {
		return "", "", err
	}

	m, err = auth.NewManager(s.cfg.RefreshSigningKey)

	if err != nil {
		return "", "", err
	}

	TTL, err = time.ParseDuration(s.cfg.RefreshTTL)

	refreshToken, err := m.NewJWT(userId, TTL)

	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s UserService) IsValidToken(validatedToken string) (bool, error) {
	m, err := auth.NewManager(s.cfg.RefreshSigningKey)
	if err != nil {
		return false, err
	}

	ok, err := m.IsValidToken(validatedToken)

	return ok, err
}
