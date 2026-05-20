package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/pipe01/flydigictl/pkg/dbus/server"
	"github.com/pipe01/flydigictl/pkg/flydigi"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	golog "log"
)

func main() {
	prettyLogging := flag.Bool("pretty-logs", false, "Enable human-readable colored logs")
	protocolMode := flag.String("mode", envDefault("FLYDIGI_MODE", string(flydigi.ProtocolModeAuto)), "Controller protocol mode: auto, dinput, xinput")
	flag.Parse()

	if *prettyLogging {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	golog.SetOutput(io.Discard) // Supress github.com/google/gousb logging

	mode, err := flydigi.ParseProtocolMode(*protocolMode)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	log.Info().Str("mode", string(mode)).Msg("using controller protocol mode")

	srv := server.NewWithProtocolMode(mode)

	if err := srv.Listen(); err != nil {
		log.Fatal().Err(err).Msg("failed to start dbus server")
	}
}

func envDefault(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}

	return fallback
}
