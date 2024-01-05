package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/sstallion/go-hid"
)

type CommandNumber byte

const (
	CommandGetDeviceInfoInAndroid CommandNumber = 236
	CommandReadConfig             CommandNumber = 235
)

type commandCallbackFunc func(data []byte)

type Gamepad struct {
	dev *hid.Device

	cmdCallbacks map[CommandNumber]commandCallbackFunc
}

func OpenGamepad() (*Gamepad, error) {
	var devices []string

	hid.Enumerate(0x04b4, 0x2412, func(info *hid.DeviceInfo) error {
		devices = append(devices, info.Path)
		return nil
	})

	if len(devices) == 0 {
		return nil, errors.New("can't find gamepad")
	}

	dev, err := hid.OpenPath(devices[3])
	if err != nil {
		return nil, fmt.Errorf("open device: %w", err)
	}

	return &Gamepad{
		dev:          dev,
		cmdCallbacks: make(map[CommandNumber]commandCallbackFunc),
	}, nil
}

func (g *Gamepad) Close() error {
	return g.dev.Close()
}

func (g *Gamepad) readLoop() {
	buf := make([]byte, 32)

	for {
		n, err := g.dev.Read(buf)
		if err != nil {
			break
		}

		data := buf[:n]

		if err := g.resolveUsbData(data); err != nil {
			log.Err(err).Msg("failed to handle usb data")
		}
	}
}

func (g *Gamepad) callHandler(cmd CommandNumber, data []byte) error {
	hnd, ok := g.cmdCallbacks[cmd]
	if ok {
		hnd(data)
	}

	return nil
}

func (g *Gamepad) resolveUsbData(p []byte) error {
	if p[15] == 235 {
		log.Debug().Str("handler", "GamepadConfigReadCB").Msg("got usb response")
		println(p[3])
		return g.callHandler(CommandReadConfig, p)
	}

	if p[15] == 229 {
		log.Debug().Str("handler", "LEDConfigReadCB").Msg("got usb response")
		return nil
	}

	if p[0] == 4 && p[1] == 17 {
		log.Debug().Str("handler", "HandleDongleInfo").Msg("got usb response")
		return nil
	}

	if p[15] == 234 || p[15] == 231 || p[15] == 51 {
		log.Debug().Str("handler", "WriteGamepadConfigCBK").Msg("got usb response")
		return nil
	}

	if p[15] == 236 {
		//TODO: Check acks
		log.Debug().Str("handler", "GamePadInfo").Msg("got usb response")
		return g.callHandler(CommandGetDeviceInfoInAndroid, p)
	}

	if p[3] == 245 && p[4] == 1 {
		//TODO: Check acks
		log.Debug().Str("handler", "ExtensionChipInfo").Msg("got usb response")
		return nil
	}

	if p[3] == 242 && p[4] == 3 {
		log.Debug().Str("handler", "ScreenInfo").Msg("got usb response")
		return nil
	}

	if p[3] == 242 && p[4] == 4 {
		log.Debug().Str("handler", "ScreenInfoSleepTime").Msg("got usb response")
		return nil
	}

	if p[0] == 4 && (p[1] == 240 || p[1] == 35) {
		log.Debug().Str("handler", "PicData").Msg("got usb response")
		return nil
	}

	if p[0] == 90 && p[1] == 165 && p[2] == 209 && p[4] == 0 {
		log.Debug().Str("handler", "WritePicData").Msg("got usb response")
		return nil
	}

	if p[3] == 250 && p[4] == 160 {
		log.Debug().Str("handler", "UUIDCB").Msg("got usb response")
		return nil
	}

	return nil
}

func (g *Gamepad) write(p []byte) error {
	_, err := g.dev.Write(p)
	return err
}

func (g *Gamepad) SendCommand(cmd CommandNumber, args ...byte) (resp []byte, err error) {
	ch := make(chan []byte)

	g.cmdCallbacks[cmd] = func(data []byte) {
		ch <- data
	}
	defer delete(g.cmdCallbacks, cmd)

	buf := make([]byte, 12)
	buf[0] = 5
	buf[1] = byte(cmd)
	copy(buf[2:], args)

	retriesLeft := 5

	for retriesLeft > 0 {
		_, err = g.dev.Write(buf)
		if err != nil {
			return nil, fmt.Errorf("write data: %w", err)
		}

		select {
		case resp := <-ch:
			return resp, nil

		case <-time.After(3 * time.Second):
			retriesLeft--
		}
	}

	return nil, errors.New("received no response")
}

func main() {
	dev, err := OpenGamepad()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to gamepad")
	}
	defer dev.Close()

	go dev.readLoop()

	resp, err := dev.SendCommand(CommandReadConfig, 0)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to send command")
	}

	packageIndex := resp[3]
	println(packageIndex)

	select {}
}
