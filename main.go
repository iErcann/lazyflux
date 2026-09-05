// lazyflux - set the screen color temperature, from the command line or a tray icon.
//
// Usage:
//
//	lazyflux                 start the tray icon (re-applies the last saved setting)
//	lazyflux TEMP            set temperature in kelvin, e.g. lazyflux 4500
//	lazyflux TEMP BRIGHT     also set brightness 0.1..1.0, e.g. lazyflux 4500 0.8
package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	minTemp     = 1000
	maxTemp     = 10000
	defaultTemp = 6500
)

// Current temperature and brightness, saved to disk.
type Setting struct {
	Temp       int
	Brightness float64
}

// Sends the setting to the screen, picking X11 or GNOME.
func Apply(s Setting) error {
	s.Temp = min(max(s.Temp, minTemp), maxTemp)
	s.Brightness = min(max(s.Brightness, 0.1), 1.0)

	switch {
	case strings.Contains(os.Getenv("XDG_CURRENT_DESKTOP"), "GNOME"):
		return setGnome(s.Temp)
	case os.Getenv("WAYLAND_DISPLAY") != "":
		return errors.New("Wayland is only supported on GNOME")
	default:
		return setX11(s.Temp, s.Brightness)
	}
}

func stateFile() string {
	dir, _ := os.UserConfigDir()
	return filepath.Join(dir, "lazyflux", "state")
}

// Reads the saved setting, or the default if there is none.
func Load() Setting {
	s := Setting{Temp: defaultTemp, Brightness: 1.0}
	data, err := os.ReadFile(stateFile())
	if err == nil {
		fmt.Sscan(string(data), &s.Temp, &s.Brightness)
	}
	return s
}

func Save(s Setting) {
	os.MkdirAll(filepath.Dir(stateFile()), 0o755)
	os.WriteFile(stateFile(), []byte(fmt.Sprintf("%d %.2f\n", s.Temp, s.Brightness)), 0o644)
}

func main() {
	if len(os.Args) == 1 {
		runTray()
		return
	}

	s := Load()
	temp, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "usage: lazyflux [temperature] [brightness]\n")
		os.Exit(2)
	}
	s.Temp = temp
	if len(os.Args) > 2 {
		s.Brightness, _ = strconv.ParseFloat(os.Args[2], 64)
	}

	if err := Apply(s); err != nil {
		fmt.Fprintln(os.Stderr, "lazyflux:", err)
		os.Exit(1)
	}
	Save(s)
}
