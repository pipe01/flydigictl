package commands

import (
	"fmt"

	"github.com/pipe01/flydigictl/pkg/version"
	"github.com/spf13/cobra"
)

var versionCommand = &cobra.Command{
	Use:   "version",
	Short: "Shows version information for flydigictl and flydigid",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if terseOutput {
			fmt.Println(version.Version)
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
