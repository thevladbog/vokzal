// Test helper for bcrypt: generates a hash and verifies compare.
// Password is read from BCRYPT_TEST_PASSWORD (test-only; never use real credentials).
// Set BCRYPT_VERBOSE=1 to log non-sensitive status (hash only; password is never logged).
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
	logger, err := zap.NewProduction()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to create logger: %v\n", err)
		os.Exit(1)
	}

	password := os.Getenv("BCRYPT_TEST_PASSWORD")
	if password == "" {
		password = defaultTestPassword
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("failed to generate bcrypt hash", zap.Error(fmt.Errorf("generate bcrypt hash: %w", err)))
		syncAndExit(logger, 1)
	}
	if os.Getenv("BCRYPT_VERBOSE") == "1" {
		logger.Info("generated hash OK", zap.String("hash", string(hash)))
	}

	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		logger.Error("compare failed", zap.Error(fmt.Errorf("compare hash and password: %w", err)), zap.String("hash", string(hash)))
		syncAndExit(logger, 1)
	}
	if os.Getenv("BCRYPT_VERBOSE") == "1" {
		logger.Info("compare OK", zap.String("hash", string(hash)))
	}

	// Intentionally public fixture hash for this test utility only: verifies defaultTestPassword.
	// Do not use for production or real credentials; only for run-time check that bcrypt compare works.
	fixtureHash := "$2a$10$GMHnZVK5kA0QRmCLqJBo3.zKoWmuR0yYZFjUSXaiCtKWEduv1eSTe"
	err = bcrypt.CompareHashAndPassword([]byte(fixtureHash), []byte(password))
	if err != nil && password == defaultTestPassword {
		logger.Error("fixture compare failed (fixture is for default test password only)",
			zap.Error(fmt.Errorf("compare fixture hash and password: %w", err)), zap.String("hash", fixtureHash))
		syncAndExit(logger, 1)
	}
	if err != nil && os.Getenv("BCRYPT_VERBOSE") == "1" {
		logger.Warn("fixture compare skipped (custom password)", zap.Error(err))
	}
	if err == nil && os.Getenv("BCRYPT_VERBOSE") == "1" {
		logger.Info("fixture compare OK", zap.String("hash", fixtureHash))
	}

	syncAndExit(logger, 0)
}

func syncAndExit(logger *zap.Logger, code int) {
	if syncErr := logger.Sync(); syncErr != nil {
		_, _ = fmt.Fprintf(os.Stderr, "logger sync: %v\n", syncErr)
	}
	os.Exit(code)
}
