package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Firebase FirebaseConfig
	Database DatabaseConfig
	Session  SessionConfig
}

type ServerConfig struct {
	Port        string
	FrontendURL string
}

type FirebaseConfig struct {
	CredentialsPath string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type SessionConfig struct {
	Secret string
	Name   string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	config := &Config{
		Server: ServerConfig{
			Port:        os.Getenv("PORT"),
			FrontendURL: os.Getenv("FRONTEND_URL"),
		},
		Firebase: FirebaseConfig{
			CredentialsPath: os.Getenv("FIREBASE_CREDENTIALS_PATH"),
		},
		Database: DatabaseConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName:   os.Getenv("DB_NAME"),
			SSLMode:  os.Getenv("DB_SSLMODE"),
		},
		Session: SessionConfig{
			Secret: os.Getenv("SESSION_SECRET"),
			Name:   os.Getenv("SESSION_NAME"),
		},
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) Validate() error {
	if c.Server.Port == "" {
		return fmt.Errorf("PORT is required")
	}
	if c.Server.FrontendURL == "" {
		return fmt.Errorf("FRONTEND_URL is required")
	}
	if c.Firebase.CredentialsPath == "" {
		return fmt.Errorf("FIREBASE_CREDENTIALS_PATH is required")
	}
	if c.Database.Host == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if c.Database.Port == "" {
		return fmt.Errorf("DB_PORT is required")
	}
	if c.Database.User == "" {
		return fmt.Errorf("DB_USER is required")
	}
	if c.Database.Password == "" {
		return fmt.Errorf("DB_PASSWORD is required")
	}
	if c.Database.DBName == "" {
		return fmt.Errorf("DB_NAME is required")
	}
	if c.Database.SSLMode == "" {
		return fmt.Errorf("DB_SSLMODE is required")
	}
	if c.Session.Secret == "" {
		return fmt.Errorf("SESSION_SECRET is required")
	}
	if c.Session.Name == "" {
		return fmt.Errorf("SESSION_NAME is required")
	}
	return nil
}

func (c *Config) GetDatabaseDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}
