package server

import (
	"errors"
	"fmt"
	"sync"

	"github.com/pipe01/flydigi-linux/pkg/dbus/pb"
	"github.com/pipe01/flydigi-linux/pkg/flydigi"
	"github.com/pipe01/flydigi-linux/pkg/flydigi/protocol"
	"google.golang.org/protobuf/proto"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
	"github.com/rs/zerolog/log"

	common "github.com/pipe01/flydigi-linux/pkg/dbus"
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

func (s *Server) GetConfiguration() ([]byte, *dbus.Error) {
	if err := s.checkConnected(); err != nil {
		return nil, err
	}

	conf, err := s.gp.GetConfig()
	if err != nil {
		return nil, makeError(common.ErrorGamepadReadingFault, err)
	}

	prot := pb.GetGamepadConfiguration(conf)

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
				Name: common.InterfaceName,
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

	err = conn.Export(s, common.ObjectPath, common.InterfaceName)
	conn.Export(intros, common.ObjectPath, "org.freedesktop.DBus.Introspectable")

	reply, err := conn.RequestName(common.InterfaceName, dbus.NameFlagDoNotQueue)
	if err != nil {
		return fmt.Errorf("request dbus name: %w", err)
	}
	if reply != dbus.RequestNameReplyPrimaryOwner {
		return fmt.Errorf("dbus name already taken")
	}

	select {}
}

func makeError(name string, err error) *dbus.Error {
	var errstr string
	if err != nil {
		errstr = err.Error()
	}
	return dbus.NewError(name, []interface{}{errstr})
}
