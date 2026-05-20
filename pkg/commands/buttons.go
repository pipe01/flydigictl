package commands

import (
	"fmt"
	"sort"
	"strings"

	"github.com/pipe01/flydigictl/pkg/dbus/pb"
	"github.com/spf13/cobra"
)

var buttonNameToID = map[string]int32{
	"UP":     0,
	"RIGHT":  1,
	"DOWN":   2,
	"LEFT":   3,
	"A":      4,
	"B":      5,
	"SELECT": 6,
	"X":      7,
	"Y":      8,
	"START":  9,
	"LB":     10,
	"RB":     11,
	"LT":     12,
	"RT":     13,
	"THUMBL": 14,
	"L3":     14,
	"THUMBR": 15,
	"R3":     15,
	"C":      16,
	"Z":      17,
	"M1":     18,
	"M2":     19,
	"M3":     20,
	"M4":     21,
	"M5":     22,
	"M6":     23,
	"MENU":   24,
	"HOME":   27,
	"BACK":   28,
}

var buttonIDToName = map[int32]string{
	0:  "UP",
	1:  "RIGHT",
	2:  "DOWN",
	3:  "LEFT",
	4:  "A",
	5:  "B",
	6:  "SELECT",
	7:  "X",
	8:  "Y",
	9:  "START",
	10: "LB",
	11: "RB",
	12: "LT",
	13: "RT",
	14: "THUMBL",
	15: "THUMBR",
	16: "C",
	17: "Z",
	18: "M1",
	19: "M2",
	20: "M3",
	21: "M4",
	22: "M5",
	23: "M6",
	24: "MENU",
	27: "HOME",
	28: "BACK",
}

var buttonsCommand = &cobra.Command{
	Use:   "buttons",
	Short: "Show or change controller button mappings",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return listButtonMappings()
	},
}

var buttonsMapCommand = &cobra.Command{
	Use:     "map SOURCE TARGET",
	Short:   "Maps a controller button to another gamepad button",
	Args:    cobra.ExactArgs(2),
	Example: "flydigictl buttons map M1 A",
	RunE: func(cmd *cobra.Command, args []string) error {
		sourceID, sourceName, err := parseButtonName(args[0])
		if err != nil {
			return err
		}

		targetID, targetName, err := parseButtonName(args[1])
		if err != nil {
			return err
		}

		return setButtonMapping(sourceID, sourceName, targetID, targetName)
	},
}

var buttonsResetCommand = &cobra.Command{
	Use:     "reset SOURCE",
	Short:   "Resets a controller button to its native mapping",
	Args:    cobra.ExactArgs(1),
	Example: "flydigictl buttons reset C",
	RunE: func(cmd *cobra.Command, args []string) error {
		sourceID, sourceName, err := parseButtonName(args[0])
		if err != nil {
			return err
		}

		return setButtonMapping(sourceID, sourceName, sourceID, sourceName)
	},
}

func init() {
	buttonsCommand.AddCommand(buttonsMapCommand)
	buttonsCommand.AddCommand(buttonsResetCommand)
	rootCmd.AddCommand(buttonsCommand)
}

func listButtonMappings() error {
	return useConnection(func() error {
		cfg, err := dbusClient.GetConfiguration()
		if err != nil {
			return fmt.Errorf("get config: %w", err)
		}

		for _, m := range cfg.ButtonMappings {
			if !m.IsOwn && !m.IsMapped {
				continue
			}

			owned := ""
			if !m.IsOwn {
				owned = " (not present)"
			}

			fmt.Printf("%-8s -> %-8s%s\n", displayButtonName(m.KeyId, m.Key), displayButtonName(m.MappedKeyId, m.MappedKey), owned)
		}

		return nil
	})
}

func setButtonMapping(sourceID int32, sourceName string, targetID int32, targetName string) error {
	return modifyConfiguration(func(conf *pb.GamepadConfiguration) {
		for _, m := range conf.ButtonMappings {
			if m.KeyId != sourceID {
				continue
			}

			m.Key = sourceName
			m.MappedKeyId = targetID
			m.MappedKey = targetName
			m.IsMapped = sourceID != targetID
			return
		}

		conf.ButtonMappings = append(conf.ButtonMappings, &pb.ButtonMapping{
			KeyId:       sourceID,
			Key:         sourceName,
			MappedKeyId: targetID,
			MappedKey:   targetName,
			IsOwn:       true,
			IsMapped:    sourceID != targetID,
		})
	})
}

func parseButtonName(raw string) (int32, string, error) {
	name := strings.ToUpper(strings.TrimSpace(raw))
	name = strings.ReplaceAll(name, "-", "")
	name = strings.ReplaceAll(name, "_", "")

	id, ok := buttonNameToID[name]
	if !ok {
		return 0, "", fmt.Errorf("unknown button %q; valid buttons: %s", raw, strings.Join(validButtonNames(), ", "))
	}

	return id, buttonIDToName[id], nil
}

func displayButtonName(id int32, fallback string) string {
	if name, ok := buttonIDToName[id]; ok {
		return name
	}
	if fallback != "" {
		return fallback
	}
	return fmt.Sprintf("%d", id)
}

func validButtonNames() []string {
	names := make([]string, 0, len(buttonNameToID))
	seen := map[string]bool{}

	for _, name := range buttonIDToName {
		if seen[name] {
			continue
		}
		seen[name] = true
		names = append(names, name)
	}

	sort.Strings(names)
	return names
}
