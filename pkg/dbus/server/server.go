package server

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/pipe01/flydigictl/pkg/dbus/pb"
	"github.com/pipe01/flydigictl/pkg/flydigi"
	"github.com/pipe01/flydigictl/pkg/flydigi/protocol"
	"github.com/pipe01/flydigictl/pkg/version"
	"google.golang.org/protobuf/proto"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
	"github.com/gookit/goutil/dump"
	"github.com/rs/zerolog/log"

	common "github.com/pipe01/flydigictl/pkg/dbus"
)

type Server struct {
	connectmu sync.Mutex

	gp *flydigi.Gamepad
}

func New() *Server {
	return &Server{}
}

func (s *Server) checkConnected() *dbus.Error {
	if s.gp == nil {
		return makeError(common.ErrorNotConnected, nil)
	}

	return nil
}

func (s *Server) Connect() *dbus.Error {
	s.connectmu.Lock()
	defer s.connectmu.Unlock()

	if s.gp != nil {
		return makeError(common.ErrorAlreadyConnected, nil)
	}

	dev, err := flydigi.OpenGamepad()
	if err != nil {
		if errors.Is(err, protocol.ErrGamepadNotPresent) {
			return makeError(common.ErrorGamepadNotFound, nil)
		}

		return dbus.MakeFailedError(err)
	}

	closech := make(chan struct{})
	go func() {
		defer close(closech)

		dev.NotifyClose(closech)

		<-closech

		s.gp = nil
	}()

	s.gp = dev
	return nil
}

func (s *Server) Disconnect() *dbus.Error {
	if err := s.checkConnected(); err != nil {
		return err
	}

	if err := s.gp.Close(); err != nil {
		log.Err(err).Msg("failed to close gamepad")
	}

	return nil
}

func (s *Server) GetServerVersion() (string, *dbus.Error) {
	return version.Version, nil
}

func (s *Server) DumpConfiguration() (string, *dbus.Error) {
	if err := s.checkConnected(); err != nil {
		return "", err
	}

	conf, err := s.gp.GetConfig()
	if err != nil {
		return "", makeError(common.ErrorGamepadReadingFault, err)
	}

	conf.Basic.NewLedConfig, err = s.gp.GetLEDConfig()
	if err != nil {
		return "", makeError(common.ErrorGamepadReadingFault, err)
	}

	var str strings.Builder

	dump.NewDumper(&str, 0).Dump(conf)

	return str.String(), nil
}

func (s *Server) GetConfiguration() ([]byte, *dbus.Error) {
	if err := s.checkConnected(); err != nil {
		return nil, err
	}

	conf, err := s.gp.GetConfig()
	if err != nil {
		return nil, makeError(common.ErrorGamepadReadingFault, err)
	}

	prot := pb.ConvertGamepadConfiguration(conf)

	data, err := proto.Marshal(prot)
	if err != nil {
		return nil, makeError(common.ErrorMarshallingFault, err)
	}

	return data, nil
}

func (s *Server) SetConfiguration(data []byte) *dbus.Error {
	if err := s.checkConnected(); err != nil {
		return err
	}

	var conf pb.GamepadConfiguration

	err := proto.Unmarshal(data, &conf)
	if err != nil {
		return makeError(common.ErrorMarshallingFault, err)
	}

	gpConf, err := s.gp.GetConfig()
	if err != nil {
		return makeError(common.ErrorGamepadReadingFault, err)
	}

	conf.ApplyTo(gpConf)

	err = s.gp.SaveConfig(gpConf)
	if err != nil {
		return makeError(common.ErrorGamepadWritingFault, err)
	}

	return nil
}

func (s *Server) GetLEDConfiguration() ([]byte, *dbus.Error) {
	if err := s.checkConnected(); err != nil {
		return nil, err
	}

	conf, err := s.gp.GetLEDConfig()
	if err != nil {
		return nil, makeError(common.ErrorGamepadReadingFault, err)
	}

	prot := pb.ConvertLEDConfiguration(conf)

	data, err := proto.Marshal(prot)
	if err != nil {
		return nil, makeError(common.ErrorMarshallingFault, err)
	}

	return data, nil
}

func (s *Server) SetLEDConfiguration(data []byte) *dbus.Error {
	if err := s.checkConnected(); err != nil {
		return err
	}

	var conf pb.LedsConfiguration

	err := proto.Unmarshal(data, &conf)
	if err != nil {
		return makeError(common.ErrorMarshallingFault, err)
	}

	ledConf, err := s.gp.GetLEDConfig()
	if err != nil {
		return makeError(common.ErrorGamepadReadingFault, err)
	}

	conf.ApplyTo(ledConf)

	err = s.gp.SaveLEDConfig(ledConf)
	if err != nil {
		return makeError(common.ErrorGamepadWritingFault, err)
	}

	return nil
}

func (s *Server) GetDeviceInfo() ([]byte, *dbus.Error) {
	if err := s.checkConnected(); err != nil {
		return nil, err
	}

	info, err := s.gp.GetGamepadInfo()
	if err != nil {
		return nil, makeError(common.ErrorGamepadWritingFault, err)
	}

	prot := pb.GamepadInfo{
		DeviceId:       info.DeviceId,
		BatteryPercent: info.BatteryPercent,
		ConnectionType: pb.ConnectionType(info.ConnectType),
		CpuType:        info.CpuType,
		CpuName:        info.CpuName,
	}

	data, err := proto.Marshal(&prot)
	if err != nil {
		return nil, makeError(common.ErrorMarshallingFault, err)
	}

	return data, nil
}

func (s *Server) Listen() error {
	conn, err := dbus.ConnectSystemBus()
	if err != nil {
		return fmt.Errorf("connect to system bus: %w", err)
	}
	defer conn.Close()

	intros := introspect.NewIntrospectable(&introspect.Node{
		Interfaces: []introspect.Interface{
			{
				Name: common.InterfaceName,
				Methods: []introspect.Method{
					{
						Name: "Connect",
					},
					{
						Name: "Disconnect",
					},
					{
						Name: "GetServerVersion",
						Args: []introspect.Arg{
							{Direction: "out", Type: "s"},
						},
					},
					{
						Name: "GetConfiguration",
						Args: []introspect.Arg{
							{Direction: "out", Type: "ay"},
						},
					},
					{
						Name: "SetConfiguration",
						Args: []introspect.Arg{
							{Direction: "in", Type: "ay"},
						},
					},
				},
			},
		},
	})

	err = conn.Export(s, common.ObjectPath, common.InterfaceName)
	conn.Export(intros, common.ObjectPath, "org.freedesktop.DBus.Introspectable")

	reply, err := conn.RequestName(common.InterfaceName, dbus.NameFlagDoNotQueue)
	if err != nil {
		return fmt.Errorf("request dbus name: %w", err)
	}
	if reply != dbus.RequestNameReplyPrimaryOwner {
		return fmt.Errorf("dbus name already taken")
	}

	log.Info().Msg("connected to dbus")
	select {}
}

func makeError(name string, err error) *dbus.Error {
	var errstr string
	if err != nil {
		errstr = err.Error()
	}
	return dbus.NewError(name, []interface{}{errstr})
}
