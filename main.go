package main

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/karalabe/usb"
	"github.com/rs/zerolog/log"
)

func xinputTest() error {
	devs, err := usb.EnumerateRaw(0x045e, 0x028e)
	if err != nil {
		return fmt.Errorf("enumerate devices: %w", err)
	}

	dev, err := devs[0].Open()
	if err != nil {
		return fmt.Errorf("open device: %w", err)
	}
	defer dev.Close()

	buf := make([]byte, 4096)

	var count atomic.Uint32

	go func() {
		for range time.Tick(1 * time.Second) {
			println(count.Swap(0))
		}
	}()

	pkg := make([]byte, 15)
	pkg[0] = 165
	pkg[1] = 32

	var sum byte
	for i, v := range pkg {
		if i == len(pkg)-1 {
			pkg[i] = sum
		} else {
			sum += v
		}
	}

	dev.Write(pkg)

	for {
		n, err := dev.Read(buf)
		if err != nil {
			break
		}

		count.Add(1)

		data := buf[:n]

		if data[14] == 165 {
			b := data[15]

			if b == 16 {
				// GamePadInfo

			}
		}
	}

	return nil
}

func main() {
	if err := xinputTest(); err != nil {
		log.Fatal().Err(err).Msg("failed to run test")
	}

	// dev, err := flydigi.OpenGamepad()
	// if err != nil {
	// 	log.Fatal().Err(err).Msg("failed to connect to gamepad")
	// }
	// defer dev.Close()

	// cfg, err := dev.GetConfig()
	// if err != nil {
	// 	log.Fatal().Err(err).Msg("failed to read configuration")
	// }

	// _, err = dev.GetLEDConfig()
	// if err != nil {
	// 	log.Fatal().Err(err).Msg("failed to read LED configuration")
	// }

	// if true {
	// 	cfg.Basic.NewLedConfig.SetSteady(config.LedUnit{
	// 		R: 255,
	// 		G: 0,
	// 		B: 0,
	// 	})

	// 	dev.SaveConfig(cfg)
	// }

	// f, _ := os.Create("config.dmp")
	// defer f.Close()

	// d := dump.NewDumper(f, 0)
	// d.Options.MaxDepth = 100
	// d.Dump(cfg)
	// return

	// // cfg.JoyMapping.LeftJoystic.Curve.Zero = 50

	// // dev.SaveConfig(cfg)

	// log.Print("done")
	// select {}
}
