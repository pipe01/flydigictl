package dinput

import (
	"errors"
	"flydigi-linux/flydigi/protocol"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	packageLength    = 52
	ledPackageLength = 49
	maxParcelLength  = 10
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

	transLock     sync.Mutex
	sendingConfig *configTrasmission

	configData    [packageLength * 10]byte
	ledConfigData [ledPackageLength * 10]byte
}

func Open(rw io.ReadWriteCloser) protocol.Protocol {
	p := &protocolDInput{
		rw:    rw,
		msgch: make(chan protocol.Message, 10),
	}
	go p.readLoop()

	return p
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
		return d.sendCommand(commandReadConfig, cmd.ConfigID)

	case protocol.CommandReadLEDConfig:
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
	g.transLock.Lock()
	defer g.transLock.Unlock()

	var chunks [][]byte
	if isLED {
		chunks = getLEDConfigDataParcels(data, configID)
	} else {
		chunks = getConfigDataParcels(data, configID)
	}

	g.sendingConfig = &configTrasmission{
		chunks: chunks,
		ackch:  make(chan int),
	}
	defer func() {
		g.sendingConfig = nil
	}()

	for i, chunk := range chunks {
		retriesLeft := 3
		success := false

		for !success && retriesLeft > 0 {
			retriesLeft--

			_, err := g.rw.Write(chunk)
			if err != nil {
				return fmt.Errorf("write chunk: %w", err)
			}

			for {
				select {
				case ack := <-g.sendingConfig.ackch:
					if ack < i {
						continue // Invalid ack number
					}
					success = true

				case <-time.After(1 * time.Second):
				}

				break
			}
		}

		if !success {
			return errors.New("device didn't respond")
		}
	}

	return nil
}

func (d *protocolDInput) sendCommand(cmd byte, args ...byte) error {
	log.Debug().Uint8("cmd", cmd).Bytes("args", args).Msg("sending command")

	buf := make([]byte, 12)
	buf[0] = 5
	buf[1] = byte(cmd)
	copy(buf[2:], args)

	_, err := d.rw.Write(buf)
	return err
}

func (d *protocolDInput) resolveUsbData(p []byte) (msg protocol.Message, ok bool) {
	if p[15] == 235 {
		if d.handleConfigRead(p, packageLength, d.configData[:]) {
			return protocol.MessageGamepadConfigReadCB{
				Data: d.configData[:],
			}, true
		}

		return nil, false
	}

	if p[15] == 229 {
		if d.handleConfigRead(p, ledPackageLength, d.ledConfigData[:]) {
			return protocol.MessageLEDConfigReadCB{
				Data: d.ledConfigData[:],
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
		if d.sendingConfig != nil {
			select {
			case d.sendingConfig.ackch <- int(p[3]):
			default:
			}
		}

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

func (d *protocolDInput) handleConfigRead(p []byte, packageCount int32, fullData []byte) (isDone bool) {
	packageIndex := int32(p[3])
	if packageIndex > packageCount {
		return false
	}

	println(packageIndex)
	data := p[5:15]

	copy(fullData[packageIndex*10:], data)

	if packageIndex == packageCount {
		// Wait for transmission to finish
		//TODO: Do it better
		time.Sleep(200 * time.Millisecond)
		return true
	}

	return false
}
