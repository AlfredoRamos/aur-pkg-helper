package app

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func SetupEnvironment() {
	configPath, err := os.UserConfigDir()
	if err != nil {
		slog.Error(fmt.Sprintf("Could not get user configuration path: %v", err))
		os.Exit(1)
	}

	envFile := filepath.Clean(filepath.Join(configPath, "aur-pkg-helper.env"))

	slog.Info(fmt.Sprintf("Reading environment file: %s", envFile))

	if err := godotenv.Load(envFile); err != nil {
		slog.Error(fmt.Sprintf("Could not load .env file: %v", err))
		os.Exit(1)
	}
}
