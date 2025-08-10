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
		slog.Error("Could not get executable name", slog.Any("error", err))
		os.Exit(1)
	}

	fmt.Println(strings.TrimSpace(fmt.Sprintf("%s v%s / GPL-3.0-or-later", filepath.Base(execPath), app.Version()))) //nolint:forbidigo
	fmt.Println("© Alfredo Ramos <alfredo.ramos@proton.me> (https://alfredoramos.mx)")                               //nolint:forbidigo
	fmt.Println()                                                                                                    //nolint:forbidigo
}

func main() {
	showVersionInfo()

	app.SetupEnvironment()

	if err := utils.SetupAurRepositories(); err != nil {
		slog.Error("Could not setup AUR repositories", slog.Any("error", err))
		os.Exit(1)
	}
}
