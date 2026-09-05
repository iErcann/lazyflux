package main

import (
	"os"
	"path/filepath"
)

// Standard autostart entry, works on every desktop.
func autostartFile() string {
	dir, _ := os.UserConfigDir()
	return filepath.Join(dir, "autostart", "lazyflux.desktop")
}

func autostartEnabled() bool {
	_, err := os.Stat(autostartFile())
	return err == nil
}

// Adds or removes the autostart entry for this binary.
func setAutostart(enable bool) error {
	if !enable {
		return os.Remove(autostartFile())
	}
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	entry := "[Desktop Entry]\nType=Application\nName=lazyflux\nExec=" + exe + "\nTerminal=false\n"
	os.MkdirAll(filepath.Dir(autostartFile()), 0o755)
	return os.WriteFile(autostartFile(), []byte(entry), 0o644)
}
