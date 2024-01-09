package commands

import (
	"fmt"

	"github.com/pipe01/flydigi-linux/pkg/version"
	"github.com/spf13/cobra"
)

var versionCommand = &cobra.Command{
	Use:   "version",
	Short: "Shows version information for flydigictl and flydigid",
	Run: func(cmd *cobra.Command, args []string) {
		if terseOutput {
			println(version.Version)
		} else {
			fmt.Printf("flydigictl version %s\n", version.Version)
		}

		ver, err := dbusClient.GetServerVersion()
		if err != nil {
			fmt.Println("flydigid not running")
		} else {
			if terseOutput {
				fmt.Println(ver)
			} else {
				fmt.Printf("flydigid version %s\n", ver)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCommand)
}
