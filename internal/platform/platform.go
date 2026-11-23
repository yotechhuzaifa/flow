// Package platform holds the OS-specific bits. Keep it tiny.
package platform

import (
	"os"
	"path/filepath"
	"runtime"
)

// DataDir returns the per-user application data directory, creating it.
// macOS: ~/Library/Application Support/<app>
// Linux: ~/.config/<app>   (XDG_CONFIG_HOME honoured)
// Windows: %AppData%\<app>
func DataDir(app string) (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(base, app)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return dir, nil
}

// GOOS is exposed for the frontend so it can pick native affordances.
func GOOS() string { return runtime.GOOS }
