package pb

import "github.com/pipe01/flydigi-linux/pkg/flydigi/config"

func GetGamepadConfiguration(bean *config.AllConfigBean) *GamepadConfiguration {
	return &GamepadConfiguration{
		LeftJoystick:  GetJoystickConfiguration(bean.JoyMapping.LeftJoystic),
		RightJoystick: GetJoystickConfiguration(bean.JoyMapping.RightJoystic),
	}
}

func (c *GamepadConfiguration) ApplyTo(bean *config.AllConfigBean) {
	c.LeftJoystick.ApplyTo(bean.JoyMapping.LeftJoystic)
	c.RightJoystick.ApplyTo(bean.JoyMapping.RightJoystic)
}

func GetJoystickConfiguration(bean *config.JoyStickBean) *JoystickConfiguration {
	return &JoystickConfiguration{
		Deadzone: bean.Curve.Zero,
	}
}

func (c *JoystickConfiguration) ApplyTo(bean *config.JoyStickBean) {
	bean.Curve.Zero = c.Deadzone
}
