package commands

import (
	"errors"
	"fmt"

	"github.com/pipe01/flydigictl/pkg/dbus/pb"
	"github.com/spf13/cobra"
)

var errInvalidColor = errors.New("invalid color, expected a color name or hex color value")

var ledsCommand = &cobra.Command{
	Use:   "leds",
	Short: "Configure lighting",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return useConnection(func() error {
			cfg, err := dbusClient.GetLEDConfiguration()
			if err != nil {
				return err
			}

			switch leds := cfg.Leds.(type) {
			case *pb.LedsConfiguration_Off:
				fmt.Println("Off")

			case *pb.LedsConfiguration_Steady:
				fmt.Printf("Steady #%06X\n", leds.Steady.Color.Rgb)

			case *pb.LedsConfiguration_Streamlined:
				fmt.Printf("Streamlined, speed %.0f%%\n", leds.Streamlined.Speed*100)

			case *pb.LedsConfiguration_Gradient:
				fmt.Printf("Gradient, speed %.0f%%\n", leds.Gradient.Speed*100)

				for i, c := range leds.Gradient.Colors {
					r, g, b := c.RGB()
					fmt.Printf("  %d: #%02X%02X%02X (%d, %d, %d)\n",
						i,
						r, g, b,
						r, g, b,
					)
				}

			}			

			return nil
		})
	},
}

var ledsBrightness float32

var ledsSteadyCommand = &cobra.Command{
	Use:     "steady",
	Short:   "Sets all LEDs to a solid color",
	Args:    cobra.ExactArgs(1),
	Example: "flydigictl leds steady #FF0000",
	RunE: func(cmd *cobra.Command, args []string) error {
		color, ok := pb.ColorFromName(args[0])
		if !ok {
			color, ok = pb.ColorFromHex(args[0])
			if !ok {
				return errInvalidColor
			}
		}

		return modifyLEDConfiguration(func(conf *pb.LedsConfiguration) {
			conf.Brightness = ledsBrightness
			conf.Leds = &pb.LedsConfiguration_Steady{Steady: &pb.LedsSteady{Color: color}}
		})
	},
}

var streamlinedSpeed float32
var ledsStreamlinedCommand = &cobra.Command{
	Use:   "streamlined",
	Short: "Sets LEDs to a streamlined effect",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return useConnection(func() error {
			return modifyLEDConfiguration(func(conf *pb.LedsConfiguration) {
				conf.Brightness = ledsBrightness
				conf.Leds = &pb.LedsConfiguration_Streamlined{Streamlined: &pb.LedsStreamlined{Speed: streamlinedSpeed}}
			})
		})
	},
}

var gradientSpeed float32
var ledsGradientCommand = &cobra.Command{
	Use:   "gradient <color1> <color2> [color3 ... color9]",
	Short: "Sets LEDs to a gradient effect",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("gradient requires at least 2 colors")
		}
		if len(args) > 9 {
			return fmt.Errorf("gradient supports at most 9 colors")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return useConnection(func() error {
			colors := make([]*pb.Color, 0, len(args))

			for _, hex := range args {
				color, ok := pb.ColorFromName(hex)
				if !ok {
					color, ok = pb.ColorFromHex(hex)
					if !ok {
						return errInvalidColor
					}
				}
				colors = append(colors, color)
			}

			return modifyLEDConfiguration(func(conf *pb.LedsConfiguration) {
				conf.Brightness = ledsBrightness
				conf.Leds = &pb.LedsConfiguration_Gradient{
					Gradient: &pb.LedsGradient{
						Speed:  gradientSpeed,
						Colors: colors,
					},
				}
			})
		})
	},
}


func init() {
	rootCmd.AddCommand(ledsCommand)

	ledsCommand.AddCommand(ledsSteadyCommand)
	ledsCommand.AddCommand(ledsStreamlinedCommand)
	ledsCommand.AddCommand(ledsGradientCommand)

	ledsCommand.PersistentFlags().Float32VarP(&ledsBrightness, "brightness", "b", 1, "led brightness (0.0-1.0)")

	ledsStreamlinedCommand.Flags().Float32VarP(&streamlinedSpeed, "speed", "s", 0.5, "speed for the effect (0.0-1.0)")
	ledsGradientCommand.Flags().Float32VarP(&gradientSpeed, "speed", "s", 0.5, "speed for the effect (0.0-1.0)")
}
