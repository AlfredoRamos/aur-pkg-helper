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
		slog.Error("Could not get AUR root path", slog.Any("error", err))
		return err
	}

	repos, err := os.ReadDir(rootPath)
	if err != nil {
		slog.Error(
			"Could not get content from AUR root path",
			slog.String("path", rootPath),
			slog.Any("error", err),
		)
		return err
	}

	for _, repo := range repos {
		repoPath := filepath.Clean(filepath.Join(rootPath, repo.Name()))

		// Ignore non-directories
		if !repo.IsDir() {
			continue
		}

		slog.Info("Processing", slog.String("repository", repo.Name()))

		// Ignore non-git repositories
		if stat, err := os.Stat(filepath.Clean(filepath.Join(repoPath, ".git"))); err != nil || !stat.IsDir() {
			slog.Error(" -> Not a git repository") //nolint:sloglint
			continue
		}

		if err := SetupGitConfig(repoPath); err != nil {
			slog.Error(fmt.Sprintf(" -> Git configuration [%t]", false)) //nolint:sloglint
		} else {
			slog.Info(fmt.Sprintf(" -> Git configuration [%t]", true)) //nolint:sloglint
		}

		if err := SetupGitHooks(repoPath); err != nil {
			slog.Error(fmt.Sprintf(" -> Git hooks [%t]", false)) //nolint:sloglint
		} else {
			slog.Info(fmt.Sprintf(" -> Git hooks [%t]", true)) //nolint:sloglint
		}
	}

	return nil
}

func copySourceHooks() error {
	srcHooksPath := "hooks"

	srcHooks, err := os.ReadDir(srcHooksPath)
	if err != nil {
		slog.Error("Could not get content from source hooks path", slog.Any("error", err))
		return err
	}

	hooksPath, err := HooksPath()
	if err != nil {
		slog.Error("Could not get hooks path", slog.Any("error", err))
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
				slog.Error("Error reading repository path", slog.Any("error", err))
				errs = append(errs, err)
			}

			if !stat.Mode().IsRegular() {
				slog.Error("The hook already exists and is not valid hook", slog.String("file", destHookFile))
				errs = append(errs, errors.New("invalid hook file"))
			}
		}

		copyHook := exec.Command("cp", "-af", srcHookFile, destHookFile) //#nosec:G204
		if err := copyHook.Run(); err != nil {
			slog.Error("Could not execute command", slog.Any("error", err))
			errs = append(errs, fmt.Errorf("could not copy source hook '%s'", hook.Name()))
		}
	}

	return errors.Join(errs...)
}

func HooksPath() (string, error) {
	hooksPath := filepath.Clean(os.Getenv("GIT_HOOKS_PATH"))
	if len(hooksPath) < 1 || hooksPath == "." {
		return "", errors.New("please set the Git hooks path: GIT_HOOKS_PATH")
	}

	hooksPath, err := filepath.Abs(hooksPath)
	if err != nil {
		slog.Error("Could not get Git hooks path", slog.Any("error", err))
		return "", err
	}

	if stat, err := os.Stat(hooksPath); err != nil || !stat.IsDir() {
		if err != nil {
			slog.Error("Error reading Git hooks path", slog.Any("error", err))
			return "", err
		}

		if !stat.IsDir() {
			slog.Error("The Git hooks path is not a directory", slog.String("path", hooksPath))
			return "", errors.New("the Git hooks path is not a directory")
		}
	}

	return hooksPath, nil
}
