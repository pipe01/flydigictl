package dbus

import (
	"fmt"
	"sync"

	"github.com/pipe01/flydigi-linux/flydigi"

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
						// Args: []introspect.Arg{
						// 	{Direction: "out", Type: "s"},
						// },
					},
					{
						Name: "Test",
						Args: []introspect.Arg{
							{Direction: "out", Type: "s"},
							{Direction: "out", Type: "s"},
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
