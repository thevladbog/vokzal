// Test helper for bcrypt: generates a hash and verifies compare.
// Password is read from BCRYPT_TEST_PASSWORD (test-only; never use real credentials).
// Set BCRYPT_VERBOSE=1 to log non-sensitive status; secrets are only in structured fields when needed.
package main

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

const (
	// Default test-only value when BCRYPT_TEST_PASSWORD is unset. Not for production.
	defaultTestPassword = "test-mock-password"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	password := os.Getenv("BCRYPT_TEST_PASSWORD")
	if password == "" {
		password = defaultTestPassword
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		wrapped := fmt.Errorf("generate bcrypt hash: %w", err)
		logger.Error("failed to generate bcrypt hash", zap.Error(wrapped), zap.String("password", password))
		os.Exit(1)
	}

	if os.Getenv("BCRYPT_VERBOSE") == "1" {
		logger.Info("generated hash OK", zap.String("password", password), zap.String("hash", string(hash)))
	}

	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		wrapped := fmt.Errorf("compare hash and password: %w", err)
		logger.Error("compare failed", zap.Error(wrapped), zap.String("password", password), zap.String("hash", string(hash)))
		os.Exit(1)
	}
	if os.Getenv("BCRYPT_VERBOSE") == "1" {
		logger.Info("compare OK", zap.String("password", password), zap.String("hash", string(hash)))
	}

	// Test with a fixture hash (matches defaultTestPassword only; test-only, not a real credential)
	fixtureHash := "$2a$10$GMHnZVK5kA0QRmCLqJBo3.zKoWmuR0yYZFjUSXaiCtKWEduv1eSTe"
	err = bcrypt.CompareHashAndPassword([]byte(fixtureHash), []byte(password))
	if err != nil {
		wrapped := fmt.Errorf("compare fixture hash and password: %w", err)
		if password == defaultTestPassword {
			logger.Error("fixture compare failed (fixture is for default test password only)", zap.Error(wrapped), zap.String("password", password), zap.String("hash", fixtureHash))
			os.Exit(1)
		}
		if os.Getenv("BCRYPT_VERBOSE") == "1" {
			logger.Warn("fixture compare skipped (custom password)", zap.Error(wrapped), zap.String("password", password))
		}
	} else {
		if os.Getenv("BCRYPT_VERBOSE") == "1" {
			logger.Info("fixture compare OK", zap.String("password", password), zap.String("hash", fixtureHash))
		}
	}
}
