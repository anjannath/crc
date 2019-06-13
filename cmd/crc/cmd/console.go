package cmd

import (
	"github.com/code-ready/crc/pkg/crc/constants"
	"github.com/code-ready/crc/pkg/crc/errors"
	"github.com/code-ready/crc/pkg/crc/output"
	"github.com/pkg/browser"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(consoleCmd)
	consoleCmd.Flags().BoolVar(&showURL, "url", false, "Print Web Console URL")
	consoleCmd.Flags().BoolVar(&getOAuthToken, "oauth", false, "Open OAuth token URL in browser")
}

var (
	showURL bool
	getOAuthToken bool
)

var consoleCmd = &cobra.Command{
	Use:   "console",
	Short: "Open Web console in default browser",
	Long:  "Opens OpenShift web console in default browser",
	Run: func(cmd *cobra.Command, args []string) {
		if showURL {
			printConsoleURL()
			errors.Exit(0)
		}
		if getOAuthToken {
			oauthURL := constants.DefaultWebConsoleURL + "/request/oauth"
			launchURLInDefaultBrowser(oauthURL)
			errors.Exit(0)
		}
		launchURLInDefaultBrowser(constants.DefaultWebConsoleURL)
	},
}

func launchURLInDefaultBrowser(url string) {
	if err := browser.OpenURL(url); err != nil {
		errors.ExitWithMessage(1, "Failed to open URL in browser:", url, err.Error())
	}
}

func printConsoleURL() {
	output.Out(constants.DefaultWebConsoleURL)
}