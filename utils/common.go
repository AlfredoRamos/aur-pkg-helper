package utils

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"
)

func RootPath() (string, error) {
	rootPath := filepath.Clean(os.Getenv("AUR_ROOT_PATH"))
	if len(rootPath) < 1 || rootPath == "." {
		return "", errors.New("please set the AUR root path: AUR_ROOT_PATH")
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
