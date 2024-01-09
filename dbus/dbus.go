package dbus

import (
	"fmt"
	"sync"

	"github.com/pipe01/flydigi-linux/dbus/pb"
	"github.com/pipe01/flydigi-linux/flydigi"
	"google.golang.org/protobuf/proto"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
	"github.com/rs/zerolog/log"
)

const (
	interfaceName = "com.pipe01.flydigi.Gamepad"
	objectPath    = "/com/pipe01/flydigi/Gamepad"
)

type Server struct {
	connectmu sync.Mutex

	gp *flydigi.Gamepad
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) checkConnected() *dbus.Error {
	if s.gp == nil {
		return dbus.NewError("org.pipe01.flydigi.Error.NotConnected", nil)
	}

	return nil
}

func (s *Server) Connect() *dbus.Error {
	s.connectmu.Lock()
	defer s.connectmu.Unlock()

	if s.gp != nil {
		return dbus.NewError("com.pipe01.flydigi.Error.AlreadyConnected", nil)
	}

	dev, err := flydigi.OpenGamepad()
	if err != nil {
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

func (s *Server) GetConfiguration() ([]byte, *dbus.Error) {
	if err := s.checkConnected(); err != nil {
		return nil, err
	}

	conf, err := s.gp.GetConfig()
	if err != nil {
		return nil, dbus.MakeFailedError(fmt.Errorf("get gamepad conf: %w", err))
	}

	prot := pb.GetGamepadConfiguration(conf)

	data, err := proto.Marshal(prot)
	if err != nil {
		return nil, dbus.MakeFailedError(fmt.Errorf("marshal conf: %w", err))
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
		return dbus.MakeFailedError(fmt.Errorf("unmarshal conf: %w", err))
	}

	gpConf, err := s.gp.GetConfig()
	if err != nil {
		return dbus.MakeFailedError(fmt.Errorf("get gamepad conf: %w", err))
	}

	conf.ApplyTo(gpConf)

	err = s.gp.SaveConfig(gpConf)
	if err != nil {
		return dbus.MakeFailedError(fmt.Errorf("save gamepad conf: %w", err))
	}

	return nil
}

func (s *Server) Listen() error {
	log.Info().Msg("connecting to dbus")

	conn, err := dbus.ConnectSystemBus()
	if err != nil {
		return fmt.Errorf("connect to system bus: %w", err)
	}
	defer conn.Close()

	intros := introspect.NewIntrospectable(&introspect.Node{
		Interfaces: []introspect.Interface{
			{
				Name: interfaceName,
				Methods: []introspect.Method{
					{
						Name: "Connect",
					},
					{
						Name: "Disconnect",
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

	err = conn.Export(s, objectPath, interfaceName)
	conn.Export(intros, objectPath, "org.freedesktop.DBus.Introspectable")

	reply, err := conn.RequestName(interfaceName, dbus.NameFlagDoNotQueue)
	if err != nil {
		return fmt.Errorf("request dbus name: %w", err)
	}
	if reply != dbus.RequestNameReplyPrimaryOwner {
		return fmt.Errorf("dbus name already taken")
	}

	select {}
}
