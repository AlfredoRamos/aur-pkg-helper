package utils

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"

	"alfredoramos.mx/aur-pkg-helper/config"
)

func RootPath() (string, error) {
	config := config.LoadConfig()
	rootPath := filepath.Clean(config.String("aur.root_path", ""))

	if len(rootPath) < 1 || rootPath == "." {
		return "", errors.New("please set the AUR root path: aur.root_path")
	}

	rootPath, err := filepath.Abs(rootPath)
	if err != nil {
		slog.Error("Could not get AUR root path", slog.Any("error", err))
		return "", err
	}

	if stat, err := os.Stat(rootPath); err != nil || !stat.IsDir() {
		if err != nil {
			slog.Error("Error reading AUR root path", slog.Any("error", err))
			return "", err
		}

		if !stat.IsDir() {
			slog.Error("The AUR root path is not a directory", slog.String("path", rootPath))
			return "", errors.New("the AUR root path is not a directory")
		}
	}

	return rootPath, nil
}
