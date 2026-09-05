package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"strconv"

	"fyne.io/systray"
)

var presets = []struct {
	Name string
	Temp int
}{
	{"Daylight", 6500},
	{"Neutral", 5500},
	{"Warm", 4500},
	{"Evening", 3500},
	{"Night", 2700},
}

var (
	current     Setting
	status      *systray.MenuItem   // greyed-out first line showing the live state
	presetItems []*systray.MenuItem // ticked when their temperature is active
)

func runTray() {
	systray.Run(onReady, nil)
}

func onReady() {
	systray.SetTitle("lazyflux")

	header := systray.AddMenuItem("lazyflux", "")
	header.Disable()
	status = systray.AddMenuItem("", "")
	status.Disable()
	systray.AddSeparator()

	for _, p := range presets {
		item := systray.AddMenuItemCheckbox(fmt.Sprintf("%s  %d K", p.Name, p.Temp), "", false)
		presetItems = append(presetItems, item)
		go func() {
			for range item.ClickedCh {
				current.Temp = p.Temp
				set(current)
			}
		}()
	}
	addAction("Temperature slider...", func() {
		slide("Temperature (K)", minTemp, maxTemp, 100, current.Temp, func(v int) {
			current.Temp = v
			set(current)
		})
	})

	systray.AddSeparator()
	addAction("Brightness slider...", func() {
		slide("Brightness (%)", 10, 100, 5, int(current.Brightness*100), func(v int) {
			current.Brightness = float64(v) / 100
			set(current)
		})
	})

	systray.AddSeparator()
	login := systray.AddMenuItemCheckbox("Start at login", "", autostartEnabled())
	go func() {
		for range login.ClickedCh {
			if setAutostart(!login.Checked()) == nil {
				toggle(login)
			}
		}
	}()
	addAction("Quit", systray.Quit)

	current = Load()
	set(current)
}

// Applies a setting, saves it and updates the menu.
func set(s Setting) {
	if err := Apply(s); err != nil {
		status.SetTitle("Error: " + err.Error())
		return
	}
	Save(s)
	state := fmt.Sprintf("%d K  ·  %.0f%%", s.Temp, s.Brightness*100)
	status.SetTitle(state)
	systray.SetTooltip("lazyflux  " + state)
	systray.SetIcon(iconPNG(s.Temp))
	for i, p := range presets {
		if p.Temp == s.Temp {
			presetItems[i].Check()
		} else {
			presetItems[i].Uncheck()
		}
	}
}

func toggle(item *systray.MenuItem) {
	if item.Checked() {
		item.Uncheck()
	} else {
		item.Check()
	}
}

func addAction(title string, fn func()) {
	item := systray.AddMenuItem(title, "")
	go func() {
		for range item.ClickedCh {
			fn()
		}
	}()
}

// Opens a zenity slider and calls fn as it moves. Cancel restores start.
func slide(label string, lo, hi, step, start int, fn func(int)) {
	cmd := exec.Command("zenity", "--scale", "--title=lazyflux", "--text="+label,
		"--min-value="+strconv.Itoa(lo), "--max-value="+strconv.Itoa(hi),
		"--step="+strconv.Itoa(step), "--value="+strconv.Itoa(start), "--print-partial")
	out, err := cmd.StdoutPipe()
	if err != nil || cmd.Start() != nil {
		return
	}
	for sc := bufio.NewScanner(out); sc.Scan(); {
		if v, err := strconv.Atoi(sc.Text()); err == nil {
			fn(v)
		}
	}
	if cmd.Wait() != nil { // non-zero exit means the dialog was cancelled
		fn(start)
	}
}
