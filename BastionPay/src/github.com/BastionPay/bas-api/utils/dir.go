package utils

import (
	"os"
	"path/filepath"
	"strings"
	"os/user"
	"runtime"
)

// Get the app running directory
// NOTE: if you run like "go run main.go",
// this return is a temporary directory
func GetRunDir() (string, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", nil
	}
	return strings.Replace(dir, "\\", "/", -1), nil
}

// Get the current directory
// NOTE: this return value is the same as os.Getwd()
func GetCurrentDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return strings.Replace(dir, "\\", "/", -1), nil
}

// Get the app directory
// NOTE:
func GetAppDir() (string, error) {
	// Get the OS specific home directory via the Go standard lib.
	var homeDir string
	usr, err := user.Current()
	if err == nil {
		homeDir = usr.HomeDir
	}

	// Fall back to standard HOME environment variable that works
	// for most POSIX OSes if the directory from the Go standard
	// lib failed.
	if err != nil || homeDir == "" {
		homeDir = os.Getenv("HOME")
	}

	goos := runtime.GOOS
	switch goos {
	// Attempt to use the LOCALAPPDATA or APPDATA environment variable on
	// Windows.
	case "windows":
		// Windows XP and before didn't have a LOCALAPPDATA, so fallback
		// to regular APPDATA when LOCALAPPDATA is not set.
		appData := os.Getenv("LOCALAPPDATA")
		if appData == "" {
			appData = os.Getenv("APPDATA")
		}

		if appData != "" {
			return appData, nil
		}

	case "darwin":
		if homeDir != "" {
			return filepath.Join(homeDir, "Library",
				"Application Support"), nil
		}

	default:
		if homeDir != "" {
			return filepath.Join(homeDir, "."), nil
		}
	}

	// Fall back to the current directory if all else fails.
	return ".", nil
}

// Check path(directory or file) is exist
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	//if os.IsNotExist(err) {
	//	return false, nil
	//}
	return false, err
}