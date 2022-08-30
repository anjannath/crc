package network

import (
	"errors"
	"fmt"
	"strings"

	"github.com/code-ready/crc/pkg/os/windows/powershell"
)

const hypervDefaultVirtualSwitchID = "c08cb7b8-9b3c-408e-8e30-5e16a3aeb444"

func SelectSwitchByNameOrDefault(name string) (bool, string) {
	// if named exists
	if ExistsSwitchByName(name) {
		return true, name
	}

	// else use Default
	return GetDefaultSwitchName()
}

func ExistsSwitchByName(name string) bool {
	getSwitchByNameCmd := fmt.Sprintf("Get-VMSwitch %s | ForEach-Object { $_.Name }", name)
	stdOut, stdErr, _ := powershell.Execute(getSwitchByNameCmd)

	// If stdErr contains the command then execution failed
	if strings.Contains(stdErr, "Get-VMSwitch") {
		return false
	}

	if strings.Contains(stdOut, name) {
		return true
	}
	return false
}

func GetDefaultSwitchName() (bool, string) {
	getDefaultSwitchNameCmd := fmt.Sprintf("[Console]::OutputEncoding = [Text.Encoding]::UTF8; Get-VMSwitch -Id %s | ForEach-Object { $_.Name }", hypervDefaultVirtualSwitchID)
	stdOut, stdErr, _ := powershell.Execute(getDefaultSwitchNameCmd)

	// If stdErr contains the command then execution failed
	if strings.Contains(stdErr, "Get-VMSwitch") {
		return false, ""
	}

	return true, strings.TrimSpace(stdOut)
}

// returns the IP from the first connected interface that has an IP address
func GetActiveIP() (string, error) {
	getInterfaceAliasCmd := `(Get-NetIPInterface -ConnectionState Connected -AddressFamily IPv4).InterfaceAlias`
	stdout, _, err := powershell.Execute(getInterfaceAliasCmd)
	if err != nil {
		return "", fmt.Errorf("Error trying to get interface names: %w", err)
	}

	ifaces := strings.Split(stdout, "\n")

	for _, iface := range ifaces {
		if strings.Contains(strings.ToLower(iface), "loopback") {
			continue
		}

		getIPAddrCmd := fmt.Sprintf("(Get-NetIPAddress -InterfaceAlias '%s' -AddressFamily IPv4).IPAddress", strings.TrimSpace(iface))
		stdout, _, err := powershell.Execute(getIPAddrCmd)
		if err == nil {
			return strings.TrimSpace(stdout), nil
		}
	}
	return "", errors.New("Unable to find an IP address assigned to any connected interface")
}
