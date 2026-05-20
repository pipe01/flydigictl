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

func clampUnit(v float32) float32 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

func brightnessFromScale(scale byte) float32 {
	if scale > 100 {
		return 1
	}

	return float32(scale) / 100
}

func ConvertGamepadConfiguration(bean *config.AllConfigBean) *GamepadConfiguration {
	cfg := &GamepadConfiguration{}
	if bean == nil {
		return cfg
	}

	if bean.JoyMapping != nil {
		cfg.LeftJoystick = ConvertJoystickConfiguration(bean.JoyMapping.LeftJoystic)
		cfg.RightJoystick = ConvertJoystickConfiguration(bean.JoyMapping.RightJoystic)
	}
	cfg.ButtonMappings = ConvertButtonMappings(bean.KeyMapping)

	return cfg
}

func (c *GamepadConfiguration) ApplyTo(bean *config.AllConfigBean) {
	if c == nil || bean == nil {
		return
	}

	if c.LeftJoystick != nil && bean.JoyMapping != nil {
		c.LeftJoystick.ApplyTo(bean.JoyMapping.LeftJoystic)
	}
	if c.RightJoystick != nil && bean.JoyMapping != nil {
		c.RightJoystick.ApplyTo(bean.JoyMapping.RightJoystic)
	}

	for _, mapping := range c.ButtonMappings {
		mapping.ApplyTo(bean.KeyMapping)
	}
}

func ConvertJoystickConfiguration(bean *config.JoyStickBean) *JoystickConfiguration {
	if bean == nil || bean.Curve == nil {
		return &JoystickConfiguration{}
	}

	return &JoystickConfiguration{
		Deadzone: bean.Curve.Zero,
	}
}

func (c *JoystickConfiguration) ApplyTo(bean *config.JoyStickBean) {
	if c == nil || bean == nil || bean.Curve == nil {
		return
	}

	bean.Curve.Zero = c.Deadzone
}

func ConvertButtonMappings(beans []*config.KeyMappingBean) []*ButtonMapping {
	mappings := make([]*ButtonMapping, 0, len(beans))

	for _, bean := range beans {
		if bean == nil {
			continue
		}

		mappedKeyID, _ := strconv.Atoi(bean.MapData)

		mappings = append(mappings, &ButtonMapping{
			KeyId:       bean.KeyId,
			Key:         bean.Key,
			MappedKeyId: int32(mappedKeyID),
			MappedKey:   bean.MapKey,
			IsOwn:       bean.IsOwn,
			IsMapped:    bean.IsMap,
		})
	}

	return mappings
}

func (c *ButtonMapping) ApplyTo(beans []*config.KeyMappingBean) {
	if c == nil || c.KeyId < 0 || int(c.KeyId) >= len(beans) || c.MappedKeyId < 0 || c.MappedKeyId > 255 {
		return
	}

	bean := beans[c.KeyId]
	if bean == nil {
		return
	}

	bean.MapData = strconv.Itoa(int(c.MappedKeyId))
	bean.MapKey = c.MappedKey
	bean.IsMap = c.KeyId != c.MappedKeyId
	bean.MapType = config.KeyMapTypeGamePad
	bean.Turbo = 0
	bean.TurboType = 0
}

func ConvertLEDConfiguration(bean *config.NewLedConfigBean) *LedsConfiguration {
	var leds isLedsConfiguration_Leds
	colorAt := func(group, unit int) *Color {
		if len(bean.LedGroups) <= group || len(bean.LedGroups[group].Units) <= unit {
			return ColorFromRGB(0, 0, 0)
		}

		u := bean.LedGroups[group].Units[unit]
		return ColorFromRGB(u.R, u.G, u.B)
	}

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

	case config.LedModeBreathing:
		leds = &LedsConfiguration_Breathing{
			Breathing: &LedsBreathing{
				Color: colorAt(0, 0),
				Speed: float32(100-bean.Loop_time) / 100,
			},
		}

	case config.LedModeGradient:
		leds = &LedsConfiguration_Gradient{
			Gradient: &LedsGradient{
				StartColor: colorAt(0, 0),
				EndColor:   colorAt(0, 1),
				Speed:      float32(100-bean.Loop_time) / 100,
			},
		}

	case config.LedModeFeedback:
		leds = &LedsConfiguration_Feedback{
			Feedback: &LedsFeedback{
				Speed: float32(100-bean.Loop_time) / 100,
			},
		}

	default:
		leds = nil
	}

	return &LedsConfiguration{
		Leds:       leds,
		Brightness: brightnessFromScale(bean.Light_scale),
	}
}

func (c *LedsConfiguration) ApplyTo(bean *config.NewLedConfigBean) {
	bean.Light_scale = byte(clampUnit(c.Brightness) * 100)

	switch leds := c.Leds.(type) {
	case *LedsConfiguration_Off:
		bean.SetOff()

	case *LedsConfiguration_Steady:
		bean.SetSteady(*leds.Steady.Color.LedUnit())

	case *LedsConfiguration_Streamlined:
		bean.SetStreamlined(leds.Streamlined.Speed)

	case *LedsConfiguration_Breathing:
		bean.SetBreathing(*leds.Breathing.Color.LedUnit(), leds.Breathing.Speed)

	case *LedsConfiguration_Gradient:
		bean.SetGradient(*leds.Gradient.StartColor.LedUnit(), *leds.Gradient.EndColor.LedUnit(), leds.Gradient.Speed)

	case *LedsConfiguration_Feedback:
		bean.SetFeedback(leds.Feedback.Speed)
	}
}
