package dinput

import (
	"flydigi-linux/flydigi/protocol"
	"flydigi-linux/flydigi/protocol/internal"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/karalabe/usb"
	"github.com/rs/zerolog/log"
)

const (
	packageLength    = 52
	ledPackageLength = 49
)

const (
	commandGetDongleVersion       = 17
	commandReadConfig             = 235
	commandGetDeviceInfoInAndroid = 236
	commandReadLEDConfig          = 229
)

type configTrasmission struct {
	chunks [][]byte
	ackch  chan int
}

type protocolDInput struct {
	rw    io.ReadWriteCloser
	msgch chan protocol.Message

	configWriter *internal.ConfigWriter

	configReader, ledConfigReader *internal.ConfigReader
}

func Open() (protocol.Protocol, error) {
	devs, err := usb.EnumerateHid(0x04b4, 0x2412)
	if err != nil {
		return nil, fmt.Errorf("enumerate devices: %w", err)
	}

	if len(devs) == 0 {
		return nil, os.ErrNotExist
	}

	p := &protocolDInput{
		msgch:           make(chan protocol.Message, 10),
		configReader:    internal.NewConfigReader(packageLength, 10),
		ledConfigReader: internal.NewConfigReader(ledPackageLength, 10),
	}

	for _, d := range devs {
		if d.Interface == 2 {
			dev, err := d.Open()
			if err != nil {
				return nil, fmt.Errorf("open usb device: %w", err)
			}

			p.rw = dev
			p.configWriter = internal.NewConfigWriter(dev)

			go p.readLoop()

			return p, nil
		}
	}

	return nil, os.ErrNotExist
}

func (d *protocolDInput) Close() error {
	err := d.rw.Close()
	close(d.msgch)
	return err
}

func (d *protocolDInput) Messages() <-chan protocol.Message {
	return d.msgch
}

func (d *protocolDInput) readLoop() {
	buf := make([]byte, 32)

	for {
		n, err := d.rw.Read(buf)
		if err != nil {
			break
		}

		data := buf[:n]

		msg, ok := d.resolveUsbData(data)
		if ok {
			d.msgch <- msg
		}
	}
}

func (d *protocolDInput) Send(cmd protocol.Command) error {
	switch cmd := cmd.(type) {
	case protocol.CommandGetDongleVersion:
		return d.sendCommand(commandGetDongleVersion)

	case protocol.CommandReadConfig:
		d.configReader.Reset()
		return d.sendCommand(commandReadConfig, cmd.ConfigID)

	case protocol.CommandReadLEDConfig:
		d.ledConfigReader.Reset()
		return d.sendCommand(commandReadLEDConfig, cmd.ConfigID)

	case protocol.CommandSendConfig:
		return d.sendConfig(cmd.Data, cmd.ConfigID, false)

	case protocol.CommandSendLEDConfig:
		return d.sendConfig(cmd.Data, cmd.ConfigID, true)

	default:
		return protocol.ErrUnknownCommand
	}
}

func (g *protocolDInput) sendConfig(data []byte, configID byte, isLED bool) error {
	var chunks [][]byte
	if isLED {
		chunks = getLEDConfigDataParcels(data, configID)
	} else {
		chunks = getConfigDataParcels(data, configID)
	}

	return g.configWriter.Send(chunks, 3, 3*time.Second)
}

func (d *protocolDInput) sendCommand(cmd byte, args ...byte) error {
	log.Debug().Uint8("cmd", cmd).Bytes("args", args).Msg("sending command")

	buf := make([]byte, 12)
	buf[0] = 5
	buf[1] = cmd
	copy(buf[2:], args)

	_, err := d.rw.Write(buf)
	return err
}

func (d *protocolDInput) resolveUsbData(p []byte) (msg protocol.Message, ok bool) {
	if p[15] == 235 {
		d.configReader.GotPackage(int(p[3]), p[5:15])

		if d.configReader.IsFinished() {
			// Wait for transmission to finish
			//TODO: Do it better
			time.Sleep(200 * time.Millisecond)

			return protocol.MessageGamepadConfigReadCB{
				Data: d.configReader.Data(),
			}, true
		}

		return nil, false
	}

	if p[15] == 229 {
		d.ledConfigReader.GotPackage(int(p[3]), p[5:15])

		if d.ledConfigReader.IsFinished() {
			// Wait for transmission to finish
			//TODO: Do it better
			time.Sleep(200 * time.Millisecond)

			return protocol.MessageLEDConfigReadCB{
				Data: d.ledConfigReader.Data(),
			}, true
		}

		return nil, false
	}

	if p[0] == 4 && p[1] == 17 {
		log.Debug().Str("handler", "HandleDongleInfo").Msg("got usb response")
		return protocol.MessageDongleInfo{
			FW_L: p[2],
			FW_H: p[3],
		}, true
	}

	if p[15] == 234 || p[15] == 231 || p[15] == 51 {
		d.configWriter.Ack(int(p[3]))

		return nil, false
	}

	if p[15] == 236 {
		return protocol.MessageGamePadInfo{
			DeviceID:         p[3],
			DeviceMac:        p[5:9],
			FW_L:             p[9],
			FW_H:             p[10],
			Battery:          p[11],
			MotionSensorType: p[14],
			CPUType:          p[12],
			ConnectionType:   p[13],
		}, true
	}

	if p[3] == 245 && p[4] == 1 {
		log.Debug().Str("handler", "ExtensionChipInfo").Msg("got usb response")
		return nil, false
	}

	if p[3] == 242 && p[4] == 3 {
		log.Debug().Str("handler", "ScreenInfo").Msg("got usb response")
		return nil, false
	}

	if p[3] == 242 && p[4] == 4 {
		log.Debug().Str("handler", "ScreenInfoSleepTime").Msg("got usb response")
		return nil, false
	}

	if p[0] == 4 && (p[1] == 240 || p[1] == 35) {
		log.Debug().Str("handler", "PicData").Msg("got usb response")
		return nil, false
	}

	if p[0] == 90 && p[1] == 165 && p[2] == 209 && p[4] == 0 {
		log.Debug().Str("handler", "WritePicData").Msg("got usb response")
		return nil, false
	}

	if p[3] == 250 && p[4] == 160 {
		log.Debug().Str("handler", "UUIDCB").Msg("got usb response")
		return nil, false
	}

	return nil, false
}
