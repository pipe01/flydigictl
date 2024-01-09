package commands

import (
	"fmt"
	"log"
	"os"

	common "github.com/pipe01/flydigi-linux/pkg/dbus"
	"github.com/pipe01/flydigi-linux/pkg/dbus/client"
	"github.com/pipe01/flydigi-linux/pkg/dbus/pb"
	"github.com/spf13/cobra"
)

var dbusClient *client.Client

var (
	weConnectedGamepad bool

	terseOutput bool
)

var rootCmd = &cobra.Command{
	Use:   "flydigictl",
	Short: "Control Flydigi gamepads",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := connectDBus(); err != nil {
			log.Fatalf("failed to connect to dbus service: %s", err)
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if weConnectedGamepad {
			dbusClient.Disconnect()
		}
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&terseOutput, "terse", "t", false, "Output only the queried value with no decoration")
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
		if common.IsFlydigiErr(err, common.ErrorAlreadyConnected) {
			return nil
		}

		return err
	}

	weConnectedGamepad = true
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
	defer dbusClient.Disconnect()

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
	defer dbusClient.Disconnect()

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
