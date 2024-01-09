package dbus

import "github.com/godbus/dbus/v5"

const (
	errorsPrefix = "org.pipe01.flydigi.Error."

	ErrorNotConnected        = errorsPrefix + "NotConnected"
	ErrorAlreadyConnected    = errorsPrefix + "AlreadyConnected"
	ErrorMarshallingFault    = errorsPrefix + "MarshallingFailed"
	ErrorGamepadWritingFault = errorsPrefix + "GamepadWritingFault"
	ErrorGamepadReadingFault = errorsPrefix + "GamepadReadingFault"
)

func IsFlydigiErr(err error, errName string) bool {
	if dberr, ok := err.(dbus.Error); ok {
		if dberr.Name == errName {
			return true
		}
	}

	return false
}
