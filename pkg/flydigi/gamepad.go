package flydigi

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pipe01/flydigictl/pkg/flydigi/config"
	"github.com/pipe01/flydigictl/pkg/flydigi/protocol"
	"github.com/pipe01/flydigictl/pkg/flydigi/protocol/dinput"
	"github.com/pipe01/flydigictl/pkg/flydigi/protocol/xinput"
	"github.com/pipe01/flydigictl/pkg/utils"

	"github.com/rs/zerolog/log"
)

type ProtocolMode string

const (
	ProtocolModeAuto   ProtocolMode = "auto"
	ProtocolModeDInput ProtocolMode = "dinput"
	ProtocolModeXInput ProtocolMode = "xinput"
)

func ParseProtocolMode(mode string) (ProtocolMode, error) {
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case "", "auto":
		return ProtocolModeAuto, nil
	case "dinput", "directinput", "hid":
		return ProtocolModeDInput, nil
	case "xinput":
		return ProtocolModeXInput, nil
	default:
		return "", fmt.Errorf("unknown protocol mode %q, expected auto, dinput, or xinput", mode)
	}
}

type FDGConncetType int32

const (
	FDGConncetUnknow FDGConncetType = iota
	FDGConncetWireless
	FDGConncetWired
)

type FDGDeviceInfo struct {
	ConfigVersion   int32
	KeyCateNum      []int32
	Keys            []int32
	SelCfgId        int32
	SelSwitchCfgId  int32
	FirmwareVersion string
	HidName         string
	BatteryPercent  int32
	DeviceId        int32
	PackageLength   int32
	// ExtChipInfo FDGDeviceExtChipInfo
	// ScreenInfo FDGDeviceScreenInfo
	FirmwareVersionCode int32
	DongleVersion       string
	DeviceMac           string
	MotionSensorType    string
	ConnectType         FDGConncetType
	ConnectMode         string
	IsConnect           bool
	CpuType             string
	CpuName             string
	GameHadleName       string
	ProductName         string
	ShowGameHadleName   string
	ShowEnGameHadleName string
	FirmwareName        string
	DONGLEFirmwareName  string
	LcdFirmwareName     string
	TriggerFirmwareName string
	SIFirmwareVersion   string
	UpgradeType         int32
	LedNum              byte
	ThemeFrColor        string
	ThemeFrHoverColor   string
	ResName             string
	ThemeBgColor        string
	IsIP                bool
}

type CommandNumber byte

const (
	CommandGetDongleVersion       CommandNumber = 17
	CommandReadConfig             CommandNumber = 235
	CommandGetDeviceInfoInAndroid CommandNumber = 236
	CommandReadLEDConfig          CommandNumber = 229
)

type commandCallbackFunc func(data []byte)

type Gamepad struct {
	prot         protocol.Protocol
	protocolMode ProtocolMode

	devInfo *utils.CondValue[FDGDeviceInfo]

	closech chan struct{}

	currConfig    *utils.CondValue[config.AllConfigBean]
	currLEDConfig *utils.CondValue[config.NewLedConfigBean]

	configID atomic.Uint32
}

func OpenGamepad() (*Gamepad, error) {
	return OpenGamepadWithMode(ProtocolModeAuto)
}

func OpenGamepadWithMode(mode ProtocolMode) (*Gamepad, error) {
	var prot protocol.Protocol
	actualMode := mode
	var err error

	switch mode {
	case ProtocolModeAuto:
		prot, err = dinput.Open()
		if err != nil {
			if err != protocol.ErrGamepadNotPresent {
				return nil, fmt.Errorf("open dinput device: %w", err)
			}

			prot, err = xinput.Open()
			if err != nil {
				return nil, fmt.Errorf("open xinput device: %w", err)
			}
			actualMode = ProtocolModeXInput
		} else {
			actualMode = ProtocolModeDInput
		}

	case ProtocolModeDInput:
		prot, err = dinput.Open()
		if err != nil {
			return nil, fmt.Errorf("open dinput device: %w", err)
		}

	case ProtocolModeXInput:
		prot, err = xinput.Open()
		if err != nil {
			return nil, fmt.Errorf("open xinput device: %w", err)
		}

	default:
		return nil, fmt.Errorf("unknown protocol mode %q", mode)
	}

	return newGamepad(prot, actualMode), nil
}

func newGamepad(prot protocol.Protocol, mode ProtocolMode) *Gamepad {
	gamepad := &Gamepad{
		prot:         prot,
		protocolMode: mode,

		closech:       make(chan struct{}),
		devInfo:       utils.NewCondValue[FDGDeviceInfo](&sync.Mutex{}),
		currConfig:    utils.NewCondValue[config.AllConfigBean](&sync.Mutex{}),
		currLEDConfig: utils.NewCondValue[config.NewLedConfigBean](&sync.Mutex{}),
	}
	go gamepad.readLoop()

	return gamepad
}

func (g *Gamepad) Close() error {
	return g.prot.Close()
}

func (g *Gamepad) readLoop() {
	defer close(g.closech)
	defer g.prot.Close()

	for msg := range g.prot.Messages() {
		if err := g.handleMessage(msg); err != nil {
			log.Err(err).Msg("failed to handle usb data")
		}
	}

	log.Debug().Msg("gamepad read loop exited")
}

func (g *Gamepad) NotifyClose(ch chan<- struct{}) {
	go func() {
		<-g.closech
		ch <- struct{}{}
	}()
}

func (g *Gamepad) handleMessage(msg protocol.Message) error {
	switch msg := msg.(type) {
	case protocol.MessageGamePadInfo:
		return g.handleDeviceInfo(msg)

	case protocol.MessageDongleInfo:
		return g.handleDongleInfo(msg)

	case protocol.MessageGamepadConfigID:
		return g.handleGamepadConfigID(msg)

	case protocol.MessageGamepadConfigReadCB:
		return g.handleGamepadConfigRead(msg)

	case protocol.MessageLEDConfigReadCB:
		return g.handleLEDConfigRead(msg)

	default:
		return errors.New("unknown message type")
	}
}

func (g *Gamepad) handleGamepadConfigID(msg protocol.MessageGamepadConfigID) error {
	g.configID.Store(uint32(msg.ConfigID))
	log.Debug().Uint8("config_id", msg.ConfigID).Msg("got selected config id")

	return nil
}

func (g *Gamepad) selectedConfigID() byte {
	return byte(g.configID.Load())
}

func (g *Gamepad) handleDeviceInfo(msg protocol.MessageGamePadInfo) error {
	log.Debug().Uint8("deviceid", msg.DeviceID).Uint8("battery", msg.Battery).Msg("got device info")

	devInfo := FDGDeviceInfo{}

	devInfo.DeviceId = int32(msg.DeviceID)
	devInfo.ConnectMode = string(g.protocolMode)
	// if (!GameHandleListDic.gameHandleDic.ContainsKey((int)deviceId))
	// {
	// 	return;
	// }

	// currentDeviceInfo.GameHadleName = GameHandleListDic.gameHandleDic[(int)deviceId].GameHadleName;
	// currentDeviceInfo.FirmwareName = GameHandleListDic.gameHandleDic[(int)deviceId].FirmwareName;

	devInfo.DeviceMac = net.HardwareAddr(msg.DeviceMac).String()

	fw_l := msg.FW_L & 15
	fw_l_2 := msg.FW_L >> 4
	fw_h := msg.FW_H & 15
	fw_h_2 := msg.FW_H >> 4

	devInfo.FirmwareVersionCode = int32(fw_h_2)*1000 + int32(fw_h)*100 + int32(fw_l_2)*10 + int32(fw_l)
	devInfo.FirmwareVersion = fmt.Sprintf("%d.%d.%d.%d", fw_h_2, fw_h, fw_l_2, fw_l)

	devInfo.BatteryPercent = batteryPercent(msg)

	switch msg.MotionSensorType {
	case 1:
		devInfo.MotionSensorType = "ST"
	case 2:
		devInfo.MotionSensorType = "QST"
	}

	if msg.CPUType > 0 {
		devInfo.CpuType = "wch"
	} else {
		devInfo.CpuType = "nordic"
	}

	if fw_h_2 >= 6 && fw_h >= 1 {
		devInfo.CpuType = "wch"
	}

	if devInfo.CpuType == "wch" {
		if msg.ConnectionType == 1 {
			devInfo.ConnectType = FDGConncetWired
			devInfo.CpuName = "ch573"
		} else {
			devInfo.ConnectType = FDGConncetWireless
			devInfo.CpuName = "ch571"
			g.prot.Send(protocol.CommandGetDongleVersion{})
		}
	}

	devInfo.GameHadleName = config.GameHandleName[devInfo.DeviceId]
	// currentDeviceInfo.GameHadleName = GameHandleListDic.gameHandleDic[currentDeviceInfo.DeviceId].GameHadleName;
	// currentDeviceInfo.ShowGameHadleName = GameHandleListDic.gameHandleDic[currentDeviceInfo.DeviceId].ShowGameHadleName;
	// currentDeviceInfo.FirmwareName = GameHandleListDic.gameHandleDic[currentDeviceInfo.DeviceId].FirmwareName;
	// currentDeviceInfo.DONGLEFirmwareName = GameHandleListDic.gameHandleDic[currentDeviceInfo.DeviceId].DONGLEFirmwareName;
	// currentDeviceInfo.SIFirmwareVersion = GameHandleListDic.gameHandleDic[currentDeviceInfo.DeviceId].SIFirmwareVersion;
	// currentDeviceInfo.ResName = GameHandleListDic.gameHandleDic[currentDeviceInfo.DeviceId].ResName;
	// currentDeviceInfo.IsIP = GameHandleListDic.gameHandleDic[currentDeviceInfo.DeviceId].IsIP;
	// currentDeviceInfo.ThemeBgColor = GameHandleListDic.gameHandleDic[currentDeviceInfo.DeviceId].ThemeBgColor;
	// currentDeviceInfo.ThemeFrColor = GameHandleListDic.gameHandleDic[currentDeviceInfo.DeviceId].ThemeFrColor;
	// currentDeviceInfo.ThemeFrHoverColor = GameHandleListDic.gameHandleDic[currentDeviceInfo.DeviceId].ThemeFrHoverColor;
	// currentDeviceInfo.LedNum = GameHandleListDic.gameHandleDic[currentDeviceInfo.DeviceId].LedNum;
	// currentDeviceInfo.LcdFirmwareName = GameHandleListDic.gameHandleDic[currentDeviceInfo.DeviceId].LcdFirmwareName;
	// currentDeviceInfo.TriggerFirmwareName = GameHandleListDic.gameHandleDic[currentDeviceInfo.DeviceId].TriggerFirmwareName;

	switch devInfo.DeviceId {
	case 19:
		devInfo.FirmwareName = "apex2"

	case 20, 21:
		devInfo.FirmwareName = "f1"
		if devInfo.CpuType == "wch" {
			devInfo.FirmwareName = "f1wch"
		}

	case 22, 23:
		devInfo.FirmwareName = "f1p"

	case 24:
		devInfo.FirmwareName = "k1"

	case 25:
		devInfo.FirmwareName = "fp1"
	}

	g.devInfo.Value = &devInfo
	g.devInfo.Broadcast()

	return nil
}

func batteryPercent(msg protocol.MessageGamePadInfo) int32 {
	if msg.DeviceID == 85 || msg.DeviceID == 105 {
		if msg.Battery <= 4 {
			return int32(msg.Battery) * 25
		}

		if msg.Battery <= 100 {
			return int32(msg.Battery)
		}
	}

	battery := msg.Battery
	const apex2MinBY = 98
	const apex2MaxBY = 114

	if battery < apex2MinBY {
		battery = apex2MinBY
	} else if battery > apex2MaxBY {
		battery = apex2MaxBY
	}

	return int32(100 * float32(battery-apex2MinBY) / float32(apex2MaxBY-apex2MinBY))
}

func (g *Gamepad) handleDongleInfo(msg protocol.MessageDongleInfo) error {
	fw_l := msg.FW_L & 15
	fw_l_2 := msg.FW_L >> 4
	fw_h := msg.FW_H & 15
	fw_h_2 := msg.FW_H >> 4

	if g.devInfo.Value == nil {
		g.devInfo.Value = &FDGDeviceInfo{}
	}

	if fw_l+fw_l_2+fw_h+fw_h_2 > 0 {
		g.devInfo.Value.DongleVersion = fmt.Sprintf("%d.%d.%d.%d", fw_l_2, fw_l, fw_h_2, fw_h)
		g.devInfo.Value.ConnectType = FDGConncetWireless
	} else {
		g.devInfo.Value.ConnectType = FDGConncetWired
	}

	return nil
}

func (g *Gamepad) handleGamepadConfigRead(msg protocol.MessageGamepadConfigReadCB) error {
	log.Debug().Int("length", len(msg.Data)).Msg("got gamepad configuration data")

	deviceId := int32(0)
	if g.devInfo.Value != nil {
		deviceId = g.devInfo.Value.DeviceId
	}
	cfg, err := config.ConvertGPConfigByByte(msg.Data, deviceId)
	if err != nil {
		preview := msg.Data
		if len(preview) > 32 {
			preview = preview[:32]
		}
		log.Error().Err(err).Bytes("data", preview).Int("length", len(msg.Data)).Int32("device_id", deviceId).Msg("failed to parse gamepad config")
		return fmt.Errorf("convert GP config: %w", err)
	}

	g.currConfig.Value = cfg
	g.currConfig.Broadcast()

	return nil
}

func (g *Gamepad) handleLEDConfigRead(msg protocol.MessageLEDConfigReadCB) error {
	log.Debug().Int("length", len(msg.Data)).Msg("got led configuration data")

	cfg := config.ConvertLEDConfigByByte(msg.Data)

	if g.currConfig.Value != nil {
		g.currConfig.Value.Basic.NewLedConfig = cfg
	}

	g.currLEDConfig.Value = cfg
	g.currLEDConfig.Broadcast()

	return nil
}

func (g *Gamepad) SaveConfig(cfg *config.AllConfigBean) error {
	var buf bytes.Buffer
	config.ConvertByteByGConfig(&buf, cfg)
	configID := g.selectedConfigID()

	log.Info().Int("length", buf.Len()).Uint8("config_id", configID).Msg("saving gamepad configuration")

	if err := g.prot.Send(protocol.CommandSendConfig{
		Data:     buf.Bytes(),
		ConfigID: configID,
	}); err != nil {
		return fmt.Errorf("send config: %w", err)
	}

	g.currConfig.Value = nil

	buf.Reset()

	if cfg.Basic.NewLedConfig != nil {
		if err := g.SaveLEDConfig(cfg.Basic.NewLedConfig); err != nil {
			return fmt.Errorf("save led config: %w", err)
		}
	}

	return nil
}

func (g *Gamepad) SaveLEDConfig(cfg *config.NewLedConfigBean) error {
	var buf bytes.Buffer
	config.ConvertByteByNewLedConfig(&buf, cfg)
	configID := g.selectedConfigID()

	log.Info().Int("length", buf.Len()).Uint8("config_id", configID).Uint8("led_mode", byte(cfg.LedMode)).Msg("saving led configuration")

	if err := g.prot.Send(protocol.CommandSendLEDConfig{
		Data:     buf.Bytes(),
		ConfigID: configID,
	}); err != nil {
		return fmt.Errorf("send config: %w", err)
	}

	g.currLEDConfig.Value = nil

	return nil
}

func getConfigRetry[T any](prot protocol.Protocol, v *utils.CondValue[T], cmd protocol.Command) (*T, error) {
	if v.Value == nil {
		retriesLeft := 3

		for retriesLeft > 0 {
			err := prot.Send(cmd)
			if err != nil {
				return nil, fmt.Errorf("send command: %w", err)
			}

			select {
			case <-v.NotifyChan():
			case <-time.After(2 * time.Second):
				retriesLeft--
				continue
			}

			return v.Value, nil
		}

		return nil, errors.New("device doesn't respond")
	}

	return v.Value, nil
}

func (g *Gamepad) GetConfig() (*config.AllConfigBean, error) {
	if g.devInfo.Value == nil {
		if _, err := g.GetGamepadInfo(); err != nil {
			return nil, fmt.Errorf("get gamepad info: %w", err)
		}
	}

	return getConfigRetry(g.prot, g.currConfig, protocol.CommandReadConfig{ConfigID: g.selectedConfigID()})
}

func (g *Gamepad) GetLEDConfig() (*config.NewLedConfigBean, error) {
	return getConfigRetry(g.prot, g.currLEDConfig, protocol.CommandReadLEDConfig{ConfigID: g.selectedConfigID()})
}

func (g *Gamepad) GetGamepadInfo() (*FDGDeviceInfo, error) {
	info, err := getConfigRetry(g.prot, g.devInfo, protocol.CommandGetDeviceInfo{})
	if err != nil {
		return nil, err
	}
	if info.ConnectMode == "" {
		info.ConnectMode = string(g.protocolMode)
	}

	return info, nil
}
