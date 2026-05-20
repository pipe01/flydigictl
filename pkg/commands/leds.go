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

			case *pb.LedsConfiguration_Breathing:
				fmt.Printf("Breathing #%06X, speed %.0f%%\n", leds.Breathing.Color.Rgb, leds.Breathing.Speed*100)

			case *pb.LedsConfiguration_Gradient:
				fmt.Printf("Gradient #%06X -> #%06X, speed %.0f%%\n", leds.Gradient.StartColor.Rgb, leds.Gradient.EndColor.Rgb, leds.Gradient.Speed*100)

			case *pb.LedsConfiguration_Feedback:
				fmt.Printf("Feedback, speed %.0f%%\n", leds.Feedback.Speed*100)

			default:
				fmt.Println("Unsupported LED mode")
			}

			return nil
		})
	},
}

var ledsBrightness float32

var ledsOffCommand = &cobra.Command{
	Use:   "off",
	Short: "Turns all LEDs off",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return modifyLEDConfiguration(func(conf *pb.LedsConfiguration) {
			conf.Brightness = ledsBrightness
			conf.Leds = &pb.LedsConfiguration_Off{}
		})
	},
}

var ledsSteadyCommand = &cobra.Command{
	Use:     "steady",
	Short:   "Sets all LEDs to a solid color",
	Args:    cobra.ExactArgs(1),
	Example: "flydigictl leds steady #FF0000",
	RunE: func(cmd *cobra.Command, args []string) error {
		color, ok := parseColor(args[0])
		if !ok {
			return errInvalidColor
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
		return modifyLEDConfiguration(func(conf *pb.LedsConfiguration) {
			conf.Brightness = ledsBrightness
			conf.Leds = &pb.LedsConfiguration_Streamlined{Streamlined: &pb.LedsStreamlined{Speed: streamlinedSpeed}}
		})
	},
}

var breathingSpeed float32
var ledsBreathingCommand = &cobra.Command{
	Use:     "breathing COLOR",
	Short:   "Sets LEDs to a breathing effect",
	Args:    cobra.ExactArgs(1),
	Example: "flydigictl leds breathing '#BB00FF' --speed 0.5",
	RunE: func(cmd *cobra.Command, args []string) error {
		color, ok := parseColor(args[0])
		if !ok {
			return errInvalidColor
		}

		return modifyLEDConfiguration(func(conf *pb.LedsConfiguration) {
			conf.Brightness = ledsBrightness
			conf.Leds = &pb.LedsConfiguration_Breathing{Breathing: &pb.LedsBreathing{
				Color: color,
				Speed: breathingSpeed,
			}}
		})
	},
}

var gradientSpeed float32
var ledsGradientCommand = &cobra.Command{
	Use:     "gradient START_COLOR END_COLOR",
	Short:   "Sets LEDs to a two-color gradient effect",
	Args:    cobra.ExactArgs(2),
	Example: "flydigictl leds gradient #BB00FF #00F2FF",
	RunE: func(cmd *cobra.Command, args []string) error {
		startColor, ok := parseColor(args[0])
		if !ok {
			return errInvalidColor
		}

		endColor, ok := parseColor(args[1])
		if !ok {
			return errInvalidColor
		}

		return modifyLEDConfiguration(func(conf *pb.LedsConfiguration) {
			conf.Brightness = ledsBrightness
			conf.Leds = &pb.LedsConfiguration_Gradient{Gradient: &pb.LedsGradient{
				StartColor: startColor,
				EndColor:   endColor,
				Speed:      gradientSpeed,
			}}
		})
	},
}

var feedbackSpeed float32
var ledsFeedbackCommand = &cobra.Command{
	Use:   "feedback",
	Short: "Sets LEDs to a button feedback effect",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return modifyLEDConfiguration(func(conf *pb.LedsConfiguration) {
			conf.Brightness = ledsBrightness
			conf.Leds = &pb.LedsConfiguration_Feedback{Feedback: &pb.LedsFeedback{Speed: feedbackSpeed}}
		})
	},
}

func parseColor(s string) (*pb.Color, bool) {
	color, ok := pb.ColorFromName(s)
	if ok {
		return color, true
	}

	return pb.ColorFromHex(s)
}

func init() {
	rootCmd.AddCommand(ledsCommand)

	ledsCommand.AddCommand(ledsOffCommand)
	ledsCommand.AddCommand(ledsSteadyCommand)
	ledsCommand.AddCommand(ledsStreamlinedCommand)
	ledsCommand.AddCommand(ledsBreathingCommand)
	ledsCommand.AddCommand(ledsGradientCommand)
	ledsCommand.AddCommand(ledsFeedbackCommand)

	ledsCommand.PersistentFlags().Float32VarP(&ledsBrightness, "brightness", "b", 1, "led brightness (0.0-1.0)")

	ledsStreamlinedCommand.Flags().Float32VarP(&streamlinedSpeed, "speed", "s", 0.5, "speed for the effect (0.0-1.0)")
	ledsBreathingCommand.Flags().Float32VarP(&breathingSpeed, "speed", "s", 0.5, "speed for the effect (0.0-1.0)")
	ledsGradientCommand.Flags().Float32VarP(&gradientSpeed, "speed", "s", 0.5, "speed for the effect (0.0-1.0)")
	ledsFeedbackCommand.Flags().Float32VarP(&feedbackSpeed, "speed", "s", 0.5, "speed for the effect (0.0-1.0)")
}
