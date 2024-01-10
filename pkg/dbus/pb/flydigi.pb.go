// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v3.12.4
// source: flydigi.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ConnectionType int32

const (
	ConnectionType_UNKNOWN  ConnectionType = 0
	ConnectionType_WIRELESS ConnectionType = 1
	ConnectionType_WIRED    ConnectionType = 2
)

// Enum value maps for ConnectionType.
var (
	ConnectionType_name = map[int32]string{
		0: "UNKNOWN",
		1: "WIRELESS",
		2: "WIRED",
	}
	ConnectionType_value = map[string]int32{
		"UNKNOWN":  0,
		"WIRELESS": 1,
		"WIRED":    2,
	}
)

func (x ConnectionType) Enum() *ConnectionType {
	p := new(ConnectionType)
	*p = x
	return p
}

func (x ConnectionType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ConnectionType) Descriptor() protoreflect.EnumDescriptor {
	return file_flydigi_proto_enumTypes[0].Descriptor()
}

func (ConnectionType) Type() protoreflect.EnumType {
	return &file_flydigi_proto_enumTypes[0]
}

func (x ConnectionType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ConnectionType.Descriptor instead.
func (ConnectionType) EnumDescriptor() ([]byte, []int) {
	return file_flydigi_proto_rawDescGZIP(), []int{0}
}

type GamepadInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DeviceId       int32          `protobuf:"varint,1,opt,name=device_id,json=deviceId,proto3" json:"device_id,omitempty"`
	BatteryPercent int32          `protobuf:"varint,2,opt,name=battery_percent,json=batteryPercent,proto3" json:"battery_percent,omitempty"`
	ConnectionType ConnectionType `protobuf:"varint,3,opt,name=connection_type,json=connectionType,proto3,enum=flydigi.ConnectionType" json:"connection_type,omitempty"`
	CpuType        string         `protobuf:"bytes,4,opt,name=cpu_type,json=cpuType,proto3" json:"cpu_type,omitempty"`
	CpuName        string         `protobuf:"bytes,5,opt,name=cpu_name,json=cpuName,proto3" json:"cpu_name,omitempty"`
}

func (x *GamepadInfo) Reset() {
	*x = GamepadInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_flydigi_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GamepadInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GamepadInfo) ProtoMessage() {}

func (x *GamepadInfo) ProtoReflect() protoreflect.Message {
	mi := &file_flydigi_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GamepadInfo.ProtoReflect.Descriptor instead.
func (*GamepadInfo) Descriptor() ([]byte, []int) {
	return file_flydigi_proto_rawDescGZIP(), []int{0}
}

func (x *GamepadInfo) GetDeviceId() int32 {
	if x != nil {
		return x.DeviceId
	}
	return 0
}

func (x *GamepadInfo) GetBatteryPercent() int32 {
	if x != nil {
		return x.BatteryPercent
	}
	return 0
}

func (x *GamepadInfo) GetConnectionType() ConnectionType {
	if x != nil {
		return x.ConnectionType
	}
	return ConnectionType_UNKNOWN
}

func (x *GamepadInfo) GetCpuType() string {
	if x != nil {
		return x.CpuType
	}
	return ""
}

func (x *GamepadInfo) GetCpuName() string {
	if x != nil {
		return x.CpuName
	}
	return ""
}

type GamepadConfiguration struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LeftJoystick  *JoystickConfiguration `protobuf:"bytes,1,opt,name=left_joystick,json=leftJoystick,proto3" json:"left_joystick,omitempty"`
	RightJoystick *JoystickConfiguration `protobuf:"bytes,2,opt,name=right_joystick,json=rightJoystick,proto3" json:"right_joystick,omitempty"`
}

func (x *GamepadConfiguration) Reset() {
	*x = GamepadConfiguration{}
	if protoimpl.UnsafeEnabled {
		mi := &file_flydigi_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GamepadConfiguration) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GamepadConfiguration) ProtoMessage() {}

func (x *GamepadConfiguration) ProtoReflect() protoreflect.Message {
	mi := &file_flydigi_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GamepadConfiguration.ProtoReflect.Descriptor instead.
func (*GamepadConfiguration) Descriptor() ([]byte, []int) {
	return file_flydigi_proto_rawDescGZIP(), []int{1}
}

func (x *GamepadConfiguration) GetLeftJoystick() *JoystickConfiguration {
	if x != nil {
		return x.LeftJoystick
	}
	return nil
}

func (x *GamepadConfiguration) GetRightJoystick() *JoystickConfiguration {
	if x != nil {
		return x.RightJoystick
	}
	return nil
}

type JoystickConfiguration struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Deadzone int32 `protobuf:"varint,1,opt,name=deadzone,proto3" json:"deadzone,omitempty"`
}

func (x *JoystickConfiguration) Reset() {
	*x = JoystickConfiguration{}
	if protoimpl.UnsafeEnabled {
		mi := &file_flydigi_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JoystickConfiguration) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JoystickConfiguration) ProtoMessage() {}

func (x *JoystickConfiguration) ProtoReflect() protoreflect.Message {
	mi := &file_flydigi_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JoystickConfiguration.ProtoReflect.Descriptor instead.
func (*JoystickConfiguration) Descriptor() ([]byte, []int) {
	return file_flydigi_proto_rawDescGZIP(), []int{2}
}

func (x *JoystickConfiguration) GetDeadzone() int32 {
	if x != nil {
		return x.Deadzone
	}
	return 0
}

type LedsOff struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *LedsOff) Reset() {
	*x = LedsOff{}
	if protoimpl.UnsafeEnabled {
		mi := &file_flydigi_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LedsOff) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LedsOff) ProtoMessage() {}

func (x *LedsOff) ProtoReflect() protoreflect.Message {
	mi := &file_flydigi_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LedsOff.ProtoReflect.Descriptor instead.
func (*LedsOff) Descriptor() ([]byte, []int) {
	return file_flydigi_proto_rawDescGZIP(), []int{3}
}

type LedsSteady struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Color *Color `protobuf:"bytes,1,opt,name=color,proto3" json:"color,omitempty"`
}

func (x *LedsSteady) Reset() {
	*x = LedsSteady{}
	if protoimpl.UnsafeEnabled {
		mi := &file_flydigi_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LedsSteady) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LedsSteady) ProtoMessage() {}

func (x *LedsSteady) ProtoReflect() protoreflect.Message {
	mi := &file_flydigi_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LedsSteady.ProtoReflect.Descriptor instead.
func (*LedsSteady) Descriptor() ([]byte, []int) {
	return file_flydigi_proto_rawDescGZIP(), []int{4}
}

func (x *LedsSteady) GetColor() *Color {
	if x != nil {
		return x.Color
	}
	return nil
}

type LedsStreamlined struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Speed float32 `protobuf:"fixed32,1,opt,name=speed,proto3" json:"speed,omitempty"`
}

func (x *LedsStreamlined) Reset() {
	*x = LedsStreamlined{}
	if protoimpl.UnsafeEnabled {
		mi := &file_flydigi_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LedsStreamlined) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LedsStreamlined) ProtoMessage() {}

func (x *LedsStreamlined) ProtoReflect() protoreflect.Message {
	mi := &file_flydigi_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LedsStreamlined.ProtoReflect.Descriptor instead.
func (*LedsStreamlined) Descriptor() ([]byte, []int) {
	return file_flydigi_proto_rawDescGZIP(), []int{5}
}

func (x *LedsStreamlined) GetSpeed() float32 {
	if x != nil {
		return x.Speed
	}
	return 0
}

type LedsConfiguration struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Brightness float32 `protobuf:"fixed32,1,opt,name=brightness,proto3" json:"brightness,omitempty"`
	// Types that are assignable to Leds:
	//
	//	*LedsConfiguration_Off
	//	*LedsConfiguration_Steady
	//	*LedsConfiguration_Streamlined
	Leds isLedsConfiguration_Leds `protobuf_oneof:"leds"`
}

func (x *LedsConfiguration) Reset() {
	*x = LedsConfiguration{}
	if protoimpl.UnsafeEnabled {
		mi := &file_flydigi_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LedsConfiguration) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LedsConfiguration) ProtoMessage() {}

func (x *LedsConfiguration) ProtoReflect() protoreflect.Message {
	mi := &file_flydigi_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LedsConfiguration.ProtoReflect.Descriptor instead.
func (*LedsConfiguration) Descriptor() ([]byte, []int) {
	return file_flydigi_proto_rawDescGZIP(), []int{6}
}

func (x *LedsConfiguration) GetBrightness() float32 {
	if x != nil {
		return x.Brightness
	}
	return 0
}

func (m *LedsConfiguration) GetLeds() isLedsConfiguration_Leds {
	if m != nil {
		return m.Leds
	}
	return nil
}

func (x *LedsConfiguration) GetOff() *LedsOff {
	if x, ok := x.GetLeds().(*LedsConfiguration_Off); ok {
		return x.Off
	}
	return nil
}

func (x *LedsConfiguration) GetSteady() *LedsSteady {
	if x, ok := x.GetLeds().(*LedsConfiguration_Steady); ok {
		return x.Steady
	}
	return nil
}

func (x *LedsConfiguration) GetStreamlined() *LedsStreamlined {
	if x, ok := x.GetLeds().(*LedsConfiguration_Streamlined); ok {
		return x.Streamlined
	}
	return nil
}

type isLedsConfiguration_Leds interface {
	isLedsConfiguration_Leds()
}

type LedsConfiguration_Off struct {
	Off *LedsOff `protobuf:"bytes,10,opt,name=off,proto3,oneof"`
}

type LedsConfiguration_Steady struct {
	Steady *LedsSteady `protobuf:"bytes,12,opt,name=steady,proto3,oneof"`
}

type LedsConfiguration_Streamlined struct {
	Streamlined *LedsStreamlined `protobuf:"bytes,13,opt,name=streamlined,proto3,oneof"`
}

func (*LedsConfiguration_Off) isLedsConfiguration_Leds() {}

func (*LedsConfiguration_Steady) isLedsConfiguration_Leds() {}

func (*LedsConfiguration_Streamlined) isLedsConfiguration_Leds() {}

type Color struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Rgb int32 `protobuf:"varint,1,opt,name=rgb,proto3" json:"rgb,omitempty"`
}

func (x *Color) Reset() {
	*x = Color{}
	if protoimpl.UnsafeEnabled {
		mi := &file_flydigi_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Color) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Color) ProtoMessage() {}

func (x *Color) ProtoReflect() protoreflect.Message {
	mi := &file_flydigi_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Color.ProtoReflect.Descriptor instead.
func (*Color) Descriptor() ([]byte, []int) {
	return file_flydigi_proto_rawDescGZIP(), []int{7}
}

func (x *Color) GetRgb() int32 {
	if x != nil {
		return x.Rgb
	}
	return 0
}

var File_flydigi_proto protoreflect.FileDescriptor

var file_flydigi_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x66, 0x6c, 0x79, 0x64, 0x69, 0x67, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x07, 0x66, 0x6c, 0x79, 0x64, 0x69, 0x67, 0x69, 0x22, 0xcb, 0x01, 0x0a, 0x0b, 0x47, 0x61, 0x6d,
	0x65, 0x70, 0x61, 0x64, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x1b, 0x0a, 0x09, 0x64, 0x65, 0x76, 0x69,
	0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x64, 0x65, 0x76,
	0x69, 0x63, 0x65, 0x49, 0x64, 0x12, 0x27, 0x0a, 0x0f, 0x62, 0x61, 0x74, 0x74, 0x65, 0x72, 0x79,
	0x5f, 0x70, 0x65, 0x72, 0x63, 0x65, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0e,
	0x62, 0x61, 0x74, 0x74, 0x65, 0x72, 0x79, 0x50, 0x65, 0x72, 0x63, 0x65, 0x6e, 0x74, 0x12, 0x40,
	0x0a, 0x0f, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x79, 0x70,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x17, 0x2e, 0x66, 0x6c, 0x79, 0x64, 0x69, 0x67,
	0x69, 0x2e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65,
	0x52, 0x0e, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65,
	0x12, 0x19, 0x0a, 0x08, 0x63, 0x70, 0x75, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x63, 0x70, 0x75, 0x54, 0x79, 0x70, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x63,
	0x70, 0x75, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63,
	0x70, 0x75, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0xa2, 0x01, 0x0a, 0x14, 0x47, 0x61, 0x6d, 0x65, 0x70,
	0x61, 0x64, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x43, 0x0a, 0x0d, 0x6c, 0x65, 0x66, 0x74, 0x5f, 0x6a, 0x6f, 0x79, 0x73, 0x74, 0x69, 0x63, 0x6b,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x66, 0x6c, 0x79, 0x64, 0x69, 0x67, 0x69,
	0x2e, 0x4a, 0x6f, 0x79, 0x73, 0x74, 0x69, 0x63, 0x6b, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0c, 0x6c, 0x65, 0x66, 0x74, 0x4a, 0x6f, 0x79, 0x73,
	0x74, 0x69, 0x63, 0x6b, 0x12, 0x45, 0x0a, 0x0e, 0x72, 0x69, 0x67, 0x68, 0x74, 0x5f, 0x6a, 0x6f,
	0x79, 0x73, 0x74, 0x69, 0x63, 0x6b, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x66,
	0x6c, 0x79, 0x64, 0x69, 0x67, 0x69, 0x2e, 0x4a, 0x6f, 0x79, 0x73, 0x74, 0x69, 0x63, 0x6b, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0d, 0x72, 0x69,
	0x67, 0x68, 0x74, 0x4a, 0x6f, 0x79, 0x73, 0x74, 0x69, 0x63, 0x6b, 0x22, 0x33, 0x0a, 0x15, 0x4a,
	0x6f, 0x79, 0x73, 0x74, 0x69, 0x63, 0x6b, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x65, 0x61, 0x64, 0x7a, 0x6f, 0x6e, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x64, 0x65, 0x61, 0x64, 0x7a, 0x6f, 0x6e, 0x65,
	0x22, 0x09, 0x0a, 0x07, 0x4c, 0x65, 0x64, 0x73, 0x4f, 0x66, 0x66, 0x22, 0x32, 0x0a, 0x0a, 0x4c,
	0x65, 0x64, 0x73, 0x53, 0x74, 0x65, 0x61, 0x64, 0x79, 0x12, 0x24, 0x0a, 0x05, 0x63, 0x6f, 0x6c,
	0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x66, 0x6c, 0x79, 0x64, 0x69,
	0x67, 0x69, 0x2e, 0x43, 0x6f, 0x6c, 0x6f, 0x72, 0x52, 0x05, 0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x22,
	0x27, 0x0a, 0x0f, 0x4c, 0x65, 0x64, 0x73, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x6c, 0x69, 0x6e,
	0x65, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x70, 0x65, 0x65, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x02, 0x52, 0x05, 0x73, 0x70, 0x65, 0x65, 0x64, 0x22, 0xce, 0x01, 0x0a, 0x11, 0x4c, 0x65, 0x64,
	0x73, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1e,
	0x0a, 0x0a, 0x62, 0x72, 0x69, 0x67, 0x68, 0x74, 0x6e, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x02, 0x52, 0x0a, 0x62, 0x72, 0x69, 0x67, 0x68, 0x74, 0x6e, 0x65, 0x73, 0x73, 0x12, 0x24,
	0x0a, 0x03, 0x6f, 0x66, 0x66, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x66, 0x6c,
	0x79, 0x64, 0x69, 0x67, 0x69, 0x2e, 0x4c, 0x65, 0x64, 0x73, 0x4f, 0x66, 0x66, 0x48, 0x00, 0x52,
	0x03, 0x6f, 0x66, 0x66, 0x12, 0x2d, 0x0a, 0x06, 0x73, 0x74, 0x65, 0x61, 0x64, 0x79, 0x18, 0x0c,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x66, 0x6c, 0x79, 0x64, 0x69, 0x67, 0x69, 0x2e, 0x4c,
	0x65, 0x64, 0x73, 0x53, 0x74, 0x65, 0x61, 0x64, 0x79, 0x48, 0x00, 0x52, 0x06, 0x73, 0x74, 0x65,
	0x61, 0x64, 0x79, 0x12, 0x3c, 0x0a, 0x0b, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x6c, 0x69, 0x6e,
	0x65, 0x64, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x66, 0x6c, 0x79, 0x64, 0x69,
	0x67, 0x69, 0x2e, 0x4c, 0x65, 0x64, 0x73, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x6c, 0x69, 0x6e,
	0x65, 0x64, 0x48, 0x00, 0x52, 0x0b, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x6c, 0x69, 0x6e, 0x65,
	0x64, 0x42, 0x06, 0x0a, 0x04, 0x6c, 0x65, 0x64, 0x73, 0x22, 0x19, 0x0a, 0x05, 0x43, 0x6f, 0x6c,
	0x6f, 0x72, 0x12, 0x10, 0x0a, 0x03, 0x72, 0x67, 0x62, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x03, 0x72, 0x67, 0x62, 0x2a, 0x36, 0x0a, 0x0e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57,
	0x4e, 0x10, 0x00, 0x12, 0x0c, 0x0a, 0x08, 0x57, 0x49, 0x52, 0x45, 0x4c, 0x45, 0x53, 0x53, 0x10,
	0x01, 0x12, 0x09, 0x0a, 0x05, 0x57, 0x49, 0x52, 0x45, 0x44, 0x10, 0x02, 0x42, 0x2a, 0x5a, 0x28,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x70, 0x69, 0x70, 0x65, 0x30,
	0x31, 0x2f, 0x66, 0x6c, 0x79, 0x64, 0x69, 0x67, 0x69, 0x63, 0x74, 0x6c, 0x2f, 0x70, 0x6b, 0x67,
	0x2f, 0x64, 0x62, 0x75, 0x73, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_flydigi_proto_rawDescOnce sync.Once
	file_flydigi_proto_rawDescData = file_flydigi_proto_rawDesc
)

func file_flydigi_proto_rawDescGZIP() []byte {
	file_flydigi_proto_rawDescOnce.Do(func() {
		file_flydigi_proto_rawDescData = protoimpl.X.CompressGZIP(file_flydigi_proto_rawDescData)
	})
	return file_flydigi_proto_rawDescData
}

var file_flydigi_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_flydigi_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_flydigi_proto_goTypes = []interface{}{
	(ConnectionType)(0),           // 0: flydigi.ConnectionType
	(*GamepadInfo)(nil),           // 1: flydigi.GamepadInfo
	(*GamepadConfiguration)(nil),  // 2: flydigi.GamepadConfiguration
	(*JoystickConfiguration)(nil), // 3: flydigi.JoystickConfiguration
	(*LedsOff)(nil),               // 4: flydigi.LedsOff
	(*LedsSteady)(nil),            // 5: flydigi.LedsSteady
	(*LedsStreamlined)(nil),       // 6: flydigi.LedsStreamlined
	(*LedsConfiguration)(nil),     // 7: flydigi.LedsConfiguration
	(*Color)(nil),                 // 8: flydigi.Color
}
var file_flydigi_proto_depIdxs = []int32{
	0, // 0: flydigi.GamepadInfo.connection_type:type_name -> flydigi.ConnectionType
	3, // 1: flydigi.GamepadConfiguration.left_joystick:type_name -> flydigi.JoystickConfiguration
	3, // 2: flydigi.GamepadConfiguration.right_joystick:type_name -> flydigi.JoystickConfiguration
	8, // 3: flydigi.LedsSteady.color:type_name -> flydigi.Color
	4, // 4: flydigi.LedsConfiguration.off:type_name -> flydigi.LedsOff
	5, // 5: flydigi.LedsConfiguration.steady:type_name -> flydigi.LedsSteady
	6, // 6: flydigi.LedsConfiguration.streamlined:type_name -> flydigi.LedsStreamlined
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_flydigi_proto_init() }
func file_flydigi_proto_init() {
	if File_flydigi_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_flydigi_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GamepadInfo); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_flydigi_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GamepadConfiguration); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_flydigi_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JoystickConfiguration); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_flydigi_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LedsOff); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_flydigi_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LedsSteady); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_flydigi_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LedsStreamlined); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_flydigi_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LedsConfiguration); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_flydigi_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Color); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_flydigi_proto_msgTypes[6].OneofWrappers = []interface{}{
		(*LedsConfiguration_Off)(nil),
		(*LedsConfiguration_Steady)(nil),
		(*LedsConfiguration_Streamlined)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_flydigi_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_flydigi_proto_goTypes,
		DependencyIndexes: file_flydigi_proto_depIdxs,
		EnumInfos:         file_flydigi_proto_enumTypes,
		MessageInfos:      file_flydigi_proto_msgTypes,
	}.Build()
	File_flydigi_proto = out.File
	file_flydigi_proto_rawDesc = nil
	file_flydigi_proto_goTypes = nil
	file_flydigi_proto_depIdxs = nil
}
