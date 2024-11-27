package utils

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

func RootPath() (string, error) {
	rootPath := filepath.Clean(os.Getenv("AUR_ROOT_PATH"))
	if len(rootPath) < 1 || rootPath == "." {
		return "", errors.New("Please set the AUR root path: AUR_ROOT_PATH.")
	}

	rootPath, err := filepath.Abs(rootPath)
	if err != nil {
		slog.Error(fmt.Sprintf("Could not get AUR root path: %v", err))
		return "", err
	}

	if stat, err := os.Stat(rootPath); err != nil || !stat.IsDir() {
		if err != nil {
			slog.Error(fmt.Sprintf("Error reading AUR root path: %v", err))
			return "", err
		}

		if !stat.IsDir() {
			slog.Error(fmt.Sprintf("The AUR root path is not a directory: %s", rootPath))
			return "", errors.New("The AUR root path is not a directory.")
		}
	}

	return rootPath, nil
}
