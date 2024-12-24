package utils

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
)

func SetupAurRepositories() error {
	rootPath, err := RootPath()
	if err != nil {
		slog.Error(fmt.Sprintf("Could not get AUR root path: %v", err))
		return err
	}

	repos, err := os.ReadDir(rootPath)
	if err != nil {
		slog.Error(fmt.Sprintf("Could not get content from AUR root path: %v", err))
		return err
	}

	for _, repo := range repos {
		repoPath := filepath.Clean(filepath.Join(rootPath, repo.Name()))

		// Ignore non-directories
		if !repo.IsDir() {
			continue
		}

		slog.Info(fmt.Sprintf("Processing package %s", repo.Name()))

		// Ignore non-git repositories
		if stat, err := os.Stat(filepath.Clean(filepath.Join(repoPath, ".git"))); err != nil || !stat.IsDir() {
			slog.Error(" -> Not a git repository")
			continue
		}

		if err := SetupGitConfig(repoPath); err != nil {
			slog.Error(fmt.Sprintf(" -> Git configuration [%t]", false))
		} else {
			slog.Info(fmt.Sprintf(" -> Git configuration [%t]", true))
		}

		if err := SetupGitHooks(repoPath); err != nil {
			slog.Error(fmt.Sprintf(" -> Git hooks [%t]", false))
		} else {
			slog.Info(fmt.Sprintf(" -> Git hooks [%t]", true))
		}
	}

	return nil
}

func copySourceHooks() error {
	srcHooksPath := "hooks"

	srcHooks, err := os.ReadDir(srcHooksPath)
	if err != nil {
		slog.Error(fmt.Sprintf("Could not get content from source hooks path: %v", err))
		return err
	}

	hooksPath, err := HooksPath()
	if err != nil {
		slog.Error(fmt.Sprintf("Could not get hooks path: %v", err))
		return err
	}

	errs := []error{}

	for _, hook := range srcHooks {
		// Ignore non-regular files
		if !hook.Type().IsRegular() {
			continue
		}

		// Ignore non-hook files
		if filepath.Ext(hook.Name()) != ".hook" {
			continue
		}

		srcHookFile := filepath.Clean(filepath.Join(srcHooksPath, hook.Name()))
		destHookFile := filepath.Clean(filepath.Join(hooksPath, hook.Name()))

		if stat, err := os.Stat(destHookFile); !os.IsNotExist(err) && (err != nil || !stat.Mode().IsRegular()) {
			if err != nil {
				slog.Error(fmt.Sprintf("Error reading repository path: %v", err))
				errs = append(errs, err)
			}

			if !stat.Mode().IsRegular() {
				slog.Error(fmt.Sprintf("The hook already exists and is not a valid file: %s", destHookFile))
				errs = append(errs, errors.New("Invalid hook file."))
			}
		}

		copyHook := exec.Command("cp", "-af", srcHookFile, destHookFile) //#nosec:G204
		if err := copyHook.Run(); err != nil {
			slog.Error(fmt.Sprintf("Could not execute command: %v", err))
			errs = append(errs, fmt.Errorf("Could not copy source hook '%s'.", hook.Name()))
		}
	}

	return errors.Join(errs...)
}

func HooksPath() (string, error) {
	hooksPath := filepath.Clean(os.Getenv("GIT_HOOKS_PATH"))
	if len(hooksPath) < 1 || hooksPath == "." {
		return "", errors.New("Please set the Git hooks path: GIT_HOOKS_PATH.")
	}

	hooksPath, err := filepath.Abs(hooksPath)
	if err != nil {
		slog.Error(fmt.Sprintf("Could not get Git hooks path: %v", err))
		return "", err
	}

	if stat, err := os.Stat(hooksPath); err != nil || !stat.IsDir() {
		if err != nil {
			slog.Error(fmt.Sprintf("Error reading Git hooks path: %v", err))
			return "", err
		}

		if !stat.IsDir() {
			slog.Error(fmt.Sprintf("The Git hooks path is not a directory: %s", hooksPath))
			return "", errors.New("The Git hooks path is not a directory.")
		}
	}

	return hooksPath, nil
}
