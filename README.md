<p align="center"><img src="packaging/lazyflux.svg" width="96" alt="lazyflux"></p>

# lazyflux

A small color-temperature control for Linux.

No schedules, location tracking, profiles, or daemon to configure. Set a temperature from the tray or the command line and lazyflux remembers it.

```sh
lazyflux              # start the tray, restoring the last setting
lazyflux 4500         # set color temperature to 4500 K
lazyflux 4500 0.8     # 4500 K at 80% brightness
lazyflux 6500         # reset to 6500 K
```

State is stored in:

```text
~/.config/lazyflux/state
```

## Install

Download the `.deb` from the [latest release](https://github.com/iercann/lazyflux/releases/latest) and double-click it (Ubuntu, Debian, Mint, Pop). lazyflux then appears in your application menu. Other distributions get a `.tar.gz` with the binary.

With Go installed:

```sh
go install github.com/iercann/lazyflux@latest
```

Run `lazyflux`, then enable **Start at login** from the tray menu if you want it to start with your desktop.

Don't run lazyflux alongside f.lux, GNOME Night Light, Redshift, Gammastep, or another color-temperature tool. They'll overwrite each other's display settings.

## Support

| Session            | color temperature | Brightness | Tray                  |
| ------------------ | -----------------: | ---------: | --------------------- |
| X11                |                yes |        yes | yes                   |
| GNOME / Wayland    |                yes |         no | requires AppIndicator |
| KDE / Wayland      |                 no |         no | no                    |
| Sway / Wayland     |                 no |         no | no                    |
| Hyprland / Wayland |                 no |         no | no                    |

### X11

lazyflux sets the display gamma ramps directly through RandR.

### GNOME / Wayland

Wayland doesn't give normal applications direct access to gamma ramps. On GNOME, lazyflux uses GNOME's color-temperature controls through `gsettings`.

Setting `6500 K` restores the normal display temperature.

Other Wayland compositors are currently unsupported.

## Tray controls

The tray menu has presets and optional sliders.

Sliders require [`zenity`](https://gitlab.gnome.org/GNOME/zenity), which is already installed on many desktop systems. If `zenity` isn't available, the presets and command-line interface still work.
