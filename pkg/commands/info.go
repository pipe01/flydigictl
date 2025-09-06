package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var gamepadNames = map[int32]string{
	19: "Apex 2",
	20: "Vader 2",
	21: "Vader 2",
	22: "Vader 2 Pro",
	23: "Vader 2",
	24: "Apex 3",
	25: "Direwolf",
	26: "Apex 3",
	29: "Apex 3",
	27: "Direwolf",
	30: "fp1Fate",
	28: "Vader 3",
	80: "Vader 3 Pro",
	81: "Vader 3 Pro ONE PIECE",
	82: "Direwolf 2",
	83: "fp2ip",
	84: "k2",
	85: "Vader 4",
}

var infoCommand = &cobra.Command{
	Use:   "info",
	Short: "Get information about the connected gamepad",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return useConnection(func() error {
			info, err := dbusClient.GetDeviceInfo()
			if err != nil {
				return fmt.Errorf("get device info: %w", err)
			}

			dict := []struct {
				Key   string
				Value any
			}{
				{"Device", fmt.Sprintf("%d (%s)", info.DeviceId, gamepadNames[info.DeviceId])},
				{"Battery percentage", fmt.Sprintf("%d%%", info.BatteryPercent)},
				{"Connection type", strings.ToLower(info.ConnectionType.String())},
				{"CPU", fmt.Sprintf("%s (%s)", info.CpuType, info.CpuName)},
			}

			maxLen := 0
			for _, e := range dict {
				if len(e.Key) > maxLen {
					maxLen = len(e.Key)
				}
			}

			for _, e := range dict {
				space := strings.Repeat(" ", maxLen-len(e.Key))

				fmt.Printf("%s%s : %v\n", space, e.Key, e.Value)
			}

			return nil
		})
	},
}

func init() {
	rootCmd.AddCommand(infoCommand)
}
