package cmd

import (
	"github.com/spf13/cobra"
	"github.com/code-ready/crc/pkg/crc/daemon"
)

func init() {
	rootCmd.AddCommand(daemonCmd)
}

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Run the crc daemon",
	Long:  "Run the crc daemon",
	Run: func(cmd *cobra.Command, args []string) {
		daemon.RunDaemon()
	},
}
