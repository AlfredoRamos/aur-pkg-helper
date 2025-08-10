package utils

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"alfredoramos.mx/aur-pkg-helper/types"
)

func SetupGitConfig(repoPath string) error {
	repoPath = filepath.Clean(repoPath)
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
		Name:  strings.TrimSpace(os.Getenv("GIT_USER_NAME")),
		Email: strings.TrimSpace(os.Getenv("GIT_USER_EMAIL")),
	}

	errs := []error{}

	if !gitConfig.IsValidName() {
		errs = append(errs, errors.New("please set the git user name: GIT_USER_NAME"))
	}

	if !gitConfig.IsValidEmail() {
		errs = append(errs, errors.New("please set the git user email: GIT_USER_EMAIL"))
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
	if err := copySourceHooks(); err != nil {
		slog.Error("Could not copy source hooks", slog.Any("error", err))
	}

	hooksPath, err := HooksPath()
	if err != nil {
		slog.Error("Could not get hooks path", slog.Any("error", err))
		return err
	}

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
		repoHookFile := filepath.Clean(filepath.Join(repoPath, ".git", "hooks", strings.ReplaceAll(hook.Name(), filepath.Ext(hook.Name()), "")))

		copyHook := exec.Command("cp", "-af", hookFile, repoHookFile) //#nosec:G204
		copyHook.Dir = repoPath
		if err := copyHook.Run(); err != nil {
			slog.Error("Could not execute command", slog.Any("error", err))
			errs = append(errs, fmt.Errorf("could not copy hook '%s'", hook.Name()))
		}
	}

	return errors.Join(errs...)
}
