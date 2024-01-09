package commands

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/pipe01/flydigictl/pkg/dbus/pb"
	"github.com/spf13/cobra"
)

type joystickSide string

const (
	joystickLeft  joystickSide = "left"
	joystickRight joystickSide = "right"
)

func (s joystickSide) GetBean(b *pb.GamepadConfiguration) *pb.JoystickConfiguration {
	switch s {
	case joystickLeft:
		return b.LeftJoystick
	case joystickRight:
		return b.RightJoystick
	}

	panic("invalid joystick side")
}

var joysticksCommand = &cobra.Command{
	Use:   "joystick",
	Short: "Manage joystick configuration",
}

func init() {
	joysticksCommand.AddCommand(genJoystickCommand(joystickLeft))
	joysticksCommand.AddCommand(genJoystickCommand(joystickRight))

	rootCmd.AddCommand(joysticksCommand)
}

func genJoystickCommand(side joystickSide) *cobra.Command {
	cmd := &cobra.Command{
		Use:   string(side),
		Short: fmt.Sprintf("Manage %s joystick configuration", side),
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "deadzone",
		Short: "Get or set this joystick's deadzone",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			switch len(args) {
			case 0: // Read deadzone
				deadzone, err := readConfiguration(func(conf *pb.GamepadConfiguration) int32 {
					return side.GetBean(conf).Deadzone
				})
				if err != nil {
					return fmt.Errorf("get configuration: %w", err)
				}

				if terseOutput {
					fmt.Println(deadzone)
				} else {
					fmt.Printf("%s joystick deadzone (0-100): %d\n", side, deadzone)
				}

			case 1: // Write deadzone
				val, err := strconv.Atoi(args[0])
				if err != nil {
					return fmt.Errorf("parse value: %w", err)
				}
				if val < 0 || val > 100 {
					return errors.New("value is outside range (0-100)")
				}

				return modifyConfiguration(func(conf *pb.GamepadConfiguration) {
					side.GetBean(conf).Deadzone = int32(val)
				})
			}

			return nil
		},
	})

	return cmd
}
