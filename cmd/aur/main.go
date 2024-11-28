package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"alfredoramos.mx/aur-pkg-helper/app"
	"alfredoramos.mx/aur-pkg-helper/utils"
)

func showVersionInfo() {
	execPath, err := os.Executable()
	if err != nil {
		slog.Error(fmt.Sprintf("Could not get executable name: %v", err))
		os.Exit(1)
	}

	slog.Info(strings.TrimSpace(fmt.Sprintf("%s %s", filepath.Base(execPath), app.Version())))
}

func main() {
	showVersionInfo()

	app.SetupEnvironment()

	if err := utils.SetupAurRepositories(); err != nil {
		slog.Error(fmt.Sprintf("Could not setup AUR repositories: %v", err))
		os.Exit(1)
	}
}
