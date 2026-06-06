package commands

import (
	"fmt"
	"slices"

	"github.com/pipe01/flydigictl/pkg/dbus/pb"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type triggerSide string

const (
	triggerLeft  triggerSide = "left"
	triggerRight triggerSide = "right"
)

var showDetails bool
var skipValidation bool

var raceOptions struct {
	InitialPos int
	Pressure   int
}
var recoilOptions struct {
	InitialPos      int
	InitialStrength int
	Intensity       int
	Frequency       int
	InputAfter      bool
}

var sniperOptions struct {
	InitialPos int
	Length     int
	Pressure   int
	InputAfter bool
}

var lockOptions struct {
	InitialPos int
}

var vibrationOptions struct {
	Coefficient int
	Threshold   int
	TravelRange int
	Frequency   int
}

func (s triggerSide) GetBean(b *pb.GamepadConfiguration) *pb.TriggerConfiguration {
	switch s {
	case triggerLeft:
		return b.LeftTrigger
	case triggerRight:
		return b.RightTrigger
	}
	panic("invalid trigger side")
}

func ValidateTrigger(cmd *cobra.Command, args []string) error {
	validIds := []int32{84}

	if skipValidation {
		return nil
	}

	var startPersist = persistConnection

	var response = useConnection(func() error {
		persistConnection = true
		info, err := dbusClient.GetDeviceInfo()

		if err != nil {
			return fmt.Errorf("get device info: %w", err)
		}

		if slices.Contains(validIds, info.DeviceId) == false {
			return fmt.Errorf("Trigger not supported for this device")
		}

		return nil
	})
	persistConnection = startPersist
	return response
}

var triggerCommand = &cobra.Command{
	Use:   "trigger",
	Short: "Configure triggers (Apex 4)",
}

func genTriggerCommand(side triggerSide) *cobra.Command {
	cmd := &cobra.Command{
		Use:     string(side),
		Short:   fmt.Sprintf("Manage %s trigger configuration", side),
		PreRunE: ValidateTrigger,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := readConfiguration(func(conf *pb.GamepadConfiguration) *pb.TriggerConfiguration {
				return side.GetBean(conf)
			})
			if err != nil {
				return fmt.Errorf("get configuration: %w", err)
			}

			switch mode := cfg.Mode.(type) {
			case *pb.TriggerConfiguration_Default:
				if terseOutput {
					fmt.Print("default")
					return nil
				}

				fmt.Printf("%s trigger mode: Default\n", cases.Title(language.English, cases.NoLower).String(string(side)))

			case *pb.TriggerConfiguration_Race:

				if terseOutput {
					fmt.Print("race")
					return nil
				}

				fmt.Printf("%s trigger mode: Race\n", cases.Title(language.English, cases.NoLower).String(string(side)))

				if showDetails {
					fmt.Printf("  %-22s : %v\n", "Initial position", mode.Race.InitialPos)
					fmt.Printf("  %-22s : %v\n", "Pressure", mode.Race.Pressure)
				}
			case *pb.TriggerConfiguration_Recoil:
				if terseOutput {
					fmt.Println("recoil")
					return nil
				}

				fmt.Printf("%s trigger mode: Recoil\n", cases.Title(language.English, cases.NoLower).String(string(side)))
				if showDetails {
					fmt.Printf("  %-22s : %v\n", "Initial position", mode.Recoil.InitialPos)
					fmt.Printf("  %-22s : %v\n", "Initial strength", mode.Recoil.InitialStrength)
					fmt.Printf("  %-22s : %v\n", "Intensity", mode.Recoil.Intensity)
					fmt.Printf("  %-22s : %v\n", "Frequency", mode.Recoil.Frequency)
					fmt.Printf("  %-22s : %v\n", "Input After Trigger", mode.Recoil.InputAfter)
				}
			case *pb.TriggerConfiguration_Sniper:
				if terseOutput {
					fmt.Println("sniper")
					return nil
				}
				fmt.Printf("%s trigger mode: Sniper\n", cases.Title(language.English, cases.NoLower).String(string(side)))
				if showDetails {
					fmt.Printf("  %-22s : %v\n", "Initial position", mode.Sniper.InitialPos)
					fmt.Printf("  %-22s : %v\n", "Length", mode.Sniper.Length)
					fmt.Printf("  %-22s : %v\n", "Pressure", mode.Sniper.Pressure)
					fmt.Printf("  %-22s : %v\n", "Input After Trigger", mode.Sniper.InputAfter)
				}
			case *pb.TriggerConfiguration_Lock:
				if terseOutput {
					fmt.Println("lock")
					return nil
				}
				fmt.Printf("%s trigger mode: Lock\n", cases.Title(language.English, cases.NoLower).String(string(side)))
				if showDetails {
					fmt.Printf("  %-22s : %v\n", "Initial position", mode.Lock.InitialPos)
				}
			case *pb.TriggerConfiguration_Vibration:
				if terseOutput {
					fmt.Println("vibration")
					return nil
				}
				fmt.Printf("%s trigger mode: Vibration\n", cases.Title(language.English, cases.NoLower).String(string(side)))
				if showDetails {
					fmt.Printf("  %-22s : %v\n", "Coefficient", mode.Vibration.Coefficient)
					fmt.Printf("  %-22s : %v\n", "Threshold", mode.Vibration.Threshold)
					fmt.Printf("  %-22s : %v\n", "Travel range", mode.Vibration.TravelRange)
					fmt.Printf("  %-22s : %v\n", "Frequency", mode.Vibration.Frequency)
				}
			default:
				if terseOutput {
					fmt.Println("unknown")
					return nil
				}
				fmt.Printf("%s trigger mode: Unknown\n", cases.Title(language.English, cases.NoLower).String(string(side)))
			}

			return nil
		},
	}

	cmd.AddCommand(&cobra.Command{
		Use:     "default",
		Short:   "Set this trigger to default mode",
		Args:    cobra.NoArgs,
		PreRunE: ValidateTrigger,
		RunE: func(cmd *cobra.Command, args []string) error {
			return modifyConfiguration(func(conf *pb.GamepadConfiguration) {
				cfg := side.GetBean(conf)
				cfg.Mode = &pb.TriggerConfiguration_Default{}
			})
		},
	})

	raceCmd := &cobra.Command{
		Use:     "race",
		Short:   "Set this trigger to race mode",
		Args:    cobra.NoArgs,
		PreRunE: ValidateTrigger,
		RunE: func(cmd *cobra.Command, args []string) error {
			if raceOptions.InitialPos < 0 || raceOptions.InitialPos > 192 {
				return fmt.Errorf("initial-pos must be between 0 and 192")
			}
			if raceOptions.Pressure < 1 || raceOptions.Pressure > 255 {
				return fmt.Errorf("pressure must be between 1 and 255")
			}
			return modifyConfiguration(func(conf *pb.GamepadConfiguration) {
				cfg := side.GetBean(conf)
				cfg.Mode = &pb.TriggerConfiguration_Race{
					Race: &pb.TriggerRace{
						InitialPos: int32(raceOptions.InitialPos),
						Pressure:   int32(raceOptions.Pressure),
					},
				}
			})
		},
	}
	raceCmd.Flags().IntVar(
		&raceOptions.InitialPos,
		"initial-pos",
		0,
		"Initial trigger position (0–192)",
	)
	raceCmd.Flags().IntVar(
		&raceOptions.Pressure,
		"pressure",
		30,
		"Trigger pressure (1–255)",
	)
	raceCmd.Flags().SortFlags = false
	cmd.AddCommand(raceCmd)

	recoilCmd := &cobra.Command{
		Use:     "recoil",
		Short:   "Set this trigger to recoil mode",
		Args:    cobra.NoArgs,
		PreRunE: ValidateTrigger,
		RunE: func(cmd *cobra.Command, args []string) error {
			if recoilOptions.InitialPos < 0 || recoilOptions.InitialPos > 192 {
				return fmt.Errorf("start-pos must be between 0 and 192")
			}
			if recoilOptions.InitialStrength < 1 || recoilOptions.InitialStrength > 255 {
				return fmt.Errorf("initial-strength must be between 1 and 255")
			}
			if recoilOptions.Intensity < 1 || recoilOptions.Intensity > 255 {
				return fmt.Errorf("intensity must be between 1 and 255")
			}
			if recoilOptions.Frequency < 1 || recoilOptions.Frequency > 255 {
				return fmt.Errorf("frequency must be between 1 and 255")
			}

			return modifyConfiguration(func(conf *pb.GamepadConfiguration) {
				cfg := side.GetBean(conf)
				cfg.Mode = &pb.TriggerConfiguration_Recoil{
					Recoil: &pb.TriggerRecoil{
						InitialPos:      int32(recoilOptions.InitialPos),
						InitialStrength: int32(recoilOptions.InitialStrength),
						Intensity:       int32(recoilOptions.Intensity),
						Frequency:       int32(recoilOptions.Frequency),
						InputAfter:      recoilOptions.InputAfter,
					},
				}
			})
		},
	}
	recoilCmd.Flags().IntVar(
		&recoilOptions.InitialPos,
		"start-pos",
		0,
		"Vibration start position (0–192)",
	)
	recoilCmd.Flags().IntVar(
		&recoilOptions.InitialStrength,
		"initial-strength",
		1,
		"Initial recoil strength (1–255)",
	)
	recoilCmd.Flags().IntVar(
		&recoilOptions.Intensity,
		"intensity",
		50,
		"Recoil intensity (1–255). Force required to trigger vibration once the trigger reaches the start position",
	)
	recoilCmd.Flags().IntVar(
		&recoilOptions.Frequency,
		"frequency",
		15,
		"Recoil frequency (1–255)",
	)
	recoilCmd.Flags().BoolVar(
		&recoilOptions.InputAfter,
		"input-after",
		true,
		"Only triggers the input after start position",
	)
	cmd.AddCommand(recoilCmd)
	recoilCmd.Flags().SortFlags = false

	sniperCmd := &cobra.Command{
		Use:     "sniper",
		Short:   "Set this trigger to sniper mode",
		Args:    cobra.NoArgs,
		PreRunE: ValidateTrigger,
		RunE: func(cmd *cobra.Command, args []string) error {
			if sniperOptions.InitialPos < 0 || sniperOptions.InitialPos > 192 {
				return fmt.Errorf("initial-pos must be between 0 and 192")
			}
			if sniperOptions.Length < 1 || sniperOptions.Length > 255 {
				return fmt.Errorf("length must be between 1 and 255")
			}
			if sniperOptions.Pressure < 1 || sniperOptions.Pressure > 255 {
				return fmt.Errorf("pressure must be between 1 and 255")
			}

			return modifyConfiguration(func(conf *pb.GamepadConfiguration) {
				cfg := side.GetBean(conf)
				cfg.Mode = &pb.TriggerConfiguration_Sniper{
					Sniper: &pb.TriggerSniper{
						InitialPos: int32(sniperOptions.InitialPos),
						Length:     int32(sniperOptions.Length),
						Pressure:   int32(sniperOptions.Pressure),
						InputAfter: sniperOptions.InputAfter,
					},
				}
			})
		},
	}
	sniperCmd.Flags().IntVar(
		&sniperOptions.InitialPos,
		"initial-pos",
		50,
		"Initial trigger position (0–192)",
	)
	sniperCmd.Flags().IntVar(
		&sniperOptions.Length,
		"length",
		30,
		"Trigger length (1–255)",
	)
	sniperCmd.Flags().IntVar(
		&sniperOptions.Pressure,
		"pressure",
		1,
		"Trigger pressure (1–255)",
	)
	sniperCmd.Flags().BoolVar(
		&sniperOptions.InputAfter,
		"input-after",
		true,
		"Only triggers the input after start position",
	)
	sniperCmd.Flags().SortFlags = false
	cmd.AddCommand(sniperCmd)

	lockCmd := &cobra.Command{
		Use:     "lock",
		Short:   "Set this trigger to lock mode",
		Args:    cobra.NoArgs,
		PreRunE: ValidateTrigger,
		RunE: func(cmd *cobra.Command, args []string) error {
			if lockOptions.InitialPos < 0 || lockOptions.InitialPos > 192 {
				return fmt.Errorf("initial-pos must be between 0 and 192")
			}
			return modifyConfiguration(func(conf *pb.GamepadConfiguration) {
				cfg := side.GetBean(conf)
				cfg.Mode = &pb.TriggerConfiguration_Lock{
					Lock: &pb.TriggerLock{
						InitialPos: int32(lockOptions.InitialPos),
					},
				}
			})
		},
	}
	lockCmd.Flags().IntVar(
		&lockOptions.InitialPos,
		"initial-pos",
		40,
		"Lock start position (0–192)",
	)
	lockCmd.Flags().SortFlags = false
	cmd.AddCommand(lockCmd)

	vibrationCmd := &cobra.Command{
		Use:     "vibration",
		Short:   "Set this trigger to vibration mode",
		Args:    cobra.NoArgs,
		PreRunE: ValidateTrigger,
		RunE: func(cmd *cobra.Command, args []string) error {
			if vibrationOptions.Coefficient < 0 || vibrationOptions.Coefficient > 200 {
				return fmt.Errorf("intensity must be between 0 and 200")
			}
			if vibrationOptions.Threshold < 1 || vibrationOptions.Threshold > 255 {
				return fmt.Errorf("threshold must be between 1 and 255")
			}
			if vibrationOptions.TravelRange < 1 || vibrationOptions.TravelRange > 200 {
				return fmt.Errorf("range must be between 1 and 200")
			}
			if vibrationOptions.Frequency < 1 || vibrationOptions.Frequency > 255 {
				return fmt.Errorf("frequency must be between 1 and 255")
			}

			return modifyConfiguration(func(conf *pb.GamepadConfiguration) {
				cfg := side.GetBean(conf)
				cfg.Mode = &pb.TriggerConfiguration_Vibration{
					Vibration: &pb.TriggerVibration{
						Coefficient: int32(vibrationOptions.Coefficient),
						Threshold:   int32(vibrationOptions.Threshold),
						TravelRange: int32(vibrationOptions.TravelRange),
						Frequency:   int32(vibrationOptions.Frequency),
					},
				}
			})
		},
	}
	vibrationCmd.Flags().IntVar(
		&vibrationOptions.Coefficient,
		"intensity",
		50,
		"Vibration intensity (0–200)",
	)
	vibrationCmd.Flags().IntVar(
		&vibrationOptions.Threshold,
		"threshold",
		10,
		"Vibration threshold, below this value no vibration occurs (1–255)",
	)
	vibrationCmd.Flags().IntVar(
		&vibrationOptions.TravelRange,
		"range",
		20,
		"Trigger travel range for sustained vibration feedback (1-200)",
	)
	vibrationCmd.Flags().IntVar(
		&vibrationOptions.Frequency,
		"frequency",
		90,
		"Vibration frequency (1-255)",
	)
	vibrationCmd.Flags().SortFlags = false
	cmd.AddCommand(vibrationCmd)

	return cmd
}

func init() {
	var triggerLeftCommand = genTriggerCommand(triggerLeft)
	var triggerRightCommand = genTriggerCommand(triggerRight)

	triggerLeftCommand.Flags().BoolVar(
		&showDetails,
		"details",
		false,
		"Show detailed trigger configuration",
	)
	triggerRightCommand.Flags().BoolVar(
		&showDetails,
		"details",
		false,
		"Show detailed trigger configuration",
	)
	triggerLeftCommand.PersistentFlags().BoolVar(
		&skipValidation,
		"dangerous-skip",
		false,
		"Skips the trigger device validation. Only use it if you are sure the device has adaptive triggers",
	)
	triggerRightCommand.PersistentFlags().BoolVar(
		&skipValidation,
		"dangerous-skip",
		false,
		"Skips the trigger device validation. Only use it if you are sure the device has adaptive triggers",
	)

	triggerCommand.AddCommand(triggerLeftCommand)
	triggerCommand.AddCommand(triggerRightCommand)

	rootCmd.AddCommand(triggerCommand)
}
