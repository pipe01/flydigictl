package commands

import (
	"errors"
	"fmt"
	"log"
	"os"

	godbus "github.com/godbus/dbus/v5"
	"github.com/pipe01/flydigi-linux/pkg/dbus"
	"github.com/pipe01/flydigi-linux/pkg/dbus/client"
	"github.com/pipe01/flydigi-linux/pkg/dbus/pb"
	"github.com/spf13/cobra"
)

var dbusClient *client.Client

var (
	weConnectedGamepad bool

	terseOutput          bool
	persistConnection    bool
	forceCloseConnection bool
)

var rootCmd = &cobra.Command{
	Use:   "flydigictl",
	Short: "Control Flydigi gamepads",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := connectDBus(); err != nil {
			log.Fatalf("failed to connect to dbus service: %s", err)
		}
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&terseOutput, "terse", "t", false, "Output only the queried value with no decoration")
	rootCmd.PersistentFlags().BoolVar(&persistConnection, "persist-conn", false, "Don't close the connection to the gamepad after the command exits")
	rootCmd.PersistentFlags().BoolVar(&forceCloseConnection, "force-close-conn", false, "Forcibly close the connection to the gamepad after the command exits")
}

func connectDBus() error {
	cl, err := client.Dial()
	if err != nil {
		return err
	}

	dbusClient = cl
	return nil
}

func connectGamepad() error {
	err := dbusClient.Connect()
	if err != nil {
		var ferr client.FlydigiError
		if errors.As(err, &ferr) && ferr.Name == dbus.ErrorAlreadyConnected {
			return nil
		}

		var dberr godbus.Error
		if errors.As(err, &dberr) && dberr.Name == "org.freedesktop.DBus.Error.ServiceUnknown" {
			return errors.New("flydigid is not running")
		}

		return err
	}

	weConnectedGamepad = true
	return nil
}

func disconnectGamepad() error {
	if forceCloseConnection || (weConnectedGamepad && !persistConnection) {
		return dbusClient.Disconnect()
	}

	return nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func readConfiguration[T any](fn func(conf *pb.GamepadConfiguration) T) (ret T, err error) {
	if err := connectGamepad(); err != nil {
		return ret, fmt.Errorf("connect to gamepad: %w", err)
	}
	defer disconnectGamepad()

	cfg, err := dbusClient.GetConfiguration()
	if err != nil {
		return ret, fmt.Errorf("get config: %w", err)
	}

	return fn(cfg), nil
}

func modifyConfiguration(fn func(conf *pb.GamepadConfiguration)) error {
	if err := connectGamepad(); err != nil {
		return fmt.Errorf("connect to gamepad: %w", err)
	}
	defer disconnectGamepad()

	cfg, err := dbusClient.GetConfiguration()
	if err != nil {
		return fmt.Errorf("get config: %w", err)
	}

	fn(cfg)

	err = dbusClient.SetConfiguration(cfg)
	if err != nil {
		return fmt.Errorf("set config: %w", err)
	}

	return nil
}
