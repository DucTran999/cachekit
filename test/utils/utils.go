package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func LoadEnv(envFile string) error {
	// Locate the directory containing go.mod
	goModDir, err := findGoModDir()
	if err != nil {
		return fmt.Errorf("go.mod not found: %w", err)
	}

	// Build the full path to the config file (e.g., .env)
	path := filepath.Join(goModDir, envFile)

	// Load env
	return godotenv.Load(path)
}

func findGoModDir() (string, error) {
	startPath, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}

	dir := startPath
	if fi, err := os.Stat(startPath); err == nil && !fi.IsDir() {
		dir = filepath.Dir(startPath)
	}

	for {
		goMod := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goMod); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break // reached root
		}
		dir = parent
	}

	return "", fmt.Errorf("go.mod not found from path: %s", startPath)
}
