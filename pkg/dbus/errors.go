package dbus

const (
	errorsPrefix = "org.pipe01.flydigi.Error."

	ErrorNotConnected        = errorsPrefix + "NotConnected"
	ErrorAlreadyConnected    = errorsPrefix + "AlreadyConnected"
	ErrorMarshallingFault    = errorsPrefix + "MarshallingFailed"
	ErrorGamepadWritingFault = errorsPrefix + "GamepadWritingFault"
	ErrorGamepadReadingFault = errorsPrefix + "GamepadReadingFault"
)
