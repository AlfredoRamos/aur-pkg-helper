package utils

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"alfredoramos.mx/aur-pkg-helper/config"
	"alfredoramos.mx/aur-pkg-helper/types"
)

func SetupGitConfig(repoPath string) error {
	config := config.LoadConfig()
	rootPath, err := RootPath()
	if err != nil {
		slog.Error("Could not get AUR root path", slog.Any("error", err))
		return err
	}

	repoPath = filepath.Clean(repoPath)

	// Avoid directory traversal attack
	if !strings.HasPrefix(repoPath, rootPath) {
		return errors.New("invalid repository path")
	}

	// Avoid directory traversal attack
	if rel, err := filepath.Rel(rootPath, repoPath); err != nil || strings.HasPrefix(rel, "..") || filepath.IsAbs(rel) {
		return errors.New("invalid repository file path: potential directory traversal")
	}

	if stat, err := os.Stat(repoPath); err != nil || !stat.IsDir() {
		if err != nil {
			slog.Error("Error reading repository path", slog.Any("error", err))
			return err
		}

		if !stat.IsDir() {
			slog.Error("The repository path is not a directory", slog.String("path", repoPath))
			return errors.New("invalid repository path")
		}
	}

	gitConfig := &types.GitConfig{
		Name:  strings.TrimSpace(config.String("git.user_name", "")),
		Email: strings.TrimSpace(config.String("git.user_email", "")),
	}

	errs := []error{}

	if !gitConfig.IsValidName() {
		errs = append(errs, errors.New("please set the git user name: git.user_name"))
	}

	if !gitConfig.IsValidEmail() {
		errs = append(errs, errors.New("please set the git user email: git.user_email"))
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	cmdName := exec.Command("git", "config", "--local", "--replace-all", "user.name", gitConfig.Name) //#nosec:G204
	cmdName.Dir = repoPath
	if err := cmdName.Run(); err != nil {
		slog.Error("Could not execute command", slog.Any("error", err))
		errs = append(errs, errors.New("could not set Git user name"))
	}

	cmdEmail := exec.Command("git", "config", "--local", "--replace-all", "user.email", gitConfig.Email) //#nosec:G204
	cmdEmail.Dir = repoPath
	if err := cmdEmail.Run(); err != nil {
		slog.Error("Could not execute command", slog.Any("error", err))
		errs = append(errs, errors.New("could not set Git user email"))
	}

	return errors.Join(errs...)
}

func SetupGitHooks(repoPath string) error {
	rootPath, err := RootPath()
	if err != nil {
		slog.Error("Could not get AUR root path", slog.Any("error", err))
		return err
	}

	repoPath = filepath.Clean(repoPath)

	// Avoid directory traversal attack
	if !strings.HasPrefix(repoPath, rootPath) {
		return errors.New("invalid repository path")
	}

	// Avoid directory traversal attack
	if rel, err := filepath.Rel(rootPath, repoPath); err != nil || strings.HasPrefix(rel, "..") || filepath.IsAbs(rel) {
		return errors.New("invalid repository file path: potential directory traversal")
	}

	if err := copySourceHooks(); err != nil {
		slog.Error("Could not copy source hooks", slog.Any("error", err))
	}

	hooksPath, err := HooksPath()
	if err != nil {
		slog.Error("Could not get hooks path", slog.Any("error", err))
		return err
	}

	hooksPath = filepath.Clean(hooksPath)

	hooks, err := os.ReadDir(hooksPath)
	if err != nil {
		slog.Error("Could not get content from Git hooks path", slog.Any("error", err))
		return err
	}

	errs := []error{}

	for _, hook := range hooks {
		// Ignore non-regular files
		if !hook.Type().IsRegular() {
			continue
		}

		// Ignore non-hook files
		if filepath.Ext(hook.Name()) != ".hook" {
			continue
		}

		hookFile := filepath.Clean(filepath.Join(hooksPath, hook.Name()))

		// Avoid directory traversal attack
		if !strings.HasPrefix(hookFile, hooksPath) {
			errs = append(errs, fmt.Errorf("invalid repository hook file path: %s", hook.Name()))
			continue
		}

		// Avoid directory traversal attack
		if rel, err := filepath.Rel(hooksPath, hookFile); err != nil || strings.HasPrefix(rel, "..") || filepath.IsAbs(rel) {
			errs = append(errs, fmt.Errorf("invalid repository hook file path: potential directory traversal: %s", hook.Name()))
			continue
		}

		repoHookFile := filepath.Clean(filepath.Join(repoPath, ".git", "hooks", strings.ReplaceAll(hook.Name(), filepath.Ext(hook.Name()), "")))

		// Avoid directory traversal attack
		if !strings.HasPrefix(repoHookFile, repoPath) {
			errs = append(errs, fmt.Errorf("invalid repository hook file path: %s", repoHookFile))
			continue
		}

		// Avoid directory traversal attack
		if rel, err := filepath.Rel(repoPath, repoHookFile); err != nil || strings.HasPrefix(rel, "..") || filepath.IsAbs(rel) {
			errs = append(errs, fmt.Errorf("invalid repository hook file path: potential directory traversal: %s", repoHookFile))
			continue
		}

		copyHook := exec.Command("cp", "-af", hookFile, repoHookFile) //#nosec:G204
		copyHook.Dir = repoPath
		if err := copyHook.Run(); err != nil {
			slog.Error("Could not execute command", slog.Any("error", err))
			errs = append(errs, fmt.Errorf("could not copy hook '%s'", hook.Name()))
		}
	}

	return errors.Join(errs...)
}
