package main

import (
	"math"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/randr"
	"github.com/jezek/xgb/xproto"
)

// Red, green and blue multipliers for a temperature.
func channelGains(temp int) (red, green, blue float64) {
	t := float64(temp)
	if temp < 6500 {
		l := math.Log(t - 700)
		return 1,
			clamp01(-1.47751309139817 + 0.28590164772055*l),
			clamp01(-4.38321650114872 + 0.6212158769447*l)
	}
	l := math.Log(t - 5800)
	return clamp01(1.75390204039018 - 0.1150805671482*l),
		clamp01(1.49221604915144 - 0.07513509588921*l),
		1
}

func clamp01(x float64) float64 { return min(max(x, 0), 1) }

// Writes the gamma ramp to every monitor.
func setX11(temp int, brightness float64) error {
	conn, err := xgb.NewConn()
	if err != nil {
		return err
	}
	defer conn.Close()
	if err := randr.Init(conn); err != nil {
		return err
	}

	r, g, b := channelGains(temp)

	for _, screen := range xproto.Setup(conn).Roots {
		res, err := randr.GetScreenResourcesCurrent(conn, screen.Root).Reply()
		if err != nil {
			return err
		}
		for _, crtc := range res.Crtcs {
			sizeReply, err := randr.GetCrtcGammaSize(conn, crtc).Reply()
			if err != nil || sizeReply.Size == 0 {
				continue // CRTC not in use
			}
			n := int(sizeReply.Size)
			red, green, blue := make([]uint16, n), make([]uint16, n), make([]uint16, n)
			for i := range n {
				v := 65535 * brightness * float64(i) / float64(n)
				red[i] = uint16(v*r + 0.5)
				green[i] = uint16(v*g + 0.5)
				blue[i] = uint16(v*b + 0.5)
			}
			if err := randr.SetCrtcGammaChecked(conn, crtc, uint16(n), red, green, blue).Check(); err != nil {
				return err
			}
		}
	}
	return nil
}
