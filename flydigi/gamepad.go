package flydigi

import (
	"bytes"
	"errors"
	"flydigi-linux/flydigi/config"
	"flydigi-linux/utils"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/sstallion/go-hid"
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

const (
	packageLength    = 52
	ledPackageLength = 49
	maxParcelLength  = 10
)

type CommandNumber byte

const (
	CommandGetDongleVersion       CommandNumber = 17
	CommandReadConfig             CommandNumber = 235
	CommandGetDeviceInfoInAndroid CommandNumber = 236
	CommandReadLEDConfig          CommandNumber = 229
)

type commandCallbackFunc func(data []byte)

type configTrasmission struct {
	chunks [][]byte
	ackch  chan int
}

type Gamepad struct {
	dev *hid.Device

	transLock sync.Mutex

	devInfo FDGDeviceInfo

	currConfig    *utils.CondValue[config.AllConfigBean]
	currLEDConfig *utils.CondValue[config.NewLedConfigBean]

	configID byte

	configData    [packageLength * 10]byte
	ledConfigData [ledPackageLength * 10]byte

	sendingConfig *configTrasmission
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

	gamepad := &Gamepad{
		dev:           dev,
		currConfig:    utils.NewCondValue[config.AllConfigBean](&sync.Mutex{}),
		currLEDConfig: utils.NewCondValue[config.NewLedConfigBean](&sync.Mutex{}),
	}
	go gamepad.readLoop()

	return gamepad, nil
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

func (g *Gamepad) resolveUsbData(p []byte) error {
	// log.Debug().Int("len", len(p)).Msg("got usb data")

	if p[15] == 235 {
		log.Debug().Str("handler", "GamepadConfigReadCB").Msg("got usb response")
		return g.handleGamepadConfigRead(p)
	}

	if p[15] == 229 {
		log.Debug().Str("handler", "LEDConfigReadCB").Msg("got usb response")
		return g.handleLEDConfigRead(p)
	}

	if p[0] == 4 && p[1] == 17 {
		log.Debug().Str("handler", "HandleDongleInfo").Msg("got usb response")
		return g.handleDongleInfo(p)
	}

	if p[15] == 234 || p[15] == 231 || p[15] == 51 {
		log.Debug().Str("handler", "WriteGamepadConfigCBK").Msg("got usb response")

		if g.sendingConfig != nil {
			select {
			case g.sendingConfig.ackch <- int(p[3]):
			default:
				log.Debug().Msg("dropped config write ack")
			}
		} else {
			log.Debug().Msg("missed config write ack")
		}

		return nil
	}

	if p[15] == 236 {
		//TODO: Check acks
		log.Debug().Str("handler", "GamePadInfo").Msg("got usb response")
		return g.handleDeviceInfo(p)
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

func (g *Gamepad) handleDeviceInfo(data []byte) error {
	g.devInfo = FDGDeviceInfo{}

	deviceId := data[3]

	g.devInfo.DeviceId = int32(deviceId)
	// if (!GameHandleListDic.gameHandleDic.ContainsKey((int)deviceId))
	// {
	// 	return;
	// }

	// currentDeviceInfo.GameHadleName = GameHandleListDic.gameHandleDic[(int)deviceId].GameHadleName;
	// currentDeviceInfo.FirmwareName = GameHandleListDic.gameHandleDic[(int)deviceId].FirmwareName;

	deviceMac := data[5:9]
	g.devInfo.DeviceMac = net.HardwareAddr(deviceMac).String()

	fw_l := data[9] & 15
	fw_l_2 := data[9] >> 4
	fw_h := data[10] & 15
	fw_h_2 := data[10] >> 4

	g.devInfo.FirmwareVersionCode = int32(fw_h_2)*1000 + int32(fw_h)*100 + int32(fw_l_2)*10 + int32(fw_l)
	g.devInfo.FirmwareVersion = fmt.Sprintf("%d.%d.%d.%d", fw_h_2, fw_h, fw_l_2, fw_l)

	battery := data[11]
	const apex2MinBY = 98
	const apex2MaxBY = 114

	if battery < apex2MinBY {
		battery = apex2MinBY
	} else if battery > apex2MaxBY {
		battery = apex2MaxBY
	}

	batteryPercent := int(100 * float32(battery-apex2MinBY) / float32(apex2MaxBY-apex2MinBY))
	g.devInfo.BatteryPercent = int32(batteryPercent)

	switch data[14] {
	case 1:
		g.devInfo.MotionSensorType = "ST"
	case 2:
		g.devInfo.MotionSensorType = "QST"
	}

	if data[12] > 0 {
		g.devInfo.CpuType = "wch"
	} else {
		g.devInfo.CpuType = "nordic"
	}

	if fw_h_2 >= 6 && fw_h >= 1 {
		g.devInfo.CpuType = "wch"
	}

	if g.devInfo.CpuType == "wch" {
		if data[13] == 1 {
			g.devInfo.ConnectType = FDGConncetWired
			g.devInfo.CpuName = "ch573"
		} else {
			g.devInfo.ConnectType = FDGConncetWireless
			g.devInfo.CpuName = "ch571"
			g.SendCommand(CommandGetDongleVersion)
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

func (g *Gamepad) handleDongleInfo(data []byte) error {
	fw_l := data[2] & 15
	fw_l_2 := data[2] >> 4
	fw_h := data[3] & 15
	fw_h_2 := data[3] >> 4

	if fw_l+fw_l_2+fw_h+fw_h_2 > 0 {
		g.devInfo.DongleVersion = fmt.Sprintf("%d.%d.%d.%d", fw_l_2, fw_l, fw_h_2, fw_h)
		g.devInfo.ConnectType = FDGConncetWireless
	} else {
		g.devInfo.ConnectType = FDGConncetWired
	}

	return nil
}

func (g *Gamepad) handleGamepadConfigRead(p []byte) error {
	packageIndex := int32(p[3])
	if packageIndex > packageLength {
		return nil
	}

	data := p[5:15]

	copy(g.configData[packageIndex*10:], data)

	if packageIndex != 52 {
		return nil
	}

	cfg, err := config.ConvertGPConfigByByte(g.configData[:])
	if err != nil {
		return fmt.Errorf("convert GP config: %w", err)
	}

	g.currConfig.Value = cfg
	g.currConfig.Broadcast()

	return nil
}

func (g *Gamepad) handleLEDConfigRead(p []byte) error {
	packageIndex := int32(p[3])
	if packageIndex > packageLength {
		return nil
	}

	data := p[5:15]

	copy(g.ledConfigData[packageIndex*10:], data)

	if packageIndex != ledPackageLength {
		return nil
	}

	cfg := config.ConvertLEDConfigByByte(g.ledConfigData[:])

	if g.currConfig.Value != nil {
		g.currConfig.Value.Basic.NewLedConfig = cfg
	}

	g.currLEDConfig.Value = cfg
	g.currLEDConfig.Broadcast()

	return nil
}

func (g *Gamepad) sendConfig(data [][]byte) error {
	g.transLock.Lock()
	defer g.transLock.Unlock()

	g.sendingConfig = &configTrasmission{
		chunks: data,
		ackch:  make(chan int),
	}
	defer func() {
		g.sendingConfig = nil
	}()

	for i, chunk := range data {
		retriesLeft := 3
		success := false

		for !success && retriesLeft > 0 {
			retriesLeft--

			// p := append([]byte{2}, chunk...)

			_, err := g.dev.Write(chunk)
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

func (g *Gamepad) SaveConfig(cfg *config.AllConfigBean) error {
	var buf bytes.Buffer
	config.ConvertByteByGConfig(&buf, cfg)

	chunks := getConfigDataParcels(buf.Bytes(), g.configID)

	if err := g.sendConfig(chunks); err != nil {
		return fmt.Errorf("send config: %w", err)
	}

	buf.Reset()
	config.ConvertByteByNewLedConfig(&buf, cfg.Basic.NewLedConfig)

	chunks = getLEDConfigDataParcels(buf.Bytes(), g.configID)

	if err := g.sendConfig(chunks); err != nil {
		return fmt.Errorf("send led config: %w", err)
	}

	return nil
}

func readCondRetry[T any](g *Gamepad, cmd CommandNumber, cond *utils.CondValue[T]) (*T, error) {
	g.transLock.Lock()
	defer g.transLock.Unlock()

	if cond.Value != nil {
		return cond.Value, nil
	}

	retriesLeft := 3

	for retriesLeft > 0 {
		err := g.SendCommand(cmd, g.configID)
		if err != nil {
			return nil, fmt.Errorf("send command: %w", err)
		}

		select {
		case <-cond.NotifyChan():
			time.Sleep(200 * time.Millisecond) // Wait for transmission to end
			return cond.Value, nil

		case <-time.After(3 * time.Second):
			log.Warn().Msg("retrying")
			retriesLeft--
		}
	}

	return nil, errors.New("device did not respond")
}

func (g *Gamepad) GetConfig() (*config.AllConfigBean, error) {
	return readCondRetry(g, CommandReadConfig, g.currConfig)
}

func (g *Gamepad) GetLEDConfig() (*config.NewLedConfigBean, error) {
	return readCondRetry(g, CommandReadLEDConfig, g.currLEDConfig)
}

func (g *Gamepad) SendCommand(cmd CommandNumber, args ...byte) error {
	log.Debug().Int("cmd", int(cmd)).Bytes("args", args).Msg("sending command")

	buf := make([]byte, 12)
	buf[0] = 5
	buf[1] = byte(cmd)
	copy(buf[2:], args)

	_, err := g.dev.Write(buf)
	return err
}
