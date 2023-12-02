package config

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	Port              string
	DbURI             string
	AccessTTL         string
	RefreshTTL        string
	AccessSigningKey  string
	RefreshSigningKey string
	SaltString        string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load(".env")

	if err != nil {
		return nil, err
	}

	port := os.Getenv("PORT")

	if port == "" {
		return nil, errors.New("PORT is not found")
	}

	dbURI := os.Getenv("DB_URI")

	if dbURI == "" {
		return nil, errors.New("DB_URL is not found")
	}

	accessTTL := os.Getenv("ACCESS_TTL")

	if accessTTL == "" {
		return nil, errors.New("ACCESS_TTL is not found")
	}

	refreshTTL := os.Getenv("REFRESH_TTL")

	if refreshTTL == "" {
		return nil, errors.New("REFRESH_TTL is not found")
	}

	accessSigningKey := os.Getenv("ACCESS_SIGNING_KEY")

	if accessSigningKey == "" {
		return nil, errors.New("ACCESS_SIGNING_KEY is not found")
	}

	refreshSigningKey := os.Getenv("REFRESH_SIGNING_KEY")

	if refreshSigningKey == "" {
		return nil, errors.New("REFRESH_SIGNING_KEY is not found")
	}

	saltString := os.Getenv("SALT_STRING")

	if saltString == "" {
		return nil, errors.New("SALT_STRING is not found")
	}

	return &Config{
		Port:              port,
		DbURI:             dbURI,
		AccessTTL:         accessTTL,
		RefreshTTL:        refreshTTL,
		AccessSigningKey:  accessSigningKey,
		RefreshSigningKey: refreshSigningKey,
		SaltString:        saltString,
	}, nil
}
