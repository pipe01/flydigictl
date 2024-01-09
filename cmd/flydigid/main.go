package main

import (
	"flag"
	"io"
	"os"

	"github.com/pipe01/flydigi-linux/pkg/dbus/server"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	golog "log"
)

func main() {
	prettyLogging := flag.Bool("pretty-logs", false, "Enable human-readable colored logs")
	flag.Parse()

	if *prettyLogging {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	golog.SetOutput(io.Discard) // Supress github.com/google/gousb logging

	srv := server.New()

	if err := srv.Listen(); err != nil {
		log.Fatal().Err(err).Msg("failed to start dbus server")
	}
}
