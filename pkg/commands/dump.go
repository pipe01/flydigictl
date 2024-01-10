package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var dumpNoColor bool

var dumpCommand = &cobra.Command{
	Use:   "dump",
	Short: "Dump raw gamepad configuration for debugging purposes",
	RunE: func(cmd *cobra.Command, args []string) error {
		return useConnection(func() error {
			dump, err := dbusClient.DumpConfiguration(dumpNoColor)
			if err != nil {
				return err
			}

			fmt.Println(dump)
			return nil
		})
	},
}

func init() {
	rootCmd.AddCommand(dumpCommand)

	dumpCommand.Flags().BoolVar(&dumpNoColor, "no-color", false, "disable ANSI color codes on output")
}
