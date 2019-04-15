package cmd

import (
	log "github.com/code-ready/crc/pkg/crc/logging"
	"github.com/code-ready/crc/pkg/crc/preflight"
	"github.com/code-ready/crc/pkg/os"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(setupCmd)
}

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "setup hypervisor",
	Long:  "setup hypervisor to run the cluster",
	Run: func(cmd *cobra.Command, args []string) {
		runSetup(args)
	},
	Hidden: true,
}

func runSetup(arguments []string) {
	if !os.CheckUserPrivilages() {
		log.Fatal("You need to run this command as root, try prefixing sudo to the command.")
	}
	preflight.SetupHost()
}
