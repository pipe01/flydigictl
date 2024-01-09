package client

import (
	"fmt"

	"github.com/godbus/dbus/v5"
	common "github.com/pipe01/flydigictl/pkg/dbus"
	"github.com/pipe01/flydigictl/pkg/dbus/pb"
	"google.golang.org/protobuf/proto"
)

type FlydigiError struct {
	Name    string
	Message string
}

func (e FlydigiError) Error() string {
	return e.Message
}

type Client struct {
	conn *dbus.Conn
	obj  dbus.BusObject
}

func Dial() (*Client, error) {
	conn, err := dbus.ConnectSystemBus()
	if err != nil {
		return nil, fmt.Errorf("connect to system bus: %w", err)
	}

	dbusObj := conn.Object(common.InterfaceName, common.ObjectPath)

	return &Client{
		conn: conn,
		obj:  dbusObj,
	}, nil
}

func (c *Client) call(methodName string, args []interface{}, retvalues ...interface{}) error {
	call := c.obj.Call(fmt.Sprintf("%s.%s", common.InterfaceName, methodName), 0, args...)

	return call.Store(retvalues...)
}

func (c *Client) Connect() error {
	return c.wrapError(c.call("Connect", nil))
}

func (c *Client) Disconnect() error {
	return c.wrapError(c.call("Disconnect", nil))
}

func (c *Client) GetServerVersion() (string, error) {
	var version string

	if err := c.call("GetServerVersion", nil, &version); err != nil {
		return "", err
	}

	return version, nil
}

func (c *Client) GetConfiguration() (*pb.GamepadConfiguration, error) {
	var cfgBytes []byte

	if err := c.call("GetConfiguration", nil, &cfgBytes); err != nil {
		return nil, c.wrapError(err)
	}

	var cfg pb.GamepadConfiguration

	if err := proto.Unmarshal(cfgBytes, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &cfg, nil
}

func (c *Client) SetConfiguration(cfg *pb.GamepadConfiguration) error {
	cfgBytes, err := proto.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	if err := c.call("SetConfiguration", []any{cfgBytes}); err != nil {
		return c.wrapError(err)
	}

	return nil
}

func (c *Client) GetDeviceInfo() (*pb.GamepadInfo, error) {
	var infoBytes []byte

	if err := c.call("GetDeviceInfo", nil, &infoBytes); err != nil {
		return nil, c.wrapError(err)
	}

	var info pb.GamepadInfo

	if err := proto.Unmarshal(infoBytes, &info); err != nil {
		return nil, fmt.Errorf("unmarshal info: %w", err)
	}

	return &info, nil
}

func (c *Client) wrapError(err error) error {
	switch err := err.(type) {
	case dbus.Error:
		var msg string

		switch err.Name {
		case common.ErrorNotConnected:
			msg = "gamepad not connected"
		case common.ErrorAlreadyConnected:
			msg = "gamepad already connected"
		case common.ErrorMarshallingFault:
			msg = "marshalling failed"
		case common.ErrorGamepadWritingFault:
			msg = "writing to gamepad failed"
		case common.ErrorGamepadReadingFault:
			msg = "reading from gamepad failed"
		case common.ErrorGamepadNotFound:
			msg = "gamepad not found"
		default:
			return err
		}

		return FlydigiError{Name: err.Name, Message: msg}
	}

	return err
}
