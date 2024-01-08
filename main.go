package main

import (
	"flydigi-linux/flydigi"
	"flydigi-linux/flydigi/config"
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"pault.ag/go/modprobe"

	golog "log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	golog.SetOutput(io.Discard) // Supress github.com/google/gousb logging

	err := modprobe.Remove("xpad")
	if err == nil {
		log.Info().Msg("unloaded xpad module")

		defer func() {
			log.Info().Msg("loading xpad module")

			err = modprobe.Load("xpad", "")
			if err != nil {
				log.Err(err).Msg("failed to load xpad module")
			}
		}()
	}

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

		cfg.Basic.NewLedConfig.SetSteady(config.LedUnit{
			R: 0,
			G: 255,
			B: 0,
		})
		//cfg.Basic.NewLedConfig.SetStreamlined()

		dev.SaveConfig(cfg)
	}
}
