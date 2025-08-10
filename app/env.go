package app

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func SetupEnvironment() {
	configPath, err := os.UserConfigDir()
	if err != nil {
		slog.Error("Could not get user configuration path", slog.Any("error", err))
		os.Exit(1)
	}

	envFile := filepath.Clean(filepath.Join(configPath, "aur-pkg-helper.env"))

	slog.Info("Reading environment", slog.String("file", envFile))

	if err := godotenv.Load(envFile); err != nil {
		slog.Error("Could not load .env file", slog.Any("error", err))
		os.Exit(1)
	}
}
