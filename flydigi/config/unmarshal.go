package config

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"

	"golang.org/x/text/encoding/simplifiedchinese"
)

var (
	keyStrDic = map[int]string{
		0:   "UP",
		1:   "RIGHT",
		2:   "DOWN",
		3:   "LEFT",
		4:   "A",
		5:   "B",
		6:   "SELECT",
		7:   "X",
		8:   "Y",
		9:   "START",
		10:  "LB",
		11:  "RB",
		12:  "LT",
		13:  "RT",
		14:  "THUMBL",
		15:  "THUMBR",
		16:  "C",
		17:  "Z",
		18:  "M1",
		19:  "M2",
		20:  "M3",
		21:  "M4",
		22:  "M5",
		23:  "M6",
		24:  "MENU",
		27:  "HOME",
		28:  "BACK",
		160: "JO",
		161: "JU",
		162: "JUR",
		163: "JR",
		164: "JRD",
		165: "JD",
		166: "JDL",
		167: "JL",
		168: "JLU",
	}

	gamepadKeyPriority = []byte{
		13, 14, 15, 16, 9, 10, 23, 11, 12, 24,
		19, 20, 21, 22, 17, 18, 7, 8, 1, 2,
		3, 4, 5, 6, 32, 32, 32, 32, 32, 32,
		32, 32,
	}
	gamepadKeyListPriority = []byte{
		96, 93, 95, 92, 100, 99, 90, 98, 97, 89,
		88, 87, 86, 85, 84, 83, 73, 72, 77, 80,
		78, 79, 76, 85, 32, 32, 32, 32, 32, 32,
		32, 32,
	}
	gamepadKeyPosType = []KeyPosTypeEnum{
		KeyPosTypeFront,
		KeyPosTypeFront,
		KeyPosTypeFront,
		KeyPosTypeFront,
		KeyPosTypeFront,
		KeyPosTypeFront,
		KeyPosTypeFront,
		KeyPosTypeFront,
		KeyPosTypeFront,
		KeyPosTypeFront,
		KeyPosTypeSide,
		KeyPosTypeSide,
		KeyPosTypeSide,
		KeyPosTypeSide,
		KeyPosTypeFront,
		KeyPosTypeFront,
		KeyPosTypeExt,
		KeyPosTypeExt,
		KeyPosTypeExt,
		KeyPosTypeExt,
		KeyPosTypeExt,
		KeyPosTypeExt,
		KeyPosTypeExt,
		KeyPosTypeExt,
		KeyPosTypeFront,
		KeyPosTypeFront,
		KeyPosTypeFront,
		KeyPosTypeFront,
		KeyPosTypeFront,
		KeyPosTypeFront,
		KeyPosTypeFront,
		KeyPosTypeFront,
	}
	gamepadKeyConfigDic = map[int]GamePadKeyConfig{
		0: GamePadKeyConfig{
			KeyId:      0,
			KeyName:    "UP",
			KeyPosType: KeyPosTypeFront,
			IsShow:     true,
			IsShowPic:  true,
			ShowRes:    "up",
			MapShowRes: "",
		},
		1: GamePadKeyConfig{
			KeyId:      1,
			KeyName:    "RIGHT",
			KeyPosType: KeyPosTypeFront,
			IsShow:     true,
			IsShowPic:  true,
			ShowRes:    "right",
			MapShowRes: "",
		},
		2: GamePadKeyConfig{
			KeyId:      2,
			KeyName:    "DOWN",
			KeyPosType: KeyPosTypeFront,
			IsShow:     true,
			IsShowPic:  true,
			ShowRes:    "down",
			MapShowRes: "",
		},
		3: GamePadKeyConfig{
			KeyId:      3,
			KeyName:    "LEFT",
			KeyPosType: KeyPosTypeFront,
			IsShow:     true,
			IsShowPic:  true,
			ShowRes:    "left",
			MapShowRes: "",
		},
		4: GamePadKeyConfig{
			KeyId:      4,
			KeyName:    "A",
			KeyPosType: KeyPosTypeFront,
			IsShow:     true,
			IsShowPic:  false,
			ShowRes:    "a",
			MapShowRes: "",
		},
		5: GamePadKeyConfig{
			KeyId:      5,
			KeyName:    "B",
			KeyPosType: KeyPosTypeFront,
			IsShow:     true,
			IsShowPic:  false,
			ShowRes:    "b",
			MapShowRes: "",
		},
		6: GamePadKeyConfig{
			KeyId:      6,
			KeyName:    "SELECT",
			KeyPosType: KeyPosTypeFront,
			IsShow:     true,
			IsShowPic:  false,
			ShowRes:    "select",
			MapShowRes: "",
		},
		7: GamePadKeyConfig{
			KeyId:      7,
			KeyName:    "X",
			KeyPosType: KeyPosTypeFront,
			IsShow:     true,
			IsShowPic:  false,
			ShowRes:    "x",
			MapShowRes: "",
		},
		8: GamePadKeyConfig{
			KeyId:      8,
			KeyName:    "Y",
			KeyPosType: KeyPosTypeFront,
			IsShow:     true,
			IsShowPic:  false,
			ShowRes:    "y",
			MapShowRes: "",
		},
		9: GamePadKeyConfig{
			KeyId:      9,
			KeyName:    "START",
			KeyPosType: KeyPosTypeFront,
			IsShow:     true,
			IsShowPic:  false,
			ShowRes:    "start",
			MapShowRes: "",
		},
		10: GamePadKeyConfig{
			KeyId:      10,
			KeyName:    "LB",
			KeyPosType: KeyPosTypeSide,
			IsShow:     true,
			IsShowPic:  false,
			ShowRes:    "lb",
			MapShowRes: "",
		},
		11: GamePadKeyConfig{
			KeyId:      11,
			KeyName:    "RB",
			KeyPosType: KeyPosTypeSide,
			IsShow:     true,
			IsShowPic:  false,
			ShowRes:    "rb",
			MapShowRes: "",
		},
		12: GamePadKeyConfig{
			KeyId:      12,
			KeyName:    "LT",
			KeyPosType: KeyPosTypeSide,
			IsShow:     true,
			IsShowPic:  false,
			ShowRes:    "lt",
			MapShowRes: "",
		},
		13: GamePadKeyConfig{
			KeyId:      13,
			KeyName:    "RT",
			KeyPosType: KeyPosTypeSide,
			IsShow:     true,
			IsShowPic:  false,
			ShowRes:    "rt",
			MapShowRes: "",
		},
		14: GamePadKeyConfig{
			KeyId:      14,
			KeyName:    "THUMBL",
			KeyPosType: KeyPosTypeFront,
			IsShow:     true,
			IsShowPic:  false,
			ShowRes:    "l3",
			MapShowRes: "",
		},
		15: GamePadKeyConfig{
			KeyId:      15,
			KeyName:    "THUMBR",
			KeyPosType: KeyPosTypeFront,
			IsShow:     true,
			IsShowPic:  false,
			ShowRes:    "r3",
			MapShowRes: "",
		},
		16: GamePadKeyConfig{
			KeyId:      16,
			KeyName:    "C",
			KeyPosType: KeyPosTypeExt,
			IsShow:     true,
			IsShowPic:  false,
			ShowRes:    "c",
			MapShowRes: "",
		},
		17: GamePadKeyConfig{
			KeyId:      17,
			KeyName:    "Z",
			KeyPosType: KeyPosTypeExt,
			IsShow:     true,
			IsShowPic:  false,
			ShowRes:    "z",
			MapShowRes: "",
		},
		18: GamePadKeyConfig{
			KeyId:      18,
			KeyName:    "M1",
			KeyPosType: KeyPosTypeExt,
			IsShow:     true,
			IsShowPic:  false,
			ShowRes:    "m1",
			MapShowRes: "",
		},
		19: GamePadKeyConfig{
			KeyId:      19,
			KeyName:    "M2",
			KeyPosType: KeyPosTypeExt,
			IsShow:     true,
			IsShowPic:  false,
			ShowRes:    "m2",
			MapShowRes: "",
		},
		20: GamePadKeyConfig{
			KeyId:      20,
			KeyName:    "M3",
			KeyPosType: KeyPosTypeExt,
			IsShow:     true,
			IsShowPic:  false,
			ShowRes:    "m3",
			MapShowRes: "",
		},
		21: GamePadKeyConfig{
			KeyId:      21,
			KeyName:    "M4",
			KeyPosType: KeyPosTypeExt,
			IsShow:     true,
			IsShowPic:  false,
			ShowRes:    "m4",
			MapShowRes: "",
		},
		22: GamePadKeyConfig{
			KeyId:      22,
			KeyName:    "M5",
			KeyPosType: KeyPosTypeExt,
			IsShow:     true,
			IsShowPic:  false,
			ShowRes:    "m5",
			MapShowRes: "",
		},
		23: GamePadKeyConfig{
			KeyId:      23,
			KeyName:    "M6",
			KeyPosType: KeyPosTypeExt,
			IsShow:     true,
			IsShowPic:  false,
			ShowRes:    "m6",
			MapShowRes: "",
		},
		24: GamePadKeyConfig{
			KeyId:      24,
			KeyName:    "Menu",
			KeyPosType: KeyPosTypeFn,
			IsShow:     false,
			IsShowPic:  false,
			ShowRes:    "menu",
			MapShowRes: "",
		},
		25: GamePadKeyConfig{
			KeyId:      25,
			KeyName:    "unknow",
			KeyPosType: KeyPosTypeFn,
			IsShow:     false,
			IsShowPic:  false,
			ShowRes:    "unknow",
			MapShowRes: "",
		},
		26: GamePadKeyConfig{
			KeyId:      26,
			KeyName:    "unknow",
			KeyPosType: KeyPosTypeFn,
			IsShow:     false,
			IsShowPic:  false,
			ShowRes:    "unknow",
			MapShowRes: "",
		},
		27: GamePadKeyConfig{
			KeyId:      27,
			KeyName:    "Home",
			KeyPosType: KeyPosTypeFn,
			IsShow:     false,
			IsShowPic:  false,
			ShowRes:    "home",
			MapShowRes: "",
		},
		28: GamePadKeyConfig{
			KeyId:      28,
			KeyName:    "Back",
			KeyPosType: KeyPosTypeFn,
			IsShow:     false,
			IsShowPic:  false,
			ShowRes:    "back",
			MapShowRes: "",
		},
		29: GamePadKeyConfig{
			KeyId:      29,
			KeyName:    "unknow",
			KeyPosType: KeyPosTypeFn,
			IsShow:     false,
			IsShowPic:  false,
			ShowRes:    "unknow",
			MapShowRes: "",
		},
		30: GamePadKeyConfig{
			KeyId:      30,
			KeyName:    "unknow",
			KeyPosType: KeyPosTypeFn,
			IsShow:     false,
			IsShowPic:  false,
			ShowRes:    "unknow",
			MapShowRes: "",
		},
		31: GamePadKeyConfig{
			KeyId:      31,
			KeyName:    "unknow",
			KeyPosType: KeyPosTypeFn,
			IsShow:     false,
			IsShowPic:  false,
			ShowRes:    "unknow",
			MapShowRes: "",
		},
		160: GamePadKeyConfig{
			KeyId:      160,
			KeyName:    "JO",
			KeyPosType: KeyPosTypeVirtual,
			IsShow:     false,
			IsShowPic:  false,
			ShowRes:    "圆心",
			MapShowRes: "",
		},
		161: GamePadKeyConfig{
			KeyId:      161,
			KeyName:    "JU",
			KeyPosType: KeyPosTypeVirtual,
			IsShow:     false,
			IsShowPic:  false,
			ShowRes:    "上",
			MapShowRes: "",
		},
		162: GamePadKeyConfig{
			KeyId:      162,
			KeyName:    "JUR",
			KeyPosType: KeyPosTypeVirtual,
			IsShow:     false,
			IsShowPic:  false,
			ShowRes:    "右上",
			MapShowRes: "",
		},
		163: GamePadKeyConfig{
			KeyId:      163,
			KeyName:    "JR",
			KeyPosType: KeyPosTypeVirtual,
			IsShow:     false,
			IsShowPic:  false,
			ShowRes:    "右",
			MapShowRes: "",
		},
		164: GamePadKeyConfig{
			KeyId:      164,
			KeyName:    "JRD",
			KeyPosType: KeyPosTypeVirtual,
			IsShow:     false,
			IsShowPic:  false,
			ShowRes:    "右下",
			MapShowRes: "",
		},
		165: GamePadKeyConfig{
			KeyId:      165,
			KeyName:    "JD",
			KeyPosType: KeyPosTypeVirtual,
			IsShow:     false,
			IsShowPic:  false,
			ShowRes:    "下",
			MapShowRes: "",
		},
		166: GamePadKeyConfig{
			KeyId:      166,
			KeyName:    "JDL",
			KeyPosType: KeyPosTypeVirtual,
			IsShow:     false,
			IsShowPic:  false,
			ShowRes:    "左下",
			MapShowRes: "",
		},
		167: GamePadKeyConfig{
			KeyId:      167,
			KeyName:    "JL",
			KeyPosType: KeyPosTypeVirtual,
			IsShow:     false,
			IsShowPic:  false,
			ShowRes:    "左",
			MapShowRes: "",
		},
		168: GamePadKeyConfig{
			KeyId:      168,
			KeyName:    "JLU",
			KeyPosType: KeyPosTypeVirtual,
			IsShow:     false,
			IsShowPic:  false,
			ShowRes:    "左上",
			MapShowRes: "",
		},
	}
	keyOwnSetsDic = map[int][]bool{
		24: []bool{
			true, true, true, true, true, true, true, true, true, true,
			true, true, true, true, true, true, false, false, true, true,
			true, true, false, false, false, false, false, true, false, false,
			false, false,
		},
		26: []bool{
			true, true, true, true, true, true, true, true, true, true,
			true, true, true, true, true, true, false, false, true, true,
			true, true, false, false, false, false, false, true, false, false,
			false, false,
		},
		29: []bool{
			true, true, true, true, true, true, true, true, true, true,
			true, true, true, true, true, true, false, false, true, true,
			true, true, false, false, false, false, false, true, false, false,
			false, false,
		},
		28: []bool{
			true, true, true, true, true, true, true, true, true, true,
			true, true, true, true, true, true, true, true, true, true,
			true, true, false, false, true, false, false, true, false, false,
			false, false,
		},
		80: []bool{
			true, true, true, true, true, true, true, true, true, true,
			true, true, true, true, true, true, true, true, true, true,
			true, true, false, false, true, false, false, true, false, false,
			false, false,
		},
		81: []bool{
			true, true, true, true, true, true, true, true, true, true,
			true, true, true, true, true, true, true, true, true, true,
			true, true, false, false, true, false, false, true, false, false,
			false, false,
		},
		82: []bool{
			true, true, true, true, true, true, true, true, true, true,
			true, true, true, true, true, true, false, false, true, true,
			false, false, false, false, true, false, false, true, false, false,
			false, false,
		},
		83: []bool{
			true, true, true, true, true, true, true, true, true, true,
			true, true, true, true, true, true, false, false, true, true,
			false, false, false, false, true, false, false, true, false, false,
			false, false,
		},
		84: []bool{
			true, true, true, true, true, true, true, true, true, true,
			true, true, true, true, true, true, false, false, true, true,
			true, true, false, false, true, false, false, true, false, false,
			false, false,
		},
	}
	defaultKeyOwnsSetDic = []bool{
		true, true, true, true, true, true, true, true, true, true,
		true, true, true, true, true, true, false, false, true, true,
		true, true, false, false, false, true, false, false, false, false,
		false, false,
	}
)

const deviceId = 80 // Vader 3 Pro

type configLists struct {
	configVersion, configName, configPackageLength, led, keyMapping, joyMapping,
	triggerMapping, motionMapping, motor, lunpan, macroMapping, triggerMotor, autoTrigger []byte
}

func getConfigListsV2(data []byte) configLists {
	return configLists{
		configVersion:       data[0:2],
		configPackageLength: data[2:3],
		led:                 data[3:13],
		keyMapping:          data[13:109],
		joyMapping:          data[109:123],
		triggerMapping:      data[123:137],
		motionMapping:       data[137:145],
		motor:               data[145:154],
		triggerMotor:        data[154:183],
		lunpan:              data[183:185],
		autoTrigger:         data[185:225],
		macroMapping:        data[225:230],
	}
}

func ConvertGPConfigByByte(data []byte) (*AllConfigBean, error) {
	var cfg AllConfigBean

	if data[0] != 0 || data[1] < 1 {
		return nil, errors.New("default config")
	}
	if data[1] < 2 {
		return nil, errors.New("can't parse config V1")
	}

	lists := getConfigListsV2(data)

	cfg.Version = fmt.Sprintf("%d.%d", lists.configVersion[0], lists.configVersion[1])

	if len(lists.configPackageLength) == 1 {
		cfg.PackageLength = int32(lists.configPackageLength[0] - 1)
	} else {
		cfg.PackageLength = 1
	}

	n, err := simplifiedchinese.GBK.NewDecoder().Bytes([]byte(cfg.Name))
	if err != nil {
		return nil, fmt.Errorf("decode name: %w", err)
	}
	cfg.Name = string(n)
	cfg.Basic = getBasicByListByte(lists.led, lists.motor, lists.lunpan)
	cfg.KeyMapping = getKeyMappingByListByte(lists.keyMapping)
	cfg.JoyMapping = getJoyMappingByListByte(lists.joyMapping)
	cfg.TriggerMapping = getTriggerMappingByListByte(lists.triggerMapping, lists.autoTrigger, lists.triggerMotor)
	cfg.MotionMapping = getMotionMappingByListByte(lists.motionMapping)
	cfg.Macro = getMacroByListByte(lists.macroMapping)

	return &cfg, nil
}

func getBasicByListByte(ledData, motorData, lunpanData []byte) *BasicBean {
	return &BasicBean{
		Led: &LedBean{
			Header: int32(ledData[0]),
			Mode:   int32(ledData[1]),
			Peroid: int32(ledData[2]),
			Light:  int32(ledData[3]),
			RgbColor0: []int32{
				int32(ledData[4]),
				int32(ledData[5]),
				int32(ledData[6]),
			},
			RgbColor1: []int32{
				int32(ledData[7]),
				int32(ledData[8]),
				int32(ledData[9]),
			},
		},
		Motor:         getMotorByListByte(motorData),
		LunPanMapping: getLunpanMappingByListByte(lunpanData),
	}
}

func getMotorByListByte(data []byte) *MotorBean {
	if len(data) != 9 {
		return &MotorBean{}
	}

	bean := MotorBean{
		MainSwitch: data[0] == 0,
		LeftMotor: &MotorNewBean{
			Type:  int32(data[1]),
			Scale: int32(data[4]),
		},
		RightMotor: &MotorNewBean{
			Type:  int32(data[5]),
			Scale: int32(data[8]),
		},
	}

	if data[2] > data[3] {
		bean.LeftMotor.Min = int32(data[3])
		bean.LeftMotor.Max = int32(data[2])
	} else {
		bean.LeftMotor.Min = int32(data[2])
		bean.LeftMotor.Max = int32(data[3])
	}

	if data[6] > data[7] {
		bean.RightMotor.Min = int32(data[7])
		bean.RightMotor.Max = int32(data[6])
	} else {
		bean.RightMotor.Min = int32(data[6])
		bean.RightMotor.Max = int32(data[7])
	}

	return &bean
}

func getLunpanMappingByListByte(data []byte) *LunPanMappingBean {
	if len(data) == 2 {
		return &LunPanMappingBean{
			Type: int32(data[0]),
			Rev:  int32(data[1]),
		}
	}

	return &LunPanMappingBean{}
}

func getKeyMappingByListByte(data []byte) []*KeyMappingBean {
	list := make([]*KeyMappingBean, 0)

	if len(data) > 96 {
		return list
	}

	for i := 0; i < len(data)/3; i++ {
		index := i * 3

		keyStr, ok := keyStrDic[i]
		if !ok {
			keyStr = "unkonw"
		}

		keyId := int32(i)
		readMapKeyId := int32(data[index])
		mapKeyId := convertKeyId(keyId, readMapKeyId, false)

		mapping := KeyMappingBean{
			KeyId:        keyId,
			TurboType:    int32(data[index+1]),
			Turbo:        int32(data[index+2]),
			Key:          keyStr,
			MapKey:       convertKeyString(keyId, readMapKeyId, false),
			IsMap:        keyId != mapKeyId,
			MapData:      strconv.Itoa(int(mapKeyId)),
			Priority:     int32(gamepadKeyPriority[i]),
			ListPriority: int32(gamepadKeyListPriority[i]),
			KeyPosType:   gamepadKeyPosType[i],
			ShowRes:      gamepadKeyConfigDic[i].ShowRes,
			MapShowRes:   gamepadKeyConfigDic[i].MapShowRes,
			IsShow:       gamepadKeyConfigDic[i].IsShow,
			IsShowPic:    gamepadKeyConfigDic[i].IsShowPic,
		}

		if owns, ok := keyOwnSetsDic[deviceId]; ok {
			mapping.IsOwn = owns[i]
		} else {
			mapping.IsOwn = defaultKeyOwnsSetDic[i]
		}

		if mapping.Turbo > 0 {
			mapping.MapType = KeyMapTypeBurst
		} else {
			mapping.MapType = KeyMapTypeGamePad
		}

		list = append(list, &mapping)
	}

	return list
}

func convertKeyId(currentKeyId, targetKeyId int32, specialKeyConvert bool) int32 {
	if targetKeyId >= 32 {
		targetKeyId = currentKeyId
	}

	if specialKeyConvert && targetKeyId >= 18 && targetKeyId <= 21 {
		switch currentKeyId {
		case 16:
			targetKeyId = 14
		case 17:
			targetKeyId = 15
		case 18:
			targetKeyId = 0
		case 19:
			targetKeyId = 2
		case 20:
			targetKeyId = 1
		case 21:
			targetKeyId = 3
		}
	}

	return targetKeyId
}

func convertKeyString(currentKeyId, targetKeyId int32, specialKeyConvert bool) string {
	targetKeyId = convertKeyId(currentKeyId, targetKeyId, specialKeyConvert)

	return keyStrDic[int(targetKeyId)]
}

func maxByte(a, b byte) byte {
	if a > b {
		return a
	}
	return b
}

func getJoyMappingByListByte(data []byte) *JoyMappingBean {
	return &JoyMappingBean{
		LeftJoystic: &JoyStickBean{
			Curve: &Curve{
				Type:     int32(data[0]),
				Zero:     int32(data[1]),
				Point0_X: int32(maxByte(data[1], data[2])),
				Point0_Y: int32(data[3]),
				Point1_X: int32(maxByte(maxByte(data[1], data[2]), data[4])),
				Point1_Y: int32(maxByte(data[3], data[5])),
				End:      int32(data[6]),
			},
		},
		RightJoystic: &JoyStickBean{
			Curve: &Curve{
				Type:     int32(data[7]),
				Zero:     int32(data[8]),
				Point0_X: int32(maxByte(data[9], data[10])),
				Point0_Y: int32(data[9]),
				Point1_X: int32(maxByte(maxByte(data[9], data[10]), data[11])),
				Point1_Y: int32(maxByte(data[9], data[12])),
				End:      int32(data[13]),
			},
		},
	}
}

func getTriggerMappingByListByte(data, listAutoTrigger, listTriggerMotor []byte) *TriggerMappingBean {
	getTrigger := func(idx, tmIdx int) *Trigger {
		t := Trigger{
			Curve: &Curve{
				Type:     int32(data[idx]),
				Zero:     int32(data[idx+1]),
				Point0_X: int32(data[idx+2]),
				Point0_Y: int32(data[idx+3]),
				Point1_X: int32(data[idx+4]),
				Point1_Y: int32(data[idx+5]),
				End:      int32(data[idx+6]),
			},
		}

		if len(listTriggerMotor) == 29 {
			t.TriggerMotor = &TriggerMotor{
				LineGear: &TriggerMotorSet{
					Type:      int32(listTriggerMotor[tmIdx+1]),
					Min:       int32(listTriggerMotor[tmIdx+2]),
					Max:       int32(listTriggerMotor[tmIdx+3]),
					Filter:    int32(listTriggerMotor[tmIdx+4]),
					Vibrlimit: int32(listTriggerMotor[tmIdx+5]),
					Scale:     int32(listTriggerMotor[tmIdx+6]),
					TimeLimit: int32(listTriggerMotor[tmIdx+7]),
				},
				MicrGear: &TriggerMotorSet{
					Type:      int32(listTriggerMotor[tmIdx+8]),
					Min:       int32(listTriggerMotor[tmIdx+9]),
					Max:       int32(listTriggerMotor[tmIdx+10]),
					Filter:    int32(listTriggerMotor[tmIdx+11]),
					Vibrlimit: int32(listTriggerMotor[tmIdx+12]),
					Scale:     int32(listTriggerMotor[tmIdx+13]),
					TimeLimit: int32(listTriggerMotor[tmIdx+14]),
				},
			}
		}

		return &t
	}

	mapping := TriggerMappingBean{
		LeftTrigger:  getTrigger(0, 0),
		RightTrigger: getTrigger(7, 14),
	}

	if len(listTriggerMotor) == 29 {
		mapping.MainSwitch = listTriggerMotor[0] == 0
	}

	if len(listAutoTrigger) > 0 {
		getAutoTriggerByListByte(&mapping, listAutoTrigger)
	}

	return &mapping
}

func getAutoTriggerByListByte(tm *TriggerMappingBean, data []byte) {
	if len(data) < 29 {
		return
	}

	getTrigger := func(idx int) *Autotrigger {
		return &Autotrigger{
			Mode: int32(data[idx]),
			VibrationBind: &Vibrationbind{
				Type:      int32(data[idx+1]),
				MinFilter: int32(data[idx+2]),
				Scale:     int32(data[idx+3]),
				TriggerParams: []int32{
					int32(data[idx+4]),
					int32(data[idx+5]),
					int32(data[idx+6]),
					int32(data[idx+7]),
					int32(data[idx+8]),
				},
			},
			MixedBorder: int32(data[idx+9]),
			MixedParams: []int32{
				int32(data[idx+10]),
				int32(data[idx+11]),
				int32(data[idx+12]),
				int32(data[idx+13]),
				int32(data[idx+14]),
				int32(data[idx+15]),
				int32(data[idx+16]),
				int32(data[idx+17]),
				int32(data[idx+18]),
				int32(data[idx+19]),
			},
		}
	}

	tm.LeftTrigger.AutoTrigger = getTrigger(0)
	tm.RightTrigger.AutoTrigger = getTrigger(20)
}

func getMotionMappingByListByte(data []byte) *MotionMappingBean {
	if data[3] == 0 {
		data[3] = 15
	}
	if data[4] == 0 {
		data[4] = 25
	}
	if data[5] == 0 {
		data[5] = 25
	}

	mapping := MotionMappingBean{
		MapType:       int32(data[0]),
		OpenKeyId:     int32(data[1]),
		OpenMapMethod: int32(data[2]),
		Zero:          int32(data[3]),
		DeadZero:      0,
		Sensity:       float32(data[4]),
		SensityX:      float32(data[4]),
		SensityY:      float32(data[5]),
		Mode:          int32(data[6]),
		OpenKeyExtId:  int32(data[7]),
	}

	if mapping.Zero > 100 {
		mapping.Zero = 100
	}
	if mapping.Sensity < 1 {
		mapping.Sensity = 1
	} else if mapping.Sensity > 127 {
		mapping.Sensity = 127
	}

	return &mapping
}

func getMacroByListByte(data []byte) *MacroGPBean {
	var bean MacroGPBean

	if len(data) <= 3 {
		return &bean
	}

	nums := data[0]
	if nums < 1 || nums > 5 {
		return &bean
	}

	for i := 0; i < int(nums); i++ {
		bean.ListOffset = append(bean.ListOffset, int32(data[i+1]))

		front := 6
		start := front + int(data[i+1]*4)
		end := front + int(data[i+2]*4)

		if i == int(nums-1) {
			end = len(data)
		}

		list := data[start:end]

		if len(list) > 0 {
			oneMacro := getOneMacroByByteArray(list)

			if len(oneMacro.ListStep) > 0 {
				bean.ListMacro = append(bean.ListMacro, oneMacro)
			}
		}
	}

	bean.Nums = int32(len(bean.ListMacro))

	return &bean
}

func getOneMacroByByteArray(data []byte) *OneMacroBeanGP {
	macKeyNum := int32(data[1])

	if len(data) < int(macKeyNum)*4 {
		macKeyNum = int32(len(data)-4) / 4
	}

	bean := OneMacroBeanGP{
		Btn:         int32(data[0]),
		Count_l:     data[1],
		Count_h:     data[2],
		TriggerType: int32(data[3]),
	}

	for i := 0; i < int(macKeyNum); i++ {

		start := 4 + 4*i
		if start+3 >= len(data) {
			break
		}

		intValue := binary.LittleEndian.Uint16(data[start:])

		macroAction := MacroActionBeanGP{
			TriggerTime: int32(intValue),
			Btn:         int32(data[start+2]),
			State:       int32(data[start+3]),
		}

		bean.ListStep = append(bean.ListStep, &macroAction)
	}

	return &bean
}

func ConvertLEDConfigByByte(data []byte) *NewLedConfigBean {
	const ledNum = 15 // This should be 16 but it seems like the data we get is too short
	const ledGroupNum = 10

	bean := NewLedConfigBean{
		Version:     data[:2],
		Type:        data[2],
		Loop_Start:  data[3],
		Loop_End:    data[4],
		Loop_time:   data[5],
		Light_scale: data[6],
		Rgb_num:     data[7],
		LedMode:     LedMode(data[8]),
		Reserve:     data[9:20],
	}

	startIndex := 20

	for i := 0; i < ledNum; i++ {
		var group LedGroup

		for j := 0; j < ledGroupNum; j++ {
			pos := j * 3

			group.Units = append(group.Units, &LedUnit{
				R: data[startIndex+pos],
				G: data[startIndex+pos+1],
				B: data[startIndex+pos+2],
			})
		}

		bean.LedGroups = append(bean.LedGroups, &group)

		startIndex += ledGroupNum * 3
	}

	return &bean
}
