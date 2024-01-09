package config

import (
	"bytes"
	"encoding/binary"
	"strconv"
	"strings"

	"github.com/pipe01/flydigictl/pkg/utils"
)

func ConvertByteByGConfig(buf *bytes.Buffer, config *AllConfigBean) {
	writeConfigVersionList(buf, config.Version)
	buf.WriteByte(byte(config.PackageLength) + 1)
	writeLedList(buf, config.Basic.Led)
	writeKeyMappingList(buf, config.KeyMapping)
	writeJoyMappingList(buf, config.JoyMapping)
	writeTriggerMappingList(buf, config.TriggerMapping)
	writeMotionMappingList(buf, config.MotionMapping)
	writeMotorList(buf, config.Basic.Motor)
	writeTriggerMotorList(buf, config.TriggerMapping)
	writeLunpanList(buf, config.Basic.LunPanMapping)
	writeAutoTriggerList(buf, config.TriggerMapping)
	buf.Write(utils.Repeat(byte(255), 5))
	writeMacroList(buf, config.Macro)
	buf.Write(utils.Repeat(byte(255), 2))
}

func writeConfigVersionList(w *bytes.Buffer, vers string) {
	sl, sh, ok := strings.Cut(vers, ".")
	if ok {
		l, _ := strconv.ParseUint(sl, 10, 8)
		h, _ := strconv.ParseUint(sh, 10, 8)

		w.Write([]byte{byte(l), byte(h)})
	}
}

func writeLedList(w *bytes.Buffer, led *LedBean) {
	w.Write([]byte{
		32,
		byte(led.Mode),
		byte(led.Peroid),
		byte(led.Light),
		byte(led.RgbColor0[0]),
		byte(led.RgbColor0[1]),
		byte(led.RgbColor0[2]),
		byte(led.RgbColor1[0]),
		byte(led.RgbColor1[1]),
		byte(led.RgbColor1[2]),
	})
}

func writeKeyMappingList(bw *bytes.Buffer, keyMappingList []*KeyMappingBean) {
	for i, m := range keyMappingList {
		if m.IsTriggerNative && m.MapType != KeyMapTypeGamePad {
			bw.WriteByte(254)
			bw.WriteByte(byte(m.TurboType))
			bw.WriteByte(byte(m.Turbo))
		} else {
			mapKeyId := m.KeyId

			turbo := byte(m.Turbo)
			if m.MapType != KeyMapTypeBurst {
				turbo = 0
			}

			if m.MapType == KeyMapTypeGamePad || m.MapType == KeyMapTypeBurst {
				v, _ := strconv.ParseInt(m.MapData, 10, 16)
				mapKeyId = int32(v)

				if mapKeyId == int32(i) {
					bw.WriteByte(255)
				} else {
					bw.WriteByte(byte(mapKeyId))
				}

				bw.WriteByte(byte(m.TurboType))
				bw.WriteByte(turbo)
			} else {
				bw.WriteByte(254)
				bw.WriteByte(0)
				bw.WriteByte(turbo)
			}
		}
	}
}

func writeCurve(w *bytes.Buffer, c *Curve) {
	w.Write([]byte{
		byte(c.Type),
		byte(c.Zero),
		byte(c.Point0_X),
		byte(c.Point0_Y),
		byte(c.Point1_X),
		byte(c.Point1_Y),
		byte(c.End),
	})
}

func writeJoyMappingList(w *bytes.Buffer, joy *JoyMappingBean) {
	writeCurve(w, joy.LeftJoystic.Curve)
	writeCurve(w, joy.RightJoystic.Curve)
}

func writeTriggerMappingList(w *bytes.Buffer, trigger *TriggerMappingBean) {
	writeCurve(w, trigger.LeftTrigger.Curve)
	writeCurve(w, trigger.RightTrigger.Curve)
}

func writeMotionMappingList(bw *bytes.Buffer, motion *MotionMappingBean) {
	bw.WriteByte(byte(motion.MapType))
	bw.WriteByte(byte(motion.OpenKeyId))
	bw.WriteByte(byte(motion.OpenMapMethod))

	if motion.MapType == 3 {
		bw.WriteByte(byte(motion.DeadZero))
	} else {
		bw.WriteByte(byte(motion.Zero))
	}

	bw.WriteByte(byte(motion.SensityX))
	bw.WriteByte(byte(motion.SensityY))
	bw.WriteByte(byte(motion.Mode))
	bw.WriteByte(byte(motion.OpenKeyExtId))
}

func writeMotorList(w *bytes.Buffer, motor *MotorBean) {
	mainSwitch := byte(1)
	if motor.MainSwitch {
		mainSwitch = 0
	}

	w.Write([]byte{
		mainSwitch,
		byte(motor.LeftMotor.Type),
		byte(motor.LeftMotor.Min),
		byte(motor.LeftMotor.Max),
		byte(motor.LeftMotor.Scale),
		byte(motor.RightMotor.Type),
		byte(motor.RightMotor.Min),
		byte(motor.RightMotor.Max),
		byte(motor.RightMotor.Scale),
	})
}

func writeTriggerMotorSet(w *bytes.Buffer, set *TriggerMotorSet) {
	w.Write([]byte{
		byte(set.Type),
		byte(set.Min),
		byte(set.Max),
		byte(set.Filter),
		byte(set.Vibrlimit),
		byte(set.Scale),
		byte(set.TimeLimit),
	})
}

func writeTriggerMotorList(w *bytes.Buffer, trigger *TriggerMappingBean) {
	mainSwitch := byte(1)
	if trigger.MainSwitch {
		mainSwitch = 0
	}

	w.WriteByte(mainSwitch)
	writeTriggerMotorSet(w, trigger.LeftTrigger.TriggerMotor.LineGear)
	writeTriggerMotorSet(w, trigger.LeftTrigger.TriggerMotor.MicrGear)
	writeTriggerMotorSet(w, trigger.RightTrigger.TriggerMotor.LineGear)
	writeTriggerMotorSet(w, trigger.RightTrigger.TriggerMotor.MicrGear)
}

func writeLunpanList(w *bytes.Buffer, lunpan *LunPanMappingBean) {
	w.Write([]byte{
		byte(lunpan.Type),
		byte(lunpan.Rev),
	})
}

func writeAutoTrigger(bw *bytes.Buffer, trigger *Autotrigger) {
	bw.Write([]byte{
		byte(trigger.Mode),
		byte(trigger.VibrationBind.Type),
		byte(trigger.VibrationBind.MinFilter),
		byte(trigger.VibrationBind.Scale),
	})

	for _, p := range trigger.VibrationBind.TriggerParams {
		bw.WriteByte(byte(p))
	}

	bw.WriteByte(byte(trigger.MixedBorder))

	for _, p := range trigger.MixedParams {
		bw.WriteByte(byte(p))
	}
}

func writeAutoTriggerList(bw *bytes.Buffer, trigger *TriggerMappingBean) {
	writeAutoTrigger(bw, trigger.LeftTrigger.AutoTrigger)
	writeAutoTrigger(bw, trigger.RightTrigger.AutoTrigger)
}

func writeMacroList(w *bytes.Buffer, macro *MacroGPBean) {
	var data bytes.Buffer

	data.WriteByte(byte(len(macro.ListMacro)))
	data.Write(utils.Repeat(byte(0), 5))

	index := 0
	for i := 0; i < len(macro.ListMacro)-1; i++ {
		index += len(macro.ListMacro[i].ListStep) + 1

		data.Bytes()[i+2] = byte(index)
	}

	for _, m := range macro.ListMacro {
		data.WriteByte(byte(m.Btn))
		data.WriteByte(byte(m.Count_l))
		data.WriteByte(byte(m.Count_h))
		data.WriteByte(byte(m.TriggerType))

		for _, s := range m.ListStep {
			binary.Write(&data, binary.LittleEndian, int16(s.TriggerTime/10))
			data.WriteByte(byte(s.Btn))
			data.WriteByte(byte(s.State))
		}
	}

	data.WriteTo(w)
}

func ConvertByteByNewLedConfig(buf *bytes.Buffer, cfg *NewLedConfigBean) {
	buf.Write(cfg.Version)
	buf.WriteByte(cfg.Type)
	buf.WriteByte(cfg.Loop_Start)
	buf.WriteByte(cfg.Loop_End)
	buf.WriteByte(cfg.Loop_time)
	buf.WriteByte(cfg.Light_scale)
	buf.WriteByte(cfg.Rgb_num)
	buf.WriteByte(byte(cfg.LedMode))
	buf.Write(utils.Repeat(byte(255), 11))

	for _, g := range cfg.LedGroups {
		for _, u := range g.Units {
			buf.WriteByte(u.R)
			buf.WriteByte(u.G)
			buf.WriteByte(u.B)
		}
	}
}
