package flydigi

import (
	"bytes"
	"errors"
	"flydigi-linux/flydigi/config"
	"flydigi-linux/flydigi/protocol"
	"flydigi-linux/flydigi/protocol/dinput"
	"flydigi-linux/utils"
	"fmt"
	"io"
	"net"
	"os"
	"sync"

	"github.com/karalabe/usb"
	"github.com/rs/zerolog/log"
)

type deviceMode int

const (
	deviceModeXInput deviceMode = iota + 1
	deviceModeDInput
)

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
	prot protocol.Protocol

	devInfo FDGDeviceInfo

	currConfig    *utils.CondValue[config.AllConfigBean]
	currLEDConfig *utils.CondValue[config.NewLedConfigBean]

	configID byte
}

func OpenGamepad() (*Gamepad, error) {
	var prot protocol.Protocol

	dev, err := openDeviceDInput()
	if err == nil {
		prot = dinput.Open(dev)
	} else {
		if err != os.ErrNotExist {
			return nil, fmt.Errorf("open dinput device: %w", err)
		}

		dev, err = openDeviceXInput()
		if err != nil {
			return nil, fmt.Errorf("open xinput device: %w", err)
		}

		// prot = protocol.OpenXInput(dev)
	}

	gamepad := &Gamepad{
		prot:          prot,
		currConfig:    utils.NewCondValue[config.AllConfigBean](&sync.Mutex{}),
		currLEDConfig: utils.NewCondValue[config.NewLedConfigBean](&sync.Mutex{}),
	}
	go gamepad.readLoop()

	return gamepad, nil
}

func openDeviceDInput() (io.ReadWriteCloser, error) {
	devs, err := usb.EnumerateHid(0x04b4, 0x2412)
	if err != nil {
		return nil, fmt.Errorf("enumerate devices: %w", err)
	}

	if len(devs) == 0 {
		return nil, os.ErrNotExist
	}

	for _, d := range devs {
		if d.Interface == 2 {
			return d.Open()
		}
	}

	return nil, os.ErrNotExist
}

func openDeviceXInput() (io.ReadWriteCloser, error) {
	devs, err := usb.EnumerateRaw(0x045e, 0x028e)
	if err != nil {
		return nil, fmt.Errorf("enumerate devices: %w", err)
	}

	if len(devs) == 0 {
		return nil, os.ErrNotExist
	}

	return devs[0].Open()
}

func (g *Gamepad) Close() error {
	return g.prot.Close()
}

func (g *Gamepad) readLoop() {
	for msg := range g.prot.Messages() {
		if err := g.handleMessage(msg); err != nil {
			log.Err(err).Msg("failed to handle usb data")
		}
	}
}

func (g *Gamepad) handleMessage(msg protocol.Message) error {
	switch msg := msg.(type) {
	case protocol.MessageGamePadInfo:
		return g.handleDeviceInfo(msg)

	case protocol.MessageDongleInfo:
		return g.handleDongleInfo(msg)

	case protocol.MessageGamepadConfigReadCB:
		return g.handleGamepadConfigRead(msg)

	case protocol.MessageLEDConfigReadCB:
		return g.handleLEDConfigRead(msg)

	default:
		return errors.New("unknown message type")
	}
}

func (g *Gamepad) handleDeviceInfo(msg protocol.MessageGamePadInfo) error {
	g.devInfo = FDGDeviceInfo{}

	g.devInfo.DeviceId = int32(msg.DeviceID)
	// if (!GameHandleListDic.gameHandleDic.ContainsKey((int)deviceId))
	// {
	// 	return;
	// }

	// currentDeviceInfo.GameHadleName = GameHandleListDic.gameHandleDic[(int)deviceId].GameHadleName;
	// currentDeviceInfo.FirmwareName = GameHandleListDic.gameHandleDic[(int)deviceId].FirmwareName;

	g.devInfo.DeviceMac = net.HardwareAddr(msg.DeviceMac).String()

	fw_l := msg.FW_L & 15
	fw_l_2 := msg.FW_L >> 4
	fw_h := msg.FW_H & 15
	fw_h_2 := msg.FW_H >> 4

	g.devInfo.FirmwareVersionCode = int32(fw_h_2)*1000 + int32(fw_h)*100 + int32(fw_l_2)*10 + int32(fw_l)
	g.devInfo.FirmwareVersion = fmt.Sprintf("%d.%d.%d.%d", fw_h_2, fw_h, fw_l_2, fw_l)

	battery := msg.Battery
	const apex2MinBY = 98
	const apex2MaxBY = 114

	if battery < apex2MinBY {
		battery = apex2MinBY
	} else if battery > apex2MaxBY {
		battery = apex2MaxBY
	}

	batteryPercent := int(100 * float32(battery-apex2MinBY) / float32(apex2MaxBY-apex2MinBY))
	g.devInfo.BatteryPercent = int32(batteryPercent)

	switch msg.MotionSensorType {
	case 1:
		g.devInfo.MotionSensorType = "ST"
	case 2:
		g.devInfo.MotionSensorType = "QST"
	}

	if msg.CPUType > 0 {
		g.devInfo.CpuType = "wch"
	} else {
		g.devInfo.CpuType = "nordic"
	}

	if fw_h_2 >= 6 && fw_h >= 1 {
		g.devInfo.CpuType = "wch"
	}

	if g.devInfo.CpuType == "wch" {
		if msg.ConnectionType == 1 {
			g.devInfo.ConnectType = FDGConncetWired
			g.devInfo.CpuName = "ch573"
		} else {
			g.devInfo.ConnectType = FDGConncetWireless
			g.devInfo.CpuName = "ch571"
			g.prot.Send(protocol.CommandGetDongleVersion{})
		}
	}

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

	switch g.devInfo.DeviceId {
	case 19:
		g.devInfo.FirmwareName = "apex2"

	case 20, 21:
		g.devInfo.FirmwareName = "f1"
		if g.devInfo.CpuType == "wch" {
			g.devInfo.FirmwareName = "f1wch"
		}

	case 22, 23:
		g.devInfo.FirmwareName = "f1p"

	case 24:
		g.devInfo.FirmwareName = "k1"

	case 25:
		g.devInfo.FirmwareName = "fp1"
	}

	return nil
}

func (g *Gamepad) handleDongleInfo(msg protocol.MessageDongleInfo) error {
	fw_l := msg.FW_L & 15
	fw_l_2 := msg.FW_L >> 4
	fw_h := msg.FW_H & 15
	fw_h_2 := msg.FW_H >> 4

	if fw_l+fw_l_2+fw_h+fw_h_2 > 0 {
		g.devInfo.DongleVersion = fmt.Sprintf("%d.%d.%d.%d", fw_l_2, fw_l, fw_h_2, fw_h)
		g.devInfo.ConnectType = FDGConncetWireless
	} else {
		g.devInfo.ConnectType = FDGConncetWired
	}

	return nil
}

func (g *Gamepad) handleGamepadConfigRead(msg protocol.MessageGamepadConfigReadCB) error {
	cfg, err := config.ConvertGPConfigByByte(msg.Data)
	if err != nil {
		return fmt.Errorf("convert GP config: %w", err)
	}

	g.currConfig.Value = cfg
	g.currConfig.Broadcast()

	return nil
}

func (g *Gamepad) handleLEDConfigRead(msg protocol.MessageLEDConfigReadCB) error {
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

	if err := g.prot.Send(protocol.CommandSendConfig{
		Data:     buf.Bytes(),
		ConfigID: g.configID,
	}); err != nil {
		return fmt.Errorf("send config: %w", err)
	}

	buf.Reset()
	config.ConvertByteByNewLedConfig(&buf, cfg.Basic.NewLedConfig)

	if err := g.prot.Send(protocol.CommandSendLEDConfig{
		Data:     buf.Bytes(),
		ConfigID: g.configID,
	}); err != nil {
		return fmt.Errorf("send config: %w", err)
	}

	return nil
}

func (g *Gamepad) GetConfig() (*config.AllConfigBean, error) {
	if g.currConfig.Value != nil {
		return g.currConfig.Value, nil
	}

	err := g.prot.Send(protocol.CommandReadConfig{ConfigID: g.configID})
	if err != nil {
		return nil, fmt.Errorf("send command: %w", err)
	}

	<-g.currConfig.NotifyChan()
	return g.currConfig.Value, nil
}

func (g *Gamepad) GetLEDConfig() (*config.NewLedConfigBean, error) {
	if g.currLEDConfig.Value != nil {
		return g.currLEDConfig.Value, nil
	}

	err := g.prot.Send(protocol.CommandReadLEDConfig{ConfigID: g.configID})
	if err != nil {
		return nil, fmt.Errorf("send command: %w", err)
	}

	<-g.currLEDConfig.NotifyChan()
	return g.currLEDConfig.Value, nil
}
