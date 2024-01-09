package main

import (
	"io"
	"os"

	"github.com/pipe01/flydigi-linux/dbus"
	"github.com/pipe01/flydigi-linux/flydigi"
	"github.com/pipe01/flydigi-linux/flydigi/config"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	golog "log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	golog.SetOutput(io.Discard) // Supress github.com/google/gousb logging

	server := dbus.NewServer()
	server.Listen()

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
