package main

import (
	"flydigi-linux/flydigi"
	"time"

	"github.com/rs/zerolog/log"

	"pault.ag/go/modprobe"
)

func main() {
	err := modprobe.Remove("xpad")
	if err == nil {
		// defer func() {
		// 	err = modprobe.Load("xpad", "")
		// 	if err != nil {
		// 		log.Err(err).Msg("failed to load xpad module")
		// 	}
		// }()
	}
	defer time.Sleep(1 * time.Second)

	dev, err := flydigi.OpenGamepad()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to gamepad")
	}
	defer dev.Close()

	cfg, err := dev.GetConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to read configuration")
	}

	_, err = dev.GetLEDConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to read LED configuration")
	}

	if true {
		cfg.JoyMapping.LeftJoystic.Curve.Zero = 1

		// cfg.Basic.NewLedConfig.SetSteady(config.LedUnit{
		// 	R: 0,
		// 	G: 0,
		// 	B: 255,
		// })
		cfg.Basic.NewLedConfig.SetStreamlined()

		dev.SaveConfig(cfg)
	}
}
