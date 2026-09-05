package main

import (
	"fmt"
	"os/exec"
	"strconv"
)

// Uses GNOME's Night Light. No brightness there. 6500 K turns it off.
func setGnome(temp int) error {
	settings := [][2]string{
		{"night-light-schedule-automatic", "false"},
		{"night-light-schedule-from", "0.0"},
		{"night-light-schedule-to", "24.0"},
		{"night-light-temperature", strconv.Itoa(temp)},
		{"night-light-enabled", strconv.FormatBool(temp < defaultTemp)},
	}
	for _, kv := range settings {
		cmd := exec.Command("gsettings", "set", "org.gnome.settings-daemon.plugins.color", kv[0], kv[1])
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("gsettings set %s: %s", kv[0], out)
		}
	}
	return nil
}
