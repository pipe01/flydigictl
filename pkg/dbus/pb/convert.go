package pb

import (
	"strconv"

	"github.com/pipe01/flydigictl/pkg/flydigi/config"
	"golang.org/x/image/colornames"
)

func ColorFromRGB(r, g, b byte) *Color {
	return &Color{Rgb: int32(b) | (int32(g) << 8) | (int32(r) << 16)}
}

func ColorFromHex(hex string) (*Color, bool) {
	if hex[0] == '#' {
		hex = hex[1:]
	}
	if len(hex) != 6 {
		return nil, false
	}

	n, err := strconv.ParseUint(hex, 16, 32)
	if err != nil || n > 0xFFFFFF {
		return nil, false
	}

	return &Color{Rgb: int32(n)}, true
}

func ColorFromName(name string) (*Color, bool) {
	rgba, ok := colornames.Map[name]
	if !ok {
		return nil, false
	}

	return ColorFromRGB(rgba.R, rgba.G, rgba.B), true
}

func (c *Color) RGB() (r, g, b byte) {
	return byte(c.Rgb >> 16), byte(c.Rgb >> 8), byte(c.Rgb)
}

func (c *Color) LedUnit() *config.LedUnit {
	r, g, b := c.RGB()
	return &config.LedUnit{r, g, b}
}

func ConvertGamepadConfiguration(bean *config.AllConfigBean) *GamepadConfiguration {
	return &GamepadConfiguration{
		LeftJoystick:  ConvertJoystickConfiguration(bean.JoyMapping.LeftJoystic),
		RightJoystick: ConvertJoystickConfiguration(bean.JoyMapping.RightJoystic),
		LeftTrigger:   ConvertTriggerConfiguration(bean.TriggerMapping.LeftTrigger),
		RightTrigger:  ConvertTriggerConfiguration(bean.TriggerMapping.RightTrigger),
	}
}

func (c *GamepadConfiguration) ApplyTo(bean *config.AllConfigBean) {
	c.LeftJoystick.ApplyTo(bean.JoyMapping.LeftJoystic)
	c.RightJoystick.ApplyTo(bean.JoyMapping.RightJoystic)
	c.LeftTrigger.ApplyTo(bean.TriggerMapping.LeftTrigger)
	c.RightTrigger.ApplyTo(bean.TriggerMapping.RightTrigger)
}

func ConvertJoystickConfiguration(bean *config.JoyStickBean) *JoystickConfiguration {
	return &JoystickConfiguration{
		Deadzone: bean.Curve.Zero,
	}
}

func (c *JoystickConfiguration) ApplyTo(bean *config.JoyStickBean) {
	bean.Curve.Zero = c.Deadzone
}

func ConvertTriggerConfiguration(bean *config.Trigger) *TriggerConfiguration {
	var mode isTriggerConfiguration_Mode

	switch bean.AutoTrigger.Mode {
	case 0:
		mode = &TriggerConfiguration_Default{}
	case 1:
		mode = &TriggerConfiguration_Race{
			Race: &TriggerRace{
				InitialPos: bean.AutoTrigger.MixedParams[0],
				Pressure:   bean.AutoTrigger.MixedParams[1],
			},
		}
	case 2:
		mode = &TriggerConfiguration_Recoil{
			Recoil: &TriggerRecoil{
				InitialPos:      bean.AutoTrigger.MixedParams[0],
				InitialStrength: bean.AutoTrigger.MixedParams[1],
				Intensity:       bean.AutoTrigger.MixedParams[2],
				Frequency:       bean.AutoTrigger.MixedParams[3],
				InputAfter:      bean.AutoTrigger.MixedParams[4] == 1,
			},
		}
	case 3:
		mode = &TriggerConfiguration_Sniper{
			Sniper: &TriggerSniper{
				InitialPos: bean.AutoTrigger.MixedParams[0],
				Length:     bean.AutoTrigger.MixedParams[1],
				Pressure:   bean.AutoTrigger.MixedParams[2],
				InputAfter: bean.AutoTrigger.MixedParams[4] == 1,
			},
		}
	case 4:
		mode = &TriggerConfiguration_Lock{
			Lock: &TriggerLock{
				InitialPos: bean.AutoTrigger.MixedParams[0],
			},
		}
	case 5:
		mode = &TriggerConfiguration_Vibration{
			Vibration: &TriggerVibration{
				Coefficient: bean.AutoTrigger.VibrationBind.Scale,
				Threshold:   bean.AutoTrigger.VibrationBind.MinFilter,
				TravelRange: bean.AutoTrigger.VibrationBind.TriggerParams[0],
				Frequency:   bean.AutoTrigger.VibrationBind.TriggerParams[3],
			},
		}
	}

	return &TriggerConfiguration{
		Mode: mode,
	}
}

func (c *TriggerConfiguration) ApplyTo(bean *config.Trigger) {
	if bean == nil {
		return
	}

	switch mode := c.Mode.(type) {
	case *TriggerConfiguration_Default:
		bean.ApplyDefaultTrigger()
	case *TriggerConfiguration_Race:
		bean.ApplyRaceTrigger(mode.Race.InitialPos, mode.Race.Pressure)
	case *TriggerConfiguration_Recoil:
		outputInt := int32(0)
		if mode.Recoil.InputAfter {
			outputInt = 1
		}
		bean.ApplyRecoilTrigger(mode.Recoil.InitialPos, mode.Recoil.InitialStrength, mode.Recoil.Intensity, mode.Recoil.Frequency, outputInt)
	case *TriggerConfiguration_Sniper:
		outputInt := int32(0)
		if mode.Sniper.InputAfter {
			outputInt = 1
		}
		bean.ApplySniperTrigger(mode.Sniper.InitialPos, mode.Sniper.Length, mode.Sniper.Pressure, outputInt)
	case *TriggerConfiguration_Lock:
		bean.ApplyLockTrigger(mode.Lock.InitialPos)
	case *TriggerConfiguration_Vibration:
		bean.ApplyVibrationTrigger(mode.Vibration.Coefficient, mode.Vibration.Threshold, mode.Vibration.TravelRange, mode.Vibration.Frequency)
	}
}

func ConvertLEDConfiguration(bean *config.NewLedConfigBean) *LedsConfiguration {
	var leds isLedsConfiguration_Leds

	switch bean.LedMode {
	case config.LedModeOff:
		leds = &LedsConfiguration_Off{}

	case config.LedModeSteady:
		unit := bean.LedGroups[0].Units[0]

		leds = &LedsConfiguration_Steady{
			Steady: &LedsSteady{
				Color: ColorFromRGB(unit.R, unit.G, unit.B),
			},
		}

	case config.LedModeStreamlined:
		leds = &LedsConfiguration_Streamlined{
			Streamlined: &LedsStreamlined{
				Speed: float32(100-bean.Loop_time) / 100,
			},
		}
	case config.LedModeGradient:
		colors := make([]*Color, 0)

		if bean.Rgb_num > 0 && len(bean.LedGroups) > 0 {
			group := bean.LedGroups[0]

			for i := 0; i < int(bean.Rgb_num) && i < len(group.Units); i++ {
				unit := group.Units[i]
				if unit == nil {
					continue
				}

				colors = append(colors, &Color{
					Rgb: int32(unit.R)<<16 | int32(unit.G)<<8 | int32(unit.B),
				})
			}
		}

		leds = &LedsConfiguration_Gradient{
			Gradient: &LedsGradient{
				Speed:  float32(100-bean.Loop_time) / 100,
				Colors: colors,
			},
		}

	default:
		panic("TODO: implement")
	}

	return &LedsConfiguration{
		Leds:       leds,
		Brightness: float32(bean.Light_scale) / 255,
	}
}

func (c *LedsConfiguration) ApplyTo(bean *config.NewLedConfigBean) {
	bean.Light_scale = byte(c.Brightness * 100)

	switch leds := c.Leds.(type) {
	case *LedsConfiguration_Off:

	case *LedsConfiguration_Steady:
		bean.SetSteady(*leds.Steady.Color.LedUnit())

	case *LedsConfiguration_Streamlined:
		bean.SetStreamlined(leds.Streamlined.Speed)

	case *LedsConfiguration_Gradient:
		g := leds.Gradient

		colors := make([]config.LedUnit, 0, len(g.Colors))
		for _, c := range g.Colors {
			if c == nil {
				continue
			}
			colors = append(colors, *c.LedUnit())
		}

		if len(colors) < 2 {
			return
		}

		bean.SetGradient(colors, g.Speed)
	}
}
