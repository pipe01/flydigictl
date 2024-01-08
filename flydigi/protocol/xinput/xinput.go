package xinput

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
	commandReadConfig             = 33
	commandGetDeviceInfoInAndroid = 16
	commandReadLEDConfig          = 38
)

type protocolXInput struct {
	rw    io.ReadWriteCloser
	msgch chan protocol.Message

	configReader, ledConfigReader *internal.ConfigReader

	configWriter *internal.ConfigWriter
}

func Open() (protocol.Protocol, error) {
	devs, err := usb.EnumerateRaw(0x045e, 0x028e)
	if err != nil {
		return nil, fmt.Errorf("enumerate devices: %w", err)
	}

	if len(devs) == 0 {
		return nil, os.ErrNotExist
	}

	dev, err := devs[0].Open()
	if err != nil {
		return nil, fmt.Errorf("open usb device: %w", err)
	}

	p := &protocolXInput{
		rw:              dev,
		msgch:           make(chan protocol.Message, 10),
		configReader:    internal.NewConfigReader(packageLength, 10),
		ledConfigReader: internal.NewConfigReader(ledPackageLength, 10),
		configWriter:    internal.NewConfigWriter(dev),
	}
	go p.readLoop()

	return p, nil
}

func (d *protocolXInput) Close() error {
	err := d.rw.Close()
	close(d.msgch)
	return err
}

func (d *protocolXInput) Messages() <-chan protocol.Message {
	return d.msgch
}

func (d *protocolXInput) readLoop() {
	buf := make([]byte, 100)

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

func (d *protocolXInput) resolveUsbData(p []byte) (protocol.Message, bool) {
	if p[14] == 165 {
		switch p[15] {
		case 16:
			return protocol.MessageGamePadInfo{
				DeviceID:         p[16],
				DeviceMac:        p[17:21],
				FW_L:             p[21],
				FW_H:             p[22],
				Battery:          p[23],
				CPUType:          p[24],
				ConnectionType:   p[25],
				MotionSensorType: p[26],
			}, true

		case 17:
			return protocol.MessageDongleInfo{
				FW_L: p[16],
				FW_H: p[17],
			}, true

		case 32:
			// HandleGamepadConfigId

		case 34:
			// HandleGamepadConfigReadCB
			d.configReader.GotPackage(int(p[16]), p[17:28])

			if d.configReader.IsFinished() {
				time.Sleep(200 * time.Millisecond)
				return protocol.MessageGamepadConfigReadCB{
					Data: d.configReader.Data(),
				}, true
			}

		case 35, 37: // HandleStartWriteGamepadConfig
			d.configWriter.Ack(0)

		case 36: // HandleWriteGamepadConfigCBK
			d.configWriter.Ack(int(p[16]))

		case 39:
			// HandleLedConfigReadCB
			d.ledConfigReader.GotPackage(int(p[16]), p[17:28])

			if d.ledConfigReader.IsFinished() {
				time.Sleep(200 * time.Millisecond)
				return protocol.MessageLEDConfigReadCB{
					Data: d.ledConfigReader.Data(),
				}, true
			}

		case 41: // HandleWriteLEDConfigCBK
			d.configWriter.Ack(int(p[16]))

		case 42: // HandleStartWriteLEDConfig
			d.configWriter.Ack(0)
		}
	}

	return nil, false
}

func (d *protocolXInput) Send(cmd protocol.Command) error {
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

func (d *protocolXInput) sendCommand(cmd byte, args ...byte) error {
	log.Debug().Uint8("cmd", cmd).Bytes("args", args).Msg("sending command")

	pkg := make([]byte, 15)
	pkg[0] = 165
	pkg[1] = cmd
	copy(pkg[2:], args)

	_, err := d.rw.Write(crcData(pkg))
	return err
}

func (g *protocolXInput) sendConfig(data []byte, configID byte, isLED bool) error {
	var chunks [][]byte
	if isLED {
		chunks = getLEDConfigDataParcels(data, configID)
	} else {
		chunks = getConfigDataParcels(data, configID)
	}

	return g.configWriter.Send(chunks, 3, 3*time.Second)
}
