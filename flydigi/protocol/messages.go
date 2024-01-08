package protocol

import "errors"

var errUnknownMessage = errors.New("unknown message")

type raw []byte

func (m raw) Raw() []byte {
	return m
}

func (raw) message() {}

type MessageGamepadConfigReadCB struct {
	raw
	PackageIndex byte
	Data         []byte
}

type MessageLEDConfigReadCB struct {
	raw
	PackageIndex byte
	Data         []byte
}

type MessageDongleInfo struct {
	raw
	FW_L byte
	FW_H byte
}

type MessageWriteGamepadConfigCBK struct {
	raw
	AckNum byte
}

type MessageGamePadInfo struct {
	raw
	DeviceID         byte
	DeviceMac        []byte
	FW_L, FW_H       byte
	Battery          byte
	MotionSensorType byte
	CPUType          byte
	ConnectionType   byte
}
