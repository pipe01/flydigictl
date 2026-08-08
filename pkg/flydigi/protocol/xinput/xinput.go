package xinput

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"sync/atomic"
	"time"

	"github.com/pipe01/flydigictl/pkg/flydigi/protocol"
	"github.com/pipe01/flydigictl/pkg/flydigi/protocol/internal"
	"github.com/pipe01/flydigictl/pkg/utils"

	"github.com/google/gousb"
	"github.com/rs/zerolog/log"
	"pault.ag/go/modprobe"
)

const (
	packageLength    = 52
	ledPackageLength = 49
)

const (
	commandGetDongleVersion = 17
	commandReadConfig       = 33
	commandGetDeviceInfo    = 16
	commandReadLEDConfig    = 38
)

type protocolXInput struct {
	in     *gousb.InEndpoint
	out    *gousb.OutEndpoint
	closer io.Closer
	dev    *gousb.Device
	usbctx *gousb.Context

	isClosed atomic.Bool

	msgch chan protocol.Message

	configReader, ledConfigReader *internal.ConfigReader

	configWriter *internal.ConfigWriter
}

func Open() (protocol.Protocol, error) {
	ctx := gousb.NewContext()

	var closers utils.MultiCloser

	devs, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		return desc.Vendor == 0x045e && desc.Product == 0x028e
	})
	if err != nil {
		ctx.Close()
		return nil, fmt.Errorf("enumerate devices: %w", err)
	}

	log.Debug().Int("count", len(devs)).Msg("found xinput usb devices")

	if len(devs) == 0 {
		ctx.Close()
		return nil, protocol.ErrGamepadNotPresent
	}

	dev := devs[0]

	cfg, err := dev.Config(1)
	if err != nil {
		return nil, fmt.Errorf("open configuration: %w", err)
	}
	closers.AddCloser(cfg)

	if err := modprobe.Remove("xpad"); err == nil {
		log.Debug().Msg("unloaded xpad module")
	}

	intf, err := cfg.Interface(0, 0)
	if err != nil {
		return nil, fmt.Errorf("open interface: %w", err)
	}
	closers.AddFunc(intf.Close)

	outep, err := intf.OutEndpoint(5)
	if err != nil {
		return nil, fmt.Errorf("open out endpoint: %w", err)
	}

	inep, err := intf.InEndpoint(1)
	if err != nil {
		return nil, fmt.Errorf("open in endpoint: %w", err)
	}

	p := &protocolXInput{
		in:              inep,
		out:             outep,
		closer:          &closers,
		dev:             dev,
		usbctx:          ctx,
		msgch:           make(chan protocol.Message, 10),
		configReader:    internal.NewConfigReader(packageLength, 10),
		ledConfigReader: internal.NewConfigReader(ledPackageLength, 10),
		configWriter:    internal.NewConfigWriter(outep),
	}
	go p.readLoop()

	return p, nil
}

func (d *protocolXInput) Close() error {
	if d.isClosed.Swap(true) {
		return nil
	}

	err := d.closer.Close()

	// Reset the device so it re-enumerates cleanly and xpad can pick it up
	// again; without this the controller is often left unresponsive until
	// power-cycled by hand.
	if rerr := d.dev.Reset(); rerr != nil {
		log.Debug().Err(rerr).Msg("failed to reset usb device")
	}
	d.dev.Close()
	d.usbctx.Close()

	// Always reload xpad, even if this session wasn't the one that unloaded
	// it (e.g. a previous crashed session left it out) — otherwise the
	// controller stays dead as a gamepad until someone loads it by hand.
	// The system modprobe is used because the Go library can't load the
	// compressed modules (.ko.xz/.ko.zst) that most distros ship nowadays.
	log.Debug().Msg("loading xpad module")
	if out, lerr := exec.Command("modprobe", "xpad").CombinedOutput(); lerr != nil {
		log.Err(lerr).Str("output", string(out)).Msg("failed to load xpad module")
	}

	return err
}

func (d *protocolXInput) Messages() <-chan protocol.Message {
	return d.msgch
}

func (d *protocolXInput) readLoop() {
	buf := make([]byte, 100)

	defer close(d.msgch)

	for {
		n, err := d.in.Read(buf)
		if err != nil {
			if status, ok := err.(gousb.TransferStatus); !ok || status != gousb.TransferNoDevice {
				log.Err(err).Msg("failed to read data from usb")
			}

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
	// Vendor packets are at least 28 bytes; anything shorter is a regular
	// gamepad report and indexing it below would panic.
	if len(p) < 28 {
		return nil, false
	}

	if p[14] == 0xA5 {
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

func (d *protocolXInput) Send(ctx context.Context, cmd protocol.Command) error {
	switch cmd := cmd.(type) {
	case protocol.CommandGetDongleVersion:
		return d.sendCommand(ctx, commandGetDongleVersion)

	case protocol.CommandGetDeviceInfo:
		return d.sendCommand(ctx, commandGetDeviceInfo)

	case protocol.CommandReadConfig:
		d.configReader.Reset()
		return d.sendCommand(ctx, commandReadConfig, cmd.ConfigID)

	case protocol.CommandReadLEDConfig:
		d.ledConfigReader.Reset()
		return d.sendCommand(ctx, commandReadLEDConfig, cmd.ConfigID)

	case protocol.CommandSendConfig:
		return d.sendConfig(ctx, cmd.Data, cmd.ConfigID, false)

	case protocol.CommandSendLEDConfig:
		return d.sendConfig(ctx, cmd.Data, cmd.ConfigID, true)

	default:
		return protocol.ErrUnknownCommand
	}
}

func (d *protocolXInput) sendCommand(ctx context.Context, cmd byte, args ...byte) error {
	log.Debug().Uint8("cmd", cmd).Bytes("args", args).Msg("sending command")

	pkg := make([]byte, 15)
	pkg[0] = 165
	pkg[1] = cmd
	copy(pkg[2:], args)

	_, err := d.out.WriteContext(ctx, crcData(pkg))
	return err
}

func (g *protocolXInput) sendConfig(ctx context.Context, data []byte, configID byte, isLED bool) error {
	var chunks [][]byte
	if isLED {
		chunks = getLEDConfigDataParcels(data, configID)
	} else {
		chunks = getConfigDataParcels(data, configID)
	}

	return g.configWriter.Send(ctx, chunks, 3, 3*time.Second)
}
