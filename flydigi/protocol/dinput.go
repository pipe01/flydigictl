package protocol

import (
	"fmt"
	"io"

	"github.com/rs/zerolog/log"
)

const (
	commandGetDongleVersion       = 17
	commandReadConfig             = 235
	commandGetDeviceInfoInAndroid = 236
	commandReadLEDConfig          = 229
)

type protocolDInput struct {
	rw io.ReadWriter
}

func OpenDInput(rw io.ReadWriter) Protocol {
	return &protocolDInput{rw}
}

func (d *protocolDInput) Read(ch chan<- Message) error {
	buf := make([]byte, 32)

	for {
		n, err := d.rw.Read(buf)
		if err != nil {
			return fmt.Errorf("read from device: %w", err)
		}

		data := buf[:n]

		msg, err := d.resolveUsbData(data)
		if err != nil {
			if err != errUnknownMessage {
				log.Err(err).Msg("failed to handle usb data")
			}
		} else {
			ch <- msg
		}
	}
}

func (d *protocolDInput) Send(cmd Command) error {
	switch cmd := cmd.(type) {
	case CommandGetDongleVersion:
		return d.SendCommand(commandGetDongleVersion)

	case CommandReadConfig:
		return d.SendCommand(commandReadConfig, cmd.ConfigID)

	case CommandReadLEDConfig:
		return d.SendCommand(commandReadLEDConfig, cmd.ConfigID)

	default:
		return ErrUnknownCommand
	}
}

func (d *protocolDInput) SendCommand(cmd byte, args ...byte) error {
	log.Debug().Uint8("cmd", cmd).Bytes("args", args).Msg("sending command")

	buf := make([]byte, 12)
	buf[0] = 5
	buf[1] = byte(cmd)
	copy(buf[2:], args)

	_, err := d.rw.Write(buf)
	return err
}

func (d *protocolDInput) resolveUsbData(p []byte) (msg Message, err error) {
	if p[15] == 235 {
		return MessageGamepadConfigReadCB{
			raw:          p,
			PackageIndex: p[3],
			Data:         p[5:15],
		}, nil
	}

	if p[15] == 229 {
		return MessageLEDConfigReadCB{
			raw:          p,
			PackageIndex: p[3],
			Data:         p[5:15],
		}, nil
	}

	if p[0] == 4 && p[1] == 17 {
		log.Debug().Str("handler", "HandleDongleInfo").Msg("got usb response")
		return MessageDongleInfo{
			raw:  p,
			FW_L: p[2],
			FW_H: p[3],
		}, nil
	}

	if p[15] == 234 || p[15] == 231 || p[15] == 51 {
		return MessageWriteGamepadConfigCBK{
			raw:    p,
			AckNum: p[3],
		}, nil
	}

	if p[15] == 236 {
		return MessageGamePadInfo{
			raw:              p,
			DeviceID:         p[3],
			DeviceMac:        p[5:9],
			FW_L:             p[9],
			FW_H:             p[10],
			Battery:          p[11],
			MotionSensorType: p[14],
			CPUType:          p[12],
			ConnectionType:   p[13],
		}, nil
	}

	if p[3] == 245 && p[4] == 1 {
		log.Debug().Str("handler", "ExtensionChipInfo").Msg("got usb response")
		return nil, errUnknownMessage
	}

	if p[3] == 242 && p[4] == 3 {
		log.Debug().Str("handler", "ScreenInfo").Msg("got usb response")
		return nil, errUnknownMessage
	}

	if p[3] == 242 && p[4] == 4 {
		log.Debug().Str("handler", "ScreenInfoSleepTime").Msg("got usb response")
		return nil, errUnknownMessage
	}

	if p[0] == 4 && (p[1] == 240 || p[1] == 35) {
		log.Debug().Str("handler", "PicData").Msg("got usb response")
		return nil, errUnknownMessage
	}

	if p[0] == 90 && p[1] == 165 && p[2] == 209 && p[4] == 0 {
		log.Debug().Str("handler", "WritePicData").Msg("got usb response")
		return nil, errUnknownMessage
	}

	if p[3] == 250 && p[4] == 160 {
		log.Debug().Str("handler", "UUIDCB").Msg("got usb response")
		return nil, errUnknownMessage
	}

	return nil, errUnknownMessage
}
