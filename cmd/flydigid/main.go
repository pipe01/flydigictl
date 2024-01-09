package main

import (
	"io"
	"os"

	"github.com/pipe01/flydigi-linux/pkg/dbus/server"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	golog "log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	golog.SetOutput(io.Discard) // Supress github.com/google/gousb logging

	server := server.New()

	if err := server.Listen(); err != nil {
		log.Fatal().Err(err).Msg("failed to start dbus server")
	}
}
